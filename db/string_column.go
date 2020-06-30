package db

import (
	"bufio"
	"gappkit/compare"
	"io"
)

var StringType StringColumnType

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

type StringColumn struct {
	BasicColumn
}

func (self *StringColumn) Init(name string) *StringColumn {
	self.BasicColumn.Init(name, &StringType)
	return self
}
