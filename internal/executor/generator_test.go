package executor

import (
	"os"
	"path"
	"testing"
)

func TestGenerateMigration(t *testing.T) {
	curDir, err := os.Getwd()
	if err != nil {

	}
	err = os.Mkdir(path.Join(curDir, "temp"), os.ModePerm)
	if err != nil {
		t.Fatal("failed to create temp directory")
	}
	defer os.RemoveAll("temp")
	err = GenerateMigration("sql", "test", "temp")
	if err != nil {
		t.Fatal("failed to create migration")
	}
}
