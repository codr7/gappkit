package db

import (
	"gappkit/util"
	"time"
)

type TimeColumn struct {
	BasicColumn
}

func (self *TimeColumn) Init(name string) *TimeColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *TimeColumn) Compare(x, y interface{}) util.Order {
	return util.CompareTime(x.(time.Time), y.(time.Time))
}

