package main

import (
	"errors"
	"github.com/alexdemen/migra/pkg"
	"os"
)

var InvalidCommandError = errors.New("invalid command")

func main() {
	err := processArgs(os.Args)
	if err != nil {

	}
}

func processArgs(args []string) error {
	if len(args) <= 1 {
		return InvalidCommandError
	}

	switch args[1] {
	case "create":
		var name, dir string
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		if len(args) > 2 {
			name = args[2]
		}
		return pkg.CreateMigration(name, dir)
	case "up":
	case "down":
	case "redo":
	case "status":
	default:
		return InvalidCommandError
	}

	return nil
}
