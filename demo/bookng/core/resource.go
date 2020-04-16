package core

import (
	"gappkit/db"
)

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
	
	if err := self.db.Resources.Store(self.id, self); err != nil {
		return err
	}

	self.exists = true
	return nil
}

func (self *DB) NewResource() *Resource {
	r := new(Resource)
	r.Record.Init(self, self.Resources.NextId(), false)
	return r
}

func (self *DB) LoadResource(id db.RecordId) (*Resource, error) {
	r := new(Resource)
	r.Record.Init(self, id, true)
	
	if err := self.Resources.Load(id, r); err != nil {
		return nil, err
	}
	
	return r, nil
}
