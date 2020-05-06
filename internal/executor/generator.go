package executor

import (
	"errors"
	"os"
	"path"
)

func GenerateMigration(mType string, name string, dir string) error {
	if mType == "sql" {
		return createSqlTemplate(name, dir)
	}
	return errors.New("invalid migration type")
}

func createSqlTemplate(name string, dir string) error {
	fUp, err := os.Create(path.Join(dir, name+".up.sql"))
	if err != nil {
		return err
	}
	fUp.Close()

	fDown, err := os.Create(path.Join(dir, name+".down.sql"))
	if err != nil {
		os.Remove(path.Join(dir, name+".up.sql"))
		return err
	}
	defer fDown.Close()

	return nil
}
