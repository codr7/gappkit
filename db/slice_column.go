package db

import (
	"bufio"
	"gappkit/compare"
	"io"
)

type Slice = []interface{}

type SliceColumnType struct {
	itemType ColumnType
}

func SliceType(itemType ColumnType) *SliceColumnType {
	return &SliceColumnType{itemType: itemType}
}

func (self *SliceColumnType) Compare(x, y interface{}) compare.Order {
	xs, ys := x.(Slice), y.(Slice)
	max := len(ys)-1
	var i int
	var v interface{}
	
	for i, v = range xs {
		if i > max {
			return compare.Gt
		}

		if result := self.itemType.Compare(v, ys[i]); result != compare.Eq {
			return result
		}
	}

	if max > i {
		return compare.Lt
	}

	return compare.Eq
}

func (self *SliceColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	l, err := DecodeLen(in)

	if err != nil {
		return nil, err
	}

	s :=  make(Slice, l)
	
	for i := 0; i < l; i++ {
		if s[i], err = self.itemType.Decode(in); err != nil {
			return nil, err
		}
	}

	return s, nil
}
	
func (self *SliceColumnType) Encode(val interface{}, out io.Writer) error {
	s := val.(Slice)
	
	if err := EncodeInt(int64(len(s)), out); err != nil {
		return err
	}

	for _, v := range s {
		if err := self.itemType.Encode(v, out); err != nil {
			return err
		}
	}

	return nil
}

type SliceColumn struct {
	BasicColumn
}

func (self *SliceColumn) Init(name string, itemType ColumnType) *SliceColumn {
	self.BasicColumn.Init(name, SliceType(itemType))
	return self
}
