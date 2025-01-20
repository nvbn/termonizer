package model

import "time"

type Setting struct {
	ID      string
	Value   string
	Updated time.Time
}
