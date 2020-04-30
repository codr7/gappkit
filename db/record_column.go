package db

import (
	"gappkit/compare"
)

type RecordColumn struct {
	BasicColumn
}

func (self *RecordColumn) Init(name string) *RecordColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *RecordColumn) Compare(x, y interface{}) compare.Order {
	return CompareRecordId(x.(RecordId), y.(RecordId))
}

