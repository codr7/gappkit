package core

import (
	"gappkit/db"
	"time"
)

type QuantityTable struct {
	db.BasicTable
}

type Quantity struct {
	Record
	ResourceId db.RecordId
	StartTime, EndTime time.Time
	Total, Available int
}

func (self *Quantity) Store() error {
	if err := self.db.Quantity.Store(self.id, self); err != nil {
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
	q.Record.Init(self, self.Quantity.NextId(), false)
	q.ResourceId = resource.id
	return q
}

func (self *DB) LoadQuantity(id db.RecordId) (*Quantity, error) {
	q := new(Quantity)
	q.Record.Init(self, id, true)

	if err := self.Quantity.Load(id, q); err != nil {
		return nil, err
	}
	
	return q, nil
}
