package core

import (
	"gappkit/db"
	"time"
)

type QuantityTable struct {
	db.BasicTable
}

func (self *QuantityTable) Init(root *db.Root) db.Table {
	self.BasicTable.Init(root, "quantity")
	self.NewColumn("ResourceId")
	self.NewColumn("StartTime")
	self.NewColumn("EndTime")
	self.NewColumn("Total")
	self.NewColumn("Available")
	root.AddTable(self)
	return self
}

func (self *QuantityTable) Load(id db.RecordId) (db.Record, error) {
	q := new(Quantity)
	q.BasicRecord.Init(id)

	if err := self.LoadRecord(id, q); err != nil {
		return nil, err
	}
	
	return q, nil
}

type Quantity struct {
	db.BasicRecord
	ResourceId db.RecordId
	StartTime, EndTime time.Time
	Total, Available int
}

func (self *Quantity) Store(db *DB) error {
	if err := db.Quantity.Store(self); err != nil {
		return err
	}

	return nil
}

func (self *Quantity) Update(
	startTime, endTime time.Time,
	total, available int,
	out []*Quantity) ([]*Quantity, error) {
	return out, nil
}

func (self *DB) NewQuantity(resource *Resource, startTime, endTime time.Time) *Quantity {
	q := new(Quantity)
	q.BasicRecord.Init(self.Quantity.NextId())
	q.ResourceId = resource.Id()
	return q
}
