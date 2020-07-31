package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"reflect"
)

var (
	Int64Type Int64ColumnType
	int64ValueType = reflect.TypeOf(int64(0))
)

type Int64ColumnType struct {
}

func (self *Int64ColumnType) Compare(x, y interface{}) compare.Order {
	return compare.Int(x.(int64), y.(int64))
}

func (self *Int64ColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeInt(in)
}

func (self *Int64ColumnType) Encode(val interface{}, out io.Writer) error {
	return EncodeInt(val.(int64), out)
}

func (self *Int64ColumnType) ValueType() reflect.Type {
	return int64ValueType
}

type Int64Column struct {
	BasicColumn
}

func (self *Int64Column) Init(name string) *Int64Column {
	self.BasicColumn.Init(name, &Int64Type)
	return self
}
