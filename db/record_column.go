package db

import (
	"bufio"
	"gappkit/compare"
	"io"
	"reflect"
)

var (
	RecordType RecordColumnType
	recordValueType = reflect.TypeOf(RecordId(0))
)

type RecordColumnType struct {
}

func (self *RecordColumnType) Compare(x, y interface{}) compare.Order {
	return CompareRecordId(x.(RecordId), y.(RecordId))
}

func (self *RecordColumnType) Decode(in *bufio.Reader) (interface{}, error) {
	return DecodeRecordId(in)
}

func (self *RecordColumnType) Encode(val interface{}, out io.Writer) error {
	return EncodeRecordId(val.(RecordId), out)
}

func (self *RecordColumnType) ValueType() reflect.Type {
	return recordValueType
}

type RecordColumn struct {
	BasicColumn
}

func (self *RecordColumn) Init(name string) *RecordColumn {
	self.BasicColumn.Init(name, &RecordType)
	return self
}
