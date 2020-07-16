package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"reflect"
)

var (
	BoolType BoolColumnType
	boolValueType = reflect.TypeOf(false)
)

type BoolColumnType struct {
}

func (self *BoolColumnType) Compare(x, y interface{}) compare.Order {
	return compare.Bool(x.(bool), y.(bool))
}

func (self *BoolColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeBool(in)
}

func (self *BoolColumnType) Encode(val interface{}, out io.Writer) error {
	return EncodeBool(val.(bool), out)
}

func (self *BoolColumnType) ValueType() reflect.Type {
	return boolValueType
}

type BoolColumn struct {
	BasicColumn
}

func (self *BoolColumn) Init(name string) *BoolColumn {
	self.BasicColumn.Init(name, &BoolType)
	return self
}
