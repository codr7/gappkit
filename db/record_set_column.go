package db

import (
	"gappkit/compare"
)

type RecordSetColumn struct {
	BasicColumn
}

func (self *RecordSetColumn) Init(name string) *RecordSetColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *RecordSetColumn) Compare(x, y interface{}) compare.Order {
	return x.(RecordSet).Compare(y.(RecordSet))
}

