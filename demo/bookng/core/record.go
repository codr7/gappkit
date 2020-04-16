package core

import (
	"gappkit/db"
)

type Record struct {
	id db.RecordId
	exists bool
}

func (self *Record) Init(id db.RecordId, exists bool) {
	self.id = id
	self.exists = exists
}

func (self *Record) Id() db.RecordId {
	return self.id
}
