package db

import (
	"gappkit/util"
)

type RecordSetColumn struct {
	BasicColumn
}

func (self *RecordSetColumn) Init(name string) *RecordSetColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *RecordSetColumn) Compare(x, y interface{}) util.Order {
	return x.(RecordSet).Compare(y.(RecordSet))
}

