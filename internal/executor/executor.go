package executor

import (
	"errors"
	"fmt"
	"github.com/alexdemen/migra/internal/core"
	"github.com/alexdemen/migra/internal/driver"
	"time"
)

type Executor struct {
	storage driver.Storage
}

func NewExecutor(storage driver.Storage) (*Executor, error) {
	ex := Executor{
		storage: storage,
	}

	return &ex, ex.check()
}

func (e Executor) ApplyMigrations(migrations []core.Migration, policy int) chan core.Message {
	done := make(chan core.Message, len(migrations))

	go func() {
		for idx, migration := range migrations {
			err := e.storage.ExecMigration(migration)
			if err != nil {
				done <- core.Message{
					Info: fmt.Sprintf("migration %s has been apply", migration.Name),
					Time: time.Now(),
				}
				if policy == BreakPolicy {
					break
				} else if policy == RollbackPolicy {
					for i := idx - 1; i >= 0; i-- {
						err = e.storage.RollBackMigration(migrations[i].Name)
						if err != nil {
							done <- core.Message{
								Error: errors.New("can`t rollback migration " + migrations[i].Name),
							}
						}
					}
				}
			} else {
				done <- core.Message{
					Info: fmt.Sprintf("migration %s has been apply", migration.Name),
					Time: time.Now(),
				}
			}
		}

		close(done)
	}()

	return done
}

func (e Executor) check() error {
	err := e.storage.Ping()
	if err != nil {
		return err
	}

	err = e.storage.InitScheme()

	return nil
}
