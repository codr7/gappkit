package bookng

import (
	"gappkit/db"
)

type Model interface {
	db.Model
	DB() *DB
}

type BasicModel struct {
	db.BasicModel
	db *DB
}

func (self *BasicModel) Init(db *DB, id db.RecordId) {
	self.BasicModel.Init(id)
	self.db = db
}

func (self *BasicModel) DB() *DB {
	return self.db
}
