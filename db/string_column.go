package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"reflect"
)

var (
	StringType StringColumnType
	stringValueType = reflect.TypeOf("")
)

type StringColumnType struct {
}

func (self *StringColumnType) Compare(x, y interface{}) compare.Order {
	return compare.String(x.(string), y.(string))
}

func (self *StringColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeString(in)
}

func (self *StringColumnType) Encode(val interface{}, out io.Writer) error {
	return EncodeString(val.(string), out)
}

func (self *StringColumnType) ValueType() reflect.Type {
	return stringValueType
}

type StringColumn struct {
	BasicColumn
}

func (self *StringColumn) Init(name string) *StringColumn {
	self.BasicColumn.Init(name, &StringType)
	return self
}
