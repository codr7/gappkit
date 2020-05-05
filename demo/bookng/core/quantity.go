package core

import (
	"gappkit/db"
	"time"
)

type Quantity struct {
	BasicModel
	Resource db.RecordId
	StartTime, EndTime time.Time
	Total, Available int64
}

func (self *Quantity) Init(db *DB, id db.RecordId) *Quantity {
	self.BasicModel.Init(db, id)
	return self
}

func (self *Quantity) Store() error {
	var out db.Record
	out.Init(self.Id())
	
	if err := self.db.Quantity.CopyFromModel(self, &out); err != nil {
		return err
	}

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

	if err = self.Quantity.CopyToModel(*in, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (self *DB) NewQuantity(resource *Resource, startTime, endTime time.Time, total, available int64) *Quantity {
	q := new(Quantity).Init(self, self.Quantity.NextId())
	q.Resource = resource.Id()
	q.StartTime = startTime
	q.EndTime = endTime
	q.Total = total
	q.Available = available
	return q
}
