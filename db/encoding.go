package db

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(RecordSet{})
	gob.Register(time.Now())
}
