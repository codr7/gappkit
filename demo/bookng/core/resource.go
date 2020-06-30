package core

import (
	"gappkit/db"
	"github.com/pkg/errors"
	"math"
	"time"
)

type Resource struct {
	BasicModel
	Categories db.Slice
	Name string
	StartTime, EndTime time.Time
	Quantity int64
}

func (self *Resource) Init(db *DB, id db.RecordId) *Resource {
	self.BasicModel.Init(db, id)
	return self
}

func (self *Resource) AvailableQuantity(startTime, endTime time.Time) (int64, error) {
	in := self.db.QuantityIndex.FindLower(self.Id(), startTime)
	out := int64(math.MaxInt64)
	
	for in.Next() && in.Key(1).(time.Time).Before(endTime) {
		q, err := self.db.LoadQuantity(in.Value())

		if err != nil {
			return out, err
		}

		if q.Available < out {
			out = q.Available
		}
	}

	return out, nil
}

func (self *Resource) UpdateQuantity(startTime, endTime time.Time, total, available int64) error {
	in := self.db.QuantityIndex.FindLower(self.Id(), startTime)
	var out []*Quantity
	
	for in.Next() && in.Key(1).(time.Time).Before(endTime) {
		q, err := self.db.LoadQuantity(in.Value())

		if err != nil {
			return err
		}

		out = append(out, q)

		if q.StartTime.Before(startTime) {
			head := self.db.NewQuantity(self, q.StartTime, startTime, q.Total, q.Available)
			out = append(out, head)
			q.StartTime = startTime
		}

		if q.EndTime.After(endTime) {
			tail := self.db.NewQuantity(self, endTime, q.EndTime, q.Total, q.Available)
			out = append(out, tail)
			q.EndTime = endTime
		}

		q.Available += available
		q.Total += total

		if q.Available < 0 {
			return NewOverbook(self, q.StartTime, q.EndTime)
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
	
	if self.db.Resource.Exists(self.Id()) {
		if prev, err = self.db.LoadResource(self.Id()); err != nil {
			return errors.Wrapf(err, "Failed loading resource: %v", self.Id())
		}		
	} else {
		if err = self.db.NewQuantity(self, MinTime, MaxTime, 0, 0).Store(); err != nil {
			return err
		}
	}

	if err = db.Store(self); err != nil {
		return err
	}
	
	if err = self.UpdateQuantity(self.StartTime, self.EndTime, self.Quantity, self.Quantity); err != nil {
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

func (self *Resource) Table() *db.Table {
	return &self.db.Resource
}

func (self *DB) LoadResource(id db.RecordId) (*Resource, error) {
	out := new(Resource).Init(self, id)
	
	if err := db.Load(out); err != nil {
		return nil, err
	}
	
	return out, nil
}

func (self *DB) NewResource() *Resource {
	r := new(Resource).Init(self, self.Resource.NextId())
	r.StartTime = MinTime
	r.EndTime = MaxTime
	r.Quantity = 1
	return r
}
