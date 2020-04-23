package core

import (
	"gappkit/db"
)

type DB struct {
	db.Root
	Quantity QuantityTable

	QuantityResource db.RecordColumn
	QuantityStartTime db.TimeColumn
	QuantityEndTime db.TimeColumn
	QuantityTotal db.Int64Column
	QuantityAvailable db.Int64Column

	QuantityIndex db.Index
	
	Resource ResourceTable
	ResourceName db.StringColumn
	ResourceCategories db.RecordSetColumn

	ResourceNameIndex db.Index
}

func NewDB(path string) *DB {
	self := new(DB)
	self.Root.Init(path)

	self.Quantity.Init(&self.Root)
	self.Quantity.AddColumn(self.QuantityResource.Init("Resource"))
	self.Quantity.AddColumn(self.QuantityStartTime.Init("StartTime"))
	self.Quantity.AddColumn(self.QuantityEndTime.Init("EndTime"))
	self.Quantity.AddColumn(self.QuantityTotal.Init("Total"))
	self.Quantity.AddColumn(self.QuantityAvailable.Init("Available"))

	self.QuantityIndex.Init(&self.Root, "quantity", true, &self.QuantityResource, &self.QuantityStartTime)
	self.Quantity.AddIndex(&self.QuantityIndex)
	
	self.Resource.Init(&self.Root)
	self.Resource.AddColumn(self.ResourceName.Init("Name"))
	self.Resource.AddColumn(self.ResourceCategories.Init("Categories"))

	self.ResourceNameIndex.Init(&self.Root, "resource_name", true, &self.ResourceName)
	self.Resource.AddIndex(&self.ResourceNameIndex)

	return self
}
