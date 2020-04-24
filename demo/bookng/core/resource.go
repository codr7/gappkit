package core

import (
	"gappkit/db"
	"time"
)

type ResourceTable struct {
	db.BasicTable
}

func (self *ResourceTable) Init(root *db.Root) db.Table {
	self.BasicTable.Init(root, "resource")
	root.AddTable(self)
	return self
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
	Categories db.RecordSet
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
	in := db.QuantityIndex.FindLower(self.Id(), startTime)
	var out []*Quantity
	
	for in.Next() && in.Key(1).(time.Time).Before(endTime) {
		q, err := db.Quantity.Load(in.Value())
		
		if out, err = q.(*Quantity).Update(startTime, endTime, total, available, out); err != nil {
			return err
		}
	}

	for _, q := range out {
		if err := q.Store(db); err != nil {
			return err
		}
	}

	return nil
}

func (self *DB) NewResource() *Resource {
	r := new(Resource)
	r.BasicRecord.Init(self.Resource.NextId())
	r.Categories.Add(r.Id())
	return r
}
