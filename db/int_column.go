package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"reflect"
)

var (
	IntType IntColumnType
	intValueType = reflect.TypeOf(int(0))
)

type IntColumnType struct {
}

func (self *IntColumnType) Compare(x, y interface{}) compare.Order {
	return compare.Int(int64(x.(int)), int64(y.(int)))
}

func (self *IntColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	v, err := DecodeInt(in)
	return int(v), err
}

func (self *IntColumnType) Encode(val interface{}, out io.Writer) error {
	return EncodeInt(int64(val.(int)), out)
}

func (self *IntColumnType) ValueType() reflect.Type {
	return intValueType
}

type IntColumn struct {
	BasicColumn
}

func (self *IntColumn) Init(name string) *IntColumn {
	self.BasicColumn.Init(name, &IntType)
	return self
}
