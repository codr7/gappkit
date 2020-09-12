package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"reflect"
)

type Slice = []interface{}

type SliceColumnType struct {
	itemType ColumnType
	valueType reflect.Type
}

func SliceType(itemType ColumnType) *SliceColumnType {
	return &SliceColumnType{itemType: itemType, valueType: reflect.SliceOf(itemType.ValueType())}
}

func (self *SliceColumnType) Compare(x, y interface{}) compare.Order {
	xs, ys := reflect.ValueOf(x), reflect.ValueOf(y)
	max := ys.Len()-1
	var i int
	
	for i := 0; i < xs.Len(); i++ {
		if i > max {
			return compare.Gt
		}

		if result := self.itemType.Compare(xs.Index(i).Interface(), ys.Index(i).Interface());
		result != compare.Eq {
			return result
		}
	}

	if max > i {
		return compare.Lt
	}

	return compare.Eq
}

func (self *SliceColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	n, err := DecodeLen(in)

	if err != nil {
		return nil, err
	}

	s := reflect.MakeSlice(self.ValueType(), n, n)
	
	for i := 0; i < n; i++ {		
		if dv, err := self.itemType.Decode(in); err == nil {
			v := s.Index(i)
			v.Set(reflect.ValueOf(dv))
			continue
		}

		return nil, err
	}

	return s.Interface(), nil
}

func (self *SliceColumnType) Encode(val interface{}, out io.Writer) error {
	/*if val == nil {
		if err := EncodeInt(0, out); err != nil {
			return err
		}
	}*/
	
	s := reflect.ValueOf(val)
	n := s.Len()
	
	if err := EncodeInt(int64(n), out); err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if err := self.itemType.Encode(s.Index(i).Interface(), out); err != nil {
			return err
		}
	}

	return nil
}

func (self *SliceColumnType) ValueType() reflect.Type {
	return self.valueType
}

type SliceColumn struct {
	BasicColumn
}

func (self *SliceColumn) Init(name string, itemType ColumnType) *SliceColumn {
	self.BasicColumn.Init(name, SliceType(itemType))
	return self
}
