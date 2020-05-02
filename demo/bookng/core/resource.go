package core

import (
	"gappkit/db"
	"github.com/pkg/errors"
	"time"
)

type Resource struct {
	db *DB
	id db.RecordId
	Categories db.RecordSet
	Name string
	StartTime, EndTime time.Time
	Quantity int64
}

func (self *Resource) Init(db *DB, id db.RecordId) *Resource {
	self.db = db
	self.id = id
	return self
}

func (self *Resource) Id() db.RecordId {
	return self.id
}

func (self *Resource) UpdateQuantity(startTime, endTime time.Time, total, available int64) error {
	in := self.db.QuantityIndex.FindLower(self.id, startTime)
	var out []*Quantity
	
	for in.Next() && in.Key(1).(time.Time).Before(endTime) {
		q, err := self.db.LoadQuantity(in.Value())

		if err != nil {
			return err
		}
		
		if q.StartTime.Before(startTime) {
			head := self.db.NewQuantity(self, q.StartTime, startTime, q.Total, q.Available)
			out = append(out, head)
		}
		
		q.Available += available + total
		q.Total += total

		if q.Available < 0 {
			return NewOverbook(self, q.StartTime, q.EndTime, q.Available)
		}
		
		out = append(out, q)
		
		if q.EndTime.After(endTime) {
			tail := self.db.NewQuantity(self, endTime, q.EndTime, q.Total, q.Available)
			out = append(out, tail)
		}
	}

	for _, q := range out {
		if err := q.Store(); err != nil {
			return errors.Wrap(err, "Failed storing quantity")
		}
	}

	return nil
}

func (self *Resource) Store() error {
	var prev *Resource
	var err error
	
	if self.db.Resource.Exists(self.id) {
		if prev, err = self.db.LoadResource(self.id); err != nil {
			return errors.Wrapf(err, "Failed loading resource: %v", self.id)
		}		
	} else {
		if err = self.db.NewQuantity(self, MinTime, MaxTime, 0, 0).Store(); err != nil {
			return err
		}
	}

	var out db.Record
	out.Init(self.id)
	out.Set(&self.db.ResourceCategories, self.Categories)
	out.Set(&self.db.ResourceName, self.Name)
	out.Set(&self.db.ResourceStartTime, self.StartTime)
	out.Set(&self.db.ResourceEndTime, self.EndTime)
	out.Set(&self.db.ResourceQuantity, self.Quantity)
	
	if err = self.db.Resource.Store(out); err != nil {
		return err
	}

	if err = self.UpdateQuantity(self.StartTime, self.EndTime, self.Quantity, self.Quantity);
	err != nil {
		return errors.Wrap(err, "Failed updating quantity")
	}

	if prev != nil {
		q := -prev.Quantity
		
		if err = prev.UpdateQuantity(prev.StartTime, prev.EndTime, q, q); err != nil {
			return errors.Wrap(err, "Failed reverting quantity update")
		}
	}

	return nil
}

func (self *DB) LoadResource(id db.RecordId) (*Resource, error) {
	in, err := self.Resource.Load(id)
	
	if err != nil {
		return nil, errors.Wrapf(err, "Failed loading resource: %v", id)
	}

	out := new(Resource).Init(self, id)
	out.Name = in.Get(&self.ResourceName).(string)
	out.Categories = in.Get(&self.ResourceCategories).(db.RecordSet)
	out.StartTime = in.Get(&self.ResourceStartTime).(time.Time)
	out.EndTime = in.Get(&self.ResourceEndTime).(time.Time)
	out.Quantity = in.Get(&self.ResourceQuantity).(int64)
	return out, nil
}

func (self *DB) NewResource() *Resource {
	r := new(Resource).Init(self, self.Resource.NextId())
	r.Categories.Add(r.id)
	r.StartTime = MinTime
	r.EndTime = MaxTime
	r.Quantity = 1
	return r
}
