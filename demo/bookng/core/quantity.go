package core

import (
	"gappkit/db"
	"time"
)

type Quantity struct {
	db *DB
	id db.RecordId
	Resource db.RecordId
	StartTime, EndTime time.Time
	Total, Available int64
}

func (self *Quantity) Init(db *DB, id db.RecordId) *Quantity {
	self.db = db
	self.id = id
	return self
}

func (self *Quantity) Id() db.RecordId {
	return self.id
}

func (self *Quantity) Store() error {
	var out db.Record
	out.Init(self.id)
	out.Set(&self.db.QuantityResource, self.Resource)
	out.Set(&self.db.QuantityStartTime, self.StartTime)
	out.Set(&self.db.QuantityEndTime, self.EndTime)
	out.Set(&self.db.QuantityTotal, self.Total)
	out.Set(&self.db.QuantityAvailable, self.Available)

	if err := self.db.Quantity.Store(out); err != nil {
		return err
	}

	return nil
}

func (self *DB) LoadQuantity(id db.RecordId) (*Quantity, error) {
	in, err := self.Quantity.Load(id)

	if err != nil {
		return nil, err
	}

	out := new(Quantity).Init(self, id)
	out.Resource = in.Get(&self.QuantityResource).(db.RecordId)
	out.StartTime = in.Get(&self.QuantityStartTime).(time.Time)
	out.EndTime = in.Get(&self.QuantityEndTime).(time.Time)
	out.Total = in.Get(&self.QuantityTotal).(int64)
	out.Available = in.Get(&self.QuantityAvailable).(int64)
	return out, nil
}

func (self *DB) NewQuantity(resource *Resource, startTime, endTime time.Time, total, available int64) *Quantity {
	q := new(Quantity).Init(self, self.Quantity.NextId())
	q.Resource = resource.id
	q.StartTime = startTime
	q.EndTime = endTime
	q.Total = total
	q.Available = available
	return q
}
