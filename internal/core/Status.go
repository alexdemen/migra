package core

import "time"

type Status struct {
	Version int `db:"id"`
	Name    string
	Status  string
	Time    time.Time
}
