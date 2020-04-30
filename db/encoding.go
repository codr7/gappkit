package db

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(RecordSet{})
	gob.Register(StoredField{})
	gob.Register(StoredRecord{})
	gob.Register(time.Now())
}
