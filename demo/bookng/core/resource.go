package core

import (
	"gappkit/db"
	"time"
)

type ResourceTable struct {
	db.BasicTable
}

func (self *ResourceTable) Load(id db.RecordId) (interface{}, error) {
	r := new(Resource)
	r.Record.Init(id, true)

	if err := self.LoadRecord(id, r); err != nil {
		return nil, err
	}
	
	return r, nil
}

type Resource struct {
	Record
	Name string
}

func (self *Resource) Store(db *DB) error {
	if !self.exists {
		if err := db.NewQuantity(self, MinTime, MaxTime).Store(db); err != nil {
			return err
		}
	}
	
	if err := db.Resource.Store(self.id, self); err != nil {
		return err
	}

	self.exists = true
	return nil
}


func (self *Resource) UpdateQuantity(db *DB, startTime, endTime time.Time, total, available int) error {
	var in, out []*Quantity
	var err error
	
	for _, q := range in {
		if out, err = q.Update(startTime, endTime, total, available, out); err != nil {
			return err
		}
	}

	for _, q := range out {
		if err = q.Store(db); err != nil {
			return err
		}
	}

	return nil
}

func (self *DB) NewResource() *Resource {
	r := new(Resource)
	r.Record.Init(self.Resource.NextId(), false)
	return r
}
