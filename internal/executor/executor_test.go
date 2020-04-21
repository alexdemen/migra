package executor

import (
	"github.com/alexdemen/migra/internal/driver"
	"testing"
)

func TestNewExecutor(t *testing.T) {

	pg, err := driver.NewPgStorage("postgres://admin:admin@127.0.0.1:5432/testdb")
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = NewExecutor(pg)
	if err != nil {
		t.Error(err.Error())
	}
}
