package core

import (
	"gappkit/db"
)

type Record struct {
	db *DB
	id db.RecordId
	exists bool
}

func (self *Record) Init(db *DB, id db.RecordId, exists bool) {
	self.db = db
	self.id = id
	self.exists = exists
}

func (self *Record) Id() db.RecordId {
	return self.id
}
