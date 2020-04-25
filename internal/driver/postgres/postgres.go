package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alexdemen/migra/internal/core"
	"github.com/jmoiron/sqlx"
	"time"
)
import _ "github.com/jackc/pgx/stdlib"

const (
	StatusRegister = "REGISTER"
	StatusPending  = "PENDING"
	StatusReady    = "READY"
	StatusDown     = "DOWN"
	StatusRollback = "ROLLBACK"
	StatusFail     = "FAIL"

	enum = `
	CREATE TYPE migration_status AS ENUM (
    'REGISTER',
    'PENDING',
    'READY',
    'DOWN',
    'ROLLBACK',
    'FAIL'
	);`

	migrationTable = `
	create table migration
	(
		id bigserial not null,
		name varchar(200) not null,
		up varchar not null,
		down varchar,
		status migration_status not null,
		time timestamp not null
	);
	
	create unique index migration_id_uindex
	on migration (id);
	
	alter table migration
	add constraint migration_pk
	primary key (id);`

	checkQuery = `
	SELECT table_name FROM information_schema.tables
	WHERE table_schema NOT IN ('information_schema','pg_catalog')
	AND table_name = 'migration';`

	statusUpdate = `
	UPDATE  migration SET status = $1 WHERE id = $2
	`
)

type PgStorage struct {
	db          *sqlx.DB
	context     context.Context
	MessageChan chan core.Message
}

func NewPgStorage(ctx context.Context, dsn string) (*PgStorage, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	store := PgStorage{
		db:      db,
		context: ctx,
	}
	return &store, nil
}

func (p PgStorage) Ping() error {
	return p.db.Ping()
}

func (p PgStorage) InitScheme() error {
	exist, err := p.existingScheme()
	if exist || err != nil {
		return err
	}

	tx, err := p.db.Begin()

	_, err = tx.ExecContext(p.context, enum)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(p.context, migrationTable)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p PgStorage) RollbackLastMigration() error {
	getRollbackQuery := `
	SELECT id, down, name, status
	FROM migration
	WHERE id = (SELECT max(id) FROM migration);`

	migration := struct {
		Id     int64
		Query  string `db:"down"`
		Name   string
		Status string
	}{}

	err := p.db.QueryRowxContext(p.context, getRollbackQuery).StructScan(&migration)
	if err != nil {
		return err
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	p.db.ExecContext(p.context, statusUpdate, StatusDown, migration.Id)

	_, err = tx.ExecContext(p.context, migration.Query)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	p.db.ExecContext(p.context, statusUpdate, StatusRollback, migration.Id)

	return nil
}

func (p PgStorage) GetStatus() ([]core.Status, error) {
	statusQuery := `
	SELECT id, name, status, time
	FROM migration
	`

	rows, err := p.db.QueryxContext(p.context, statusQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]core.Status, 0)
	for rows.Next() {
		status := core.Status{}
		err = rows.StructScan(&status)
		if err != nil {
			return nil, err
		}
		result = append(result, status)
	}

	return result, nil
}

func (p PgStorage) ExecMigration(migration core.Migration) error {

	id, err := p.registerMigration(migration)
	if err != nil {
		return err
	}

	awaiter := time.NewTicker(time.Second)
	defer awaiter.Stop()
	err = func() error {
		for {
			select {
			case <-awaiter.C:
				{
					ready, err := p.checkOrder(migration.Name(), id)
					if err != nil {
						return err
					} else if ready {
						return nil
					} else {
						message := core.Message{
							Info: "Waiting",
							Time: time.Now(),
						}
						p.sendMessage(message)
					}
				}
			case <-p.context.Done():
				return nil
			}
		}
	}()

	err = p.exec(migration, id)

	return nil
}

func (p PgStorage) existingScheme() (bool, error) {
	resRow := p.db.QueryRowContext(p.context, checkQuery)
	var exist string
	err := resRow.Scan(&exist)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return exist != "", nil
}

func (p PgStorage) registerMigration(migration core.Migration) (int64, error) {
	sqlQuery := `
	INSERT INTO migration (name, up, down, status, time)
	VALUES ($1, $2, $3, $4, now())
	RETURNING id;
	`
	var lastId int64
	err := p.db.QueryRowContext(
		p.context,
		sqlQuery,
		migration.Name(),
		migration.Query(),
		migration.ReversQuery(),
		StatusRegister).Scan(&lastId)

	return lastId, err
}

func (p PgStorage) checkOrder(name string, id int64) (bool, error) {
	check := `
	SELECT min(id)
	FROM migration
	WHERE name = $1;`

	nextMigration := `
	SELECT min(id)
	FROM migration
	WHERE status != 'READY' AND status != 'FAIL';
	`

	var actualId int64
	err := p.db.QueryRowContext(p.context, check, name).Scan(&actualId)
	if err != nil {
		//TODO Удалить данную миграцию
		return false, err
	} else if id != actualId {
		//TODO удалить миграцию
		return false, core.MigrationExistError
	}

	var next int64
	err = p.db.QueryRowContext(p.context, nextMigration).Scan(&next)
	if err != nil {
		//TODO Удалить данную миграцию
		return false, err
	}

	if next != id {
		return false, nil
	} else {
		return true, nil
	}
}

func (p PgStorage) exec(migration core.Migration, id int64) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	tx.ExecContext(p.context, statusUpdate, StatusPending, id)

	_, err = p.db.ExecContext(p.context, migration.Query())
	if err != nil {
		tx.Rollback()
		p.db.ExecContext(p.context, statusUpdate, StatusFail, id)
		return err
	}
	tx.Commit()

	p.db.ExecContext(p.context, statusUpdate, StatusReady, id)

	return err
}

func (p PgStorage) sendMessage(message core.Message) {
	if p.MessageChan != nil {
		p.MessageChan <- message
	}
}
