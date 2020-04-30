package db

import (
	"gappkit/compare"
)

type StringColumn struct {
	BasicColumn
}

func (self *StringColumn) Init(name string) *StringColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *StringColumn) Compare(x, y interface{}) compare.Order {
	return compare.String(x.(string), y.(string))
}

