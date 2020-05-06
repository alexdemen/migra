package core

import "time"

type Message struct {
	Info  string
	Time  time.Time
	Error error
}
