package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"reflect"
)

type ColumnType interface {
	Compare(x, y interface{}) compare.Order
	Decode(in *bufio.Reader) (interface{}, error)
	Encode(val interface{}, out io.Writer) error
	ValueType() reflect.Type
}

type Column interface {
	ColumnType
	Name() string
	Type() ColumnType
}

type BasicColumn struct {
	name string
	columnType ColumnType
}

func (self *BasicColumn) Init(name string, columnType ColumnType) *BasicColumn {
	self.name = name
	self.columnType = columnType
	return self
}

func (self *BasicColumn) Name() string {
	return self.name
}

func (self *BasicColumn) Type() ColumnType {
	return self.columnType
}

func (self *BasicColumn) Compare(x, y interface{}) compare.Order {
	return self.columnType.Compare(x, y)
}

func (self *BasicColumn) Decode(in *bufio.Reader) (interface{}, error) {
	return self.columnType.Decode(in)
}

func (self *BasicColumn) Encode(val interface{}, out io.Writer) error {
	return self.columnType.Encode(val, out)
}

func (self *BasicColumn) ValueType() reflect.Type {
	return self.columnType.ValueType()
}
