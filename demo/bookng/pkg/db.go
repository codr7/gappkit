package bookng

import (
	"gappkit/db"
)

type DB struct {
	db.Root
	
	Quantity db.Table
	QuantityResource db.RecordColumn
	QuantityStartTime db.TimeColumn
	QuantityEndTime db.TimeColumn
	QuantityTotal db.Int64Column
	QuantityAvailable db.Int64Column
	QuantityIndex db.Index
	
	Resource db.Table
	ResourceName db.StringColumn
	ResourceCategories db.SliceColumn
	ResourceStartTime db.TimeColumn
	ResourceEndTime db.TimeColumn
	ResourceQuantity db.Int64Column
	ResourceNameIndex db.Index

	Item db.Table
	ItemResource db.RecordColumn
	ItemStartTime db.TimeColumn
	ItemEndTime db.TimeColumn
	ItemQuantity db.Int64Column
}

func NewDB(path string) *DB {
	self := new(DB)
	self.Root.Init(path)

	self.Quantity.Init(&self.Root, "quantity")
	self.Quantity.AddColumn(self.QuantityResource.Init("Resource"))
	self.Quantity.AddColumn(self.QuantityStartTime.Init("StartTime"))
	self.Quantity.AddColumn(self.QuantityEndTime.Init("EndTime"))
	self.Quantity.AddColumn(self.QuantityTotal.Init("Total"))
	self.Quantity.AddColumn(self.QuantityAvailable.Init("Available"))
	
	self.QuantityIndex.Init(&self.Root, "quantity",
		&self.QuantityResource, &self.QuantityEndTime, &self.QuantityStartTime)
	self.Quantity.AddIndex(&self.QuantityIndex)
	
	self.Resource.Init(&self.Root, "resource")
	self.Resource.AddColumn(self.ResourceName.Init("Name"))
	self.Resource.AddColumn(self.ResourceCategories.Init("Categories", &db.RecordType))
	self.Resource.AddColumn(self.ResourceStartTime.Init("StartTime"))
	self.Resource.AddColumn(self.ResourceEndTime.Init("EndTime"))
	self.Resource.AddColumn(self.ResourceQuantity.Init("Quantity"))

	self.ResourceNameIndex.Init(&self.Root, "resource_name", &self.ResourceName)
	self.Resource.AddIndex(&self.ResourceNameIndex)

	self.Item.Init(&self.Root, "item")
	self.Item.AddColumn(self.ItemResource.Init("Resource"))
	self.Item.AddColumn(self.ItemStartTime.Init("StartTime"))
	self.Item.AddColumn(self.ItemEndTime.Init("EndTime"))
	self.Item.AddColumn(self.ItemQuantity.Init("Quantity"))
	
	return self
}
