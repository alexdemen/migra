package driver

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)
import _ "github.com/jackc/pgx/stdlib"

const (
	enum = `
	CREATE TYPE migration_status AS ENUM (
    'REGISTER',
    'PENDING',
    'READY',
    'DOWN',
    'ROLLBACK'
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
)

type PgStorage struct {
	db      *sqlx.DB
	context context.Context
}

func NewPgStorage(dsn string) (*PgStorage, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	store := PgStorage{
		db:      db,
		context: context.Background(),
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

func (p PgStorage) existingScheme() (bool, error) {
	resRow := p.db.QueryRowContext(p.context, checkQuery)
	var exist string
	err := resRow.Scan(&exist)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return exist != "", nil
}
