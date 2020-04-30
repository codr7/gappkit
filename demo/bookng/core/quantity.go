package core

import (
	"gappkit/db"
	"time"
)

type Quantity struct {
	id db.RecordId
	Resource db.RecordId
	StartTime, EndTime time.Time
	Total, Available int
}

func (self *Quantity) Init(id db.RecordId) *Quantity {
	self.id = id
	return self
}

func (self *Quantity) Id() db.RecordId {
	return self.id
}

func (self *Quantity) Update(
	startTime, endTime time.Time,
	total, available int,
	out []*Quantity) ([]*Quantity, error) {
	return out, nil
}

func (self *DB) LoadQuantity(id db.RecordId) (*Quantity, error) {
	in, err := self.Quantity.Load(id)

	if err != nil {
		return nil, err
	}

	out := new(Quantity).Init(id)
	out.Resource = in.Get(&self.QuantityResource).(db.RecordId)
	out.StartTime = in.Get(&self.QuantityStartTime).(time.Time)
	out.EndTime = in.Get(&self.QuantityEndTime).(time.Time)
	out.Total = in.Get(&self.QuantityTotal).(int)
	out.Available = in.Get(&self.QuantityAvailable).(int)
	return out, nil
}

func (self *DB) NewQuantity(resource *Resource, startTime, endTime time.Time) *Quantity {
	q := new(Quantity).Init(self.Quantity.NextId())
	q.Resource = resource.id
	q.StartTime = startTime
	q.EndTime = endTime
	return q
}

func (self *DB) StoreQuantity(in *Quantity) error {
	var out db.Record
	out.Set(&self.QuantityResource, in.Resource)
	out.Set(&self.QuantityStartTime, in.StartTime)
	out.Set(&self.QuantityEndTime, in.EndTime)
	out.Set(&self.QuantityTotal, in.Total)
	out.Set(&self.QuantityAvailable, in.Available)

	if err := self.Quantity.Store(in.id, out); err != nil {
		return err
	}

	return nil
}
