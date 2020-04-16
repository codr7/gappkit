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

	self.Quantity.Init(&self.Root, "quantity")
	self.Quantity.NewColumn("ResourceId")
	self.Quantity.NewColumn("StartTime")
	self.Quantity.NewColumn("EndTime")
	self.Quantity.NewColumn("Total")
	self.Quantity.NewColumn("Available")
	self.AddTable(&self.Quantity)

	self.Resource.Init(&self.Root, "resource")
	self.Resource.NewColumn("Name")
	self.AddTable(&self.Resource)

	return self
}
