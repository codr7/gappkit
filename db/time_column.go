package db

import (
	"time"
)

type TimeColumn struct {
	BasicColumn
}

func (self *TimeColumn) Init(name string) *TimeColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *TimeColumn) Compare(x, y interface{}) Order {
	return CompareTime(x.(time.Time), y.(time.Time))
}

