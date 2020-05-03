package db

import (
	"bufio"
	"gappkit/compare"
	"io"
)

type StringColumn struct {
	BasicColumn
}

func (self *StringColumn) Init(name string) *StringColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *StringColumn) Compare(x, y interface{}) compare.Order {
	return compare.String(x.(string), y.(string))
}

func (self *StringColumn) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeString(in)
}

func (self *StringColumn) Encode(val interface{}, out io.Writer) error {
	return EncodeString(val.(string), out)
}

