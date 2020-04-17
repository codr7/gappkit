package core

import (
	"gappkit/db"
)

type DB struct {
	db.Root
	Quantity QuantityTable
	Resource ResourceTable
}

func NewDB(path string) *DB {
	self := new(DB)
	self.Root.Init(path)
	self.Quantity.Init(&self.Root)
	self.Resource.Init(&self.Root)
	return self
}
