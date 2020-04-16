package core

import (
	"gappkit/db"
	"time"
)

type ResourceTable struct {
	db.BasicTable
}

type Resource struct {
	Record
	Name string
}

func (self *Resource) Store() error {
	if !self.exists {
		if err := self.db.NewQuantity(self, MinTime, MaxTime).Store(); err != nil {
			return err
		}
	}
	
	if err := self.db.Resource.Store(self.id, self); err != nil {
		return err
	}

	self.exists = true
	return nil
}


func (self *Resource) UpdateQuantity(startTime, endTime time.Time, total, available int) error {
	var in, out []*Quantity
	var err error
	
	for _, q := range in {
		if out, err = q.Update(startTime, endTime, total, available, out); err != nil {
			return err
		}
	}

	for _, q := range out {
		if err = q.Store(); err != nil {
			return err
		}
	}

	return nil
}

func (self *DB) NewResource() *Resource {
	r := new(Resource)
	r.Record.Init(self, self.Resource.NextId(), false)
	return r
}

func (self *DB) LoadResource(id db.RecordId) (*Resource, error) {
	r := new(Resource)
	r.Record.Init(self, id, true)
	
	if err := self.Resource.Load(id, r); err != nil {
		return nil, err
	}
	
	return r, nil
}
