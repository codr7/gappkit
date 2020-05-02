package db

import (
	"bufio"
	"gappkit/compare"
	"io"
)

type Int64Column struct {
	BasicColumn
}

func (self *Int64Column) Init(name string) *Int64Column {
	self.BasicColumn.Init(name)
	return self
}

func (self *Int64Column) Compare(x, y interface{}) compare.Order {
	return compare.Int64(x.(int64), y.(int64))
}

func (self *Int64Column) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeInt(in)
}

func (self *Int64Column) Encode(val interface{}, out io.Writer) error {
	return EncodeInt(val.(int64), out)
}
