package core

import (
	"gappkit/db"
)

type DB struct {
	db.Root
	Calendar db.Table
	Resources db.Table
}

func NewDB(path string) *DB {
	self := new(DB)
	self.Root.Init(path)

	self.Calendar.Init(&self.Root, "calendar")
	self.Calendar.NewColumn("ResourceId")
	self.Calendar.NewColumn("StartTime")
	self.Calendar.NewColumn("EndTime")
	self.Calendar.NewColumn("Total")
	self.Calendar.NewColumn("Available")

	self.Resources.Init(&self.Root, "resource")
	self.Resources.NewColumn("Name")

	return self
}
