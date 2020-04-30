package db

import (
	"gappkit/util"
)

type Int64Column struct {
	BasicColumn
}

func (self *Int64Column) Init(name string) *Int64Column {
	self.BasicColumn.Init(name)
	return self
}

func (self *Int64Column) Compare(x, y interface{}) util.Order {
	return util.CompareInt64(x.(int64), y.(int64))
}

