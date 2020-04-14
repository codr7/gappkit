package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
)

type Record = map[string]interface{}
type RecordId = uint64
type RecordSize = uint32
type Offset = uint64
type RecordOffsetMap = map[RecordId]Offset
	
type Table struct {
	root *Root
	name string
	columns []Column
	nextRecordId RecordId
	records RecordOffsetMap
	file *os.File
	dataBuffer bytes.Buffer
	encoder, dataEncoder *gob.Encoder
	decoder *gob.Decoder
}

func (self *Table) Init(root *Root, name string) *Table {
	self.root = root
	self.name = name
	self.records = make(RecordOffsetMap)
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

	if self.file, err = os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm); err != nil {
		return err
	}

	self.encoder = gob.NewEncoder(self.file)
	self.dataEncoder = gob.NewEncoder(&self.dataBuffer)
	self.decoder = gob.NewDecoder(self.file)
	return nil
}

func (self *Table) Close() error {
	if err := self.file.Close(); err != nil {
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

	offset, err := self.file.Seek(0, io.SeekEnd)
	
	if err != nil {
		return err
	}

	self.records[id] = Offset(offset)

	if err := self.dataEncoder.Encode(data); err != nil {
		return err
	}

	size := RecordSize(self.dataBuffer.Len())
	
	if err := self.encoder.Encode(size); err != nil {
		return err
	}

	if _, err := self.file.Write(self.dataBuffer.Bytes()); err != nil {
		return err
	}

	self.dataBuffer.Reset()
	return nil
}
