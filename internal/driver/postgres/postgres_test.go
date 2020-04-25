package postgres

import (
	"context"
	"github.com/alexdemen/migra/internal/core"
	"testing"
)

func TestPgStorage_InitScheme(t *testing.T) {
	pg, err := NewPgStorage(context.Background(), "postgres://admin:admin@127.0.0.1:5432/testdb")
	if err != nil {
		t.Fatal(err.Error())
	}

	err = pg.InitScheme()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPgStorage_registerMigration(t *testing.T) {
	query := `
	create table test1(
		id bigserial not null
	);

	create table test12(
		id bigserial not null
	);
	`

	migration := core.NewMigration("test", query, "")

	pg, err := NewPgStorage(context.Background(), "postgres://admin:admin@127.0.0.1:5432/testdb")
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = pg.registerMigration(*migration)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPgStorage_ExecMigration(t *testing.T) {
	query := `
	create table test1(
		id bigserial not null
	);

	create table test12(
		id bigserial not null
	);
	`

	migration := core.NewMigration("test", query, "")

	pg, err := NewPgStorage(context.Background(), "postgres://admin:admin@127.0.0.1:5432/testdb")
	if err != nil {
		t.Fatal(err.Error())
	}

	err = pg.ExecMigration(*migration)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPgStorage_RollbackLastMigration(t *testing.T) {
	query := `
	create table test1(
		id bigserial not null
	);

	create table test12(
		id bigserial not null
	);
	`

	reversQuery := `
	DROP TABLE test1;

	DROP TABLE test12
	`

	migration := core.NewMigration("test", query, reversQuery)

	pg, err := NewPgStorage(context.Background(), "postgres://admin:admin@127.0.0.1:5432/testdb")
	if err != nil {
		t.Fatal(err.Error())
	}

	err = pg.ExecMigration(*migration)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = pg.RollbackLastMigration()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPgStorage_GetStatus(t *testing.T) {
	migration := core.NewMigration("test", "query", "")

	pg, err := NewPgStorage(context.Background(), "postgres://admin:admin@127.0.0.1:5432/testdb")
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = pg.registerMigration(*migration)
	if err != nil {
		t.Fatal(err.Error())
	}

	status, err := pg.GetStatus()
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(status) == 0 {
		t.Fatal("Invalid len if Status slice.")
	} else if status[0].Name != migration.Name() {
		t.Fatal("Invalid migration name.")
	}
}
