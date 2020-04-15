package core

import (
	"gappkit/db"
)

type DB struct {
	db.Root
	Resources db.Table
}

func NewDB(path string) *DB {
	self := new(DB)
	self.Root.Init(path)

	self.Resources.Init(&self.Root, "resource")
	self.Resources.NewColumn("Name")

	return self
}
