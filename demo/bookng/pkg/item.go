package bookng

import (
	"gappkit/db"
	"github.com/pkg/errors"
	"time"
)

type Item struct {
	BasicModel
	Resource db.RecordId
	StartTime, EndTime time.Time
	Quantity int64
}

func (self *Item) Init(db *DB, id db.RecordId) *Item {
	self.BasicModel.Init(db, id)
	return self
}

func (self *Item) SetStartTime(val time.Time) {
	self.StartTime = val
	self.EndTime = val.AddDate(0, 0, 1)
}

func (self *Item) Store() error {
	var err error
	
	if self.db.Item.Exists(self.Id()) {
		var prev *Item
		
		if prev, err = self.db.LoadItem(self.Id()); err != nil {
			return errors.Wrapf(err, "Failed loading item: %v", self.Id())
		}

		var resource *Resource
		
		if resource, err = self.db.LoadResource(prev.Resource); err != nil {
			return err
		}

		if err = resource.UpdateQuantity(prev.StartTime, prev.EndTime, 0, prev.Quantity); err != nil {
			return errors.Wrap(err, "Failed reverting quantity update")
		}
	}

	if err = db.Store(self); err != nil {
		return err
	}

	var resource *Resource
	
	if resource, err = self.db.LoadResource(self.Resource); err != nil {
		return err
	}

	if err = resource.UpdateQuantity(self.StartTime, self.EndTime, 0, -self.Quantity); err != nil {
		return errors.Wrap(err, "Failed updating quantity")
	}

	return nil
}

func (self *Item) Table() *db.Table {
	return &self.db.Item
}

func (self *DB) LoadItem(id db.RecordId) (*Item, error) {
	out := new(Item).Init(self, id)

	if err := db.Load(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (self *DB) NewItem() *Item {
	i := new(Item).Init(self, self.Item.NextId())
	i.SetStartTime(Today())
	i.Quantity = 1
	return i
}
