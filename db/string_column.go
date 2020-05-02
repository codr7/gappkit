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
	l, err := DecodeLen(in)

	if err != nil {
		return nil, err
	}

	v := make([]byte, l)

	if _, err = in.Read(v); err != nil {
		return nil, err
	}
	
	return string(v), nil
}

func (self *StringColumn) Encode(val interface{}, out io.Writer) error {
	v := []byte(val.(string))
	
	if err := EncodeLen(v, out); err != nil {
		return err
	}

	_, err := out.Write(v)
	return err
}

