package db

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(Field{})
	gob.Register(Record{})
	gob.Register(RecordSet{})
	gob.Register(time.Now())
}
