package db

import (
	"bufio"
	"gappkit/compare"
	"io"
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

func (self *TimeColumn) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeTime(in)
}

func (self *TimeColumn) Encode(val interface{}, out io.Writer) error {
	return EncodeTime(val.(time.Time), out)
}
