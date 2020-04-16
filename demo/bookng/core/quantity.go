package core

import (
	"gappkit/db"
	"time"
)

type QuantityTable struct {
	db.BasicTable
}

func (self *QuantityTable) Load(id db.RecordId) (interface{}, error) {
	q := new(Quantity)
	q.Record.Init(id, true)

	if err := self.LoadRecord(id, q); err != nil {
		return nil, err
	}
	
	return q, nil
}

type Quantity struct {
	Record
	ResourceId db.RecordId
	StartTime, EndTime time.Time
	Total, Available int
}

func (self *Quantity) Store(db *DB) error {
	if err := db.Quantity.Store(self.id, self); err != nil {
		return err
	}

	self.exists = true
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
	q.Record.Init(self.Quantity.NextId(), false)
	q.ResourceId = resource.id
	return q
}
