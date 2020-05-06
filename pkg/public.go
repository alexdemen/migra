package pkg

import (
	"errors"
	"github.com/alexdemen/migra/internal/core"
)

func CreateMigration(name string, dir string) error {
	return errors.New("not implemented")
}

func UpMigrations(dir string) error {
	return errors.New("not implemented")
}

func DownMigration() error {
	return errors.New("not implemented")
}

func RedoMigration() error {
	return errors.New("not implemented")
}

func StatusMigrations() ([]core.Status, error) {
	return nil, errors.New("not implemented")
}
