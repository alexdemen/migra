package core

import "errors"

var MigrationExistError = errors.New("migration already exists")
var LastMigrationStatusError = errors.New("last migration not ready status")
