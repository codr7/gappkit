package db

import (
	"gappkit/util"
)

type StringColumn struct {
	BasicColumn
}

func (self *StringColumn) Init(name string) *StringColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *StringColumn) Compare(x, y interface{}) util.Order {
	return util.CompareString(x.(string), y.(string))
}

