package db

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
)

type Record = map[string]interface{}
type RecordId = uint64
type RecordOffsetMap = map[RecordId]int64
	
type Table struct {
	root *Root
	name string
	columns []Column
	nextRecordId RecordId
	recordOffsets RecordOffsetMap
	reader,
	writer *os.File
	decoder *gob.Decoder
	encoder *gob.Encoder
}

func (self *Table) Init(root *Root, name string) *Table {
	self.root = root
	self.name = name
	self.recordOffsets = make(RecordOffsetMap)
	return self
}

func (self *Table) AddColumn(column Column) {
	self.columns = append(self.columns, column)
}

func (self *Table) NewColumn(name string) Column {
	c := new(ColumnBase).Init(name)
	self.AddColumn(c)
	return c
}

func (self *Table) Open() error {
	p := path.Join(self.root.path, fmt.Sprintf("%v.table", self.name))
	var err error
	
	if self.writer, err = os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err != nil {
		return err
	}

	self.encoder = gob.NewEncoder(self.writer)
	
	if self.reader, err = os.OpenFile(p, os.O_RDONLY, os.ModePerm); err != nil {
		return err
	}

	self.decoder = gob.NewDecoder(self.reader)
	return nil
}

func (self *Table) Close() error {
	if err := self.reader.Close(); err != nil {
		return err
	}

	if err := self.writer.Close(); err != nil {
		return err
	}
	
	return nil
}

func (self *Table) NextId() RecordId {
	self.nextRecordId++
	return self.nextRecordId
}

func (self *Table) Store(id RecordId, record interface{}) error {
	data := make(Record)

	for _, c := range self.columns {
		data[c.Name()] = c.Get(record)
	}

	if err := self.encoder.Encode(id); err != nil {
		return err
	}

	offset, err := self.writer.Seek(0, io.SeekCurrent)
	
	if err != nil {
		return err
	}

	self.recordOffsets[id] = offset

	if err := self.encoder.Encode(data); err != nil {
		return err
	}
	
	return nil
}
