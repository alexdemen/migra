package driver

import "github.com/alexdemen/migra/internal/core"

type Storage interface {
	Ping() error
	InitScheme() error
	ExecMigration(migration core.Migration) error
	RollbackLastMigration() error
	GetStatus() ([]core.Status, error)
}
