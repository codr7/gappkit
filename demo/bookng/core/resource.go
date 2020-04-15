package core

import (
	"gappkit/db"
)

type Resource struct {
	Id db.RecordId
	Name string
}

func (self *Resource) Store(db *DB) error {
	return db.Resources.Store(self.Id, self)
}

func (self *DB) NewResource() *Resource {
	return &Resource{Id: self.Resources.NextId()}
}

func (self *DB) LoadResource(id db.RecordId) (*Resource, error) {
	r := &Resource{Id: id}

	if err := self.Resources.Load(id, r); err != nil {
		return nil, err
	}
	
	return r, nil
}
