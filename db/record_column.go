package db

import (
	"bufio"
	"gappkit/compare"
	"io"
)

type RecordColumn struct {
	BasicColumn
}

func (self *RecordColumn) Init(name string) *RecordColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *RecordColumn) Compare(x, y interface{}) compare.Order {
	return CompareRecordId(x.(RecordId), y.(RecordId))
}

func (self *RecordColumn) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeRecordId(in)
}

func (self *RecordColumn) Encode(val interface{}, out io.Writer) error {
	return EncodeRecordId(val.(RecordId), out)
}
