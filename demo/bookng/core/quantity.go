package core

import (
	"gappkit/db"
	"time"
)

type Quantity struct {
	Record
	ResourceId db.RecordId
	StartTime, EndTime time.Time
	Total, Available int
}

func (self *Quantity) Store() error {
	if err := self.db.Calendar.Store(self.id, self); err != nil {
		return err
	}

	self.exists = true
	return nil
}

func (self *DB) NewQuantity(resource *Resource, startTime, endTime time.Time) *Quantity {
	q := new(Quantity)
	q.Record.Init(self, self.Calendar.NextId(), false)
	q.ResourceId = resource.id
	return q
}

func (self *DB) LoadQuantity(id db.RecordId) (*Quantity, error) {
	q := new(Quantity)
	q.Record.Init(self, id, true)

	if err := self.Calendar.Load(id, q); err != nil {
		return nil, err
	}
	
	return q, nil
}
