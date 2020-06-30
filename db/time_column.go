package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"reflect"
	"time"
)

var (
	TimeType TimeColumnType
	timeValueType = reflect.TypeOf(time.Now())
)

type TimeColumnType struct {
}

func (self *TimeColumnType) Compare(x, y interface{}) compare.Order {
	return compare.Time(x.(time.Time), y.(time.Time))
}

func (self *TimeColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeTime(in)
}

func (self *TimeColumnType) Encode(val interface{}, out io.Writer) error {
	return EncodeTime(val.(time.Time), out)
}

func (self *TimeColumnType) ValueType() reflect.Type {
	return timeValueType
}

type TimeColumn struct {
	BasicColumn
}

func (self *TimeColumn) Init(name string) *TimeColumn {
	self.BasicColumn.Init(name, &TimeType)
	return self
}
