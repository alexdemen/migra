package executor

import "github.com/alexdemen/migra/internal/driver"

type executor struct {
	storage driver.Storage
}

func (e executor) check() error {
	err := e.storage.Ping()
	if err != nil {
		return err
	}

	err = e.storage.InitScheme()

	return nil
}

func NewExecutor(storage driver.Storage) (*executor, error) {
	ex := executor{
		storage: storage,
	}

	ex.check()

	return &ex, nil
}
