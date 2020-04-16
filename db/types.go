package db

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(time.Now())
}
