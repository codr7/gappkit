package db

import (
	"bufio"
	"gappkit/compare"
	"io"
)

type RecordSetColumn struct {
	BasicColumn
}

func (self *RecordSetColumn) Init(name string) *RecordSetColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *RecordSetColumn) Compare(x, y interface{}) compare.Order {
	return x.(RecordSet).Compare(y.(RecordSet))
}

func (self *RecordSetColumn) Decode(in *bufio.Reader) (interface{}, error) {
	l, err := DecodeLen(in)

	if err != nil {
		return nil, err
	}

	var v RecordSet
	v.Items = make([]RecordId, l)
	
	for i := 0; i < l; i++ {
		if v.Items[i], err = DecodeRecordId(in); err != nil {
			return nil, err
		}
	}

	return v, nil
}
	
func (self *RecordSetColumn) Encode(val interface{}, out io.Writer) error {
	v := val.(RecordSet)
	
	if err := EncodeInt(int64(v.Len()), out); err != nil {
		return err
	}

	for _, id := range v.Items {
		if err := EncodeRecordId(id, out); err != nil {
			return err
		}
	}

	return nil
}

