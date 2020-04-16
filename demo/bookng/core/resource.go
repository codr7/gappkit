package core

import (
	"gappkit/db"
	"time"
)

type ResourceTable struct {
	db.BasicTable
}

func (self *ResourceTable) Load(id db.RecordId) (db.Record, error) {
	r := new(Resource)
	r.BasicRecord.Init(id)

	if err := self.LoadRecord(id, r); err != nil {
		return nil, err
	}
	
	return r, nil
}

type Resource struct {
	db.BasicRecord
	Name string
}

func (self *Resource) Store(db *DB) error {
	if !db.Resource.Exists(self.Id()) {
		if err := db.NewQuantity(self, MinTime, MaxTime).Store(db); err != nil {
			return err
		}
	}
	
	if err := db.Resource.Store(self); err != nil {
		return err
	}

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
	r.BasicRecord.Init(self.Resource.NextId())
	return r
}
