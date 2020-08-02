package bookng

import (
	"gappkit/db"
	"time"
)

type Quantity struct {
	BasicModel
	Resource db.RecordId
	StartTime, EndTime time.Time
	Total, Available int64
}

func (self *Quantity) Init(db *DB, id db.RecordId) *Quantity {
	self.BasicModel.Init(db, id)
	return self
}

func (self *Quantity) Store() error {
	_, err := db.Store(self)
	return err
}

func (self *Quantity) Table() *db.Table {
	return &self.db.Quantity
}

func (self *DB) LoadQuantity(id db.RecordId) (*Quantity, error) {
	out := new(Quantity).Init(self, id)

	if err := db.Load(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (self *DB) NewQuantity(resource *Resource, startTime, endTime time.Time, total, available int64) *Quantity {
	q := new(Quantity).Init(self, self.Quantity.NextId())
	q.Resource = resource.Id()
	q.StartTime = startTime
	q.EndTime = endTime
	q.Total = total
	q.Available = available
	return q
}
