package core

import (
	"gappkit/db"
	"time"
)

type Resource struct {
	id db.RecordId
	Categories db.RecordSet
	Name string
}

func (self *Resource) Init(id db.RecordId) *Resource {
	self.id = id
	return self
}

func (self *Resource) Id() db.RecordId {
	return self.id
}

func (self *Resource) UpdateQuantity(db *DB, startTime, endTime time.Time, total, available int) error {
	in := db.QuantityIndex.FindLower(self.id, startTime)
	var out []*Quantity
	
	for in.Next() && in.Key(1).(time.Time).Before(endTime) {
		q, err := db.LoadQuantity(in.Value())
		
		if out, err = q.Update(startTime, endTime, total, available, out); err != nil {
			return err
		}
	}

	for _, q := range out {
		if err := db.StoreQuantity(q); err != nil {
			return err
		}
	}

	return nil
}

func (self *DB) LoadResource(id db.RecordId) (*Resource, error) {
	in, err := self.Resource.Load(id)
	
	if err != nil {
		return nil, err
	}

	out := new(Resource).Init(id)
	out.Name = in.Get(&self.ResourceName).(string)
	out.Categories = in.Get(&self.ResourceCategories).(db.RecordSet)
	return out, nil
}

func (self *DB) NewResource() *Resource {
	r := new(Resource).Init(self.Resource.NextId())
	r.Categories.Add(r.id)
	return r
}

func (self *DB) StoreResource(in *Resource) error {
	if !self.Resource.Exists(in.id) {
		if err := self.StoreQuantity(self.NewQuantity(in, MinTime, MaxTime)); err != nil {
			return err
		}
	}

	var out db.Record
	out.Set(&self.ResourceCategories, in.Categories)
	out.Set(&self.ResourceName, in.Name)
	
	if err := self.Resource.Store(in.id, out); err != nil {
		return err
	}

	return nil
}
