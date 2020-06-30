package db

import (
	"bufio"
	"gappkit/compare"
	"io"
)

var Int64Type Int64ColumnType

type Int64ColumnType struct {
}

func (self *Int64ColumnType) Compare(x, y interface{}) compare.Order {
	return compare.Int64(x.(int64), y.(int64))
}

func (self *Int64ColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeInt(in)
}

func (self *Int64ColumnType) Encode(val interface{}, out io.Writer) error {
	return EncodeInt(val.(int64), out)
}


type Int64Column struct {
	BasicColumn
}

func (self *Int64Column) Init(name string) *Int64Column {
	self.BasicColumn.Init(name, &Int64Type)
	return self
}
