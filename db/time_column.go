package db

import (
	"gappkit/compare"
	"time"
)

type TimeColumn struct {
	BasicColumn
}

func (self *TimeColumn) Init(name string) *TimeColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *TimeColumn) Compare(x, y interface{}) compare.Order {
	return compare.Time(x.(time.Time), y.(time.Time))
}

