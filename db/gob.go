package db

import (
	"encoding/gob"
	"time"
)

func init() {
	Register(time.Now())
}

func Register(value interface{}) {
	gob.Register(value)
}
