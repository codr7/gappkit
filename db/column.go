package db

import (
	"bufio"
	"gappkit/compare"
	"io"
)

type Column interface {
	Compare(x, y interface{}) compare.Order
	Decode(in *bufio.Reader) (interface{}, error)
	Encode(val interface{}, out io.Writer) error
	Name() string
}

type BasicColumn struct {
	name string
}

func (self *BasicColumn) Init(name string) *BasicColumn {
	self.name = name
	return self
}

func (self *BasicColumn) Name() string {
	return self.name
}
