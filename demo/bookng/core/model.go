package core

import (
	"gappkit/db"
)

type Model interface {
	Init(db *DB, id db.RecordId)
	DB() *DB
	Id() db.RecordId
	Store() error
}

type BasicModel struct {
	db *DB
	id db.RecordId
}

func (self *BasicModel) Init(db *DB, id db.RecordId) {
	self.db = db
	self.id = id
}

func (self *BasicModel) DB() *DB {
	return self.db
}

func (self *BasicModel) Id() db.RecordId {
	return self.id
}
