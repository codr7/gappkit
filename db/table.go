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
type Offset = int64

type Table struct {
	root *Root
	name string
	columns []Column
	nextRecordId RecordId
	records map[RecordId]Offset
	keyFile, dataFile *os.File
	keyEncoder *gob.Encoder
}

func (self *Table) Init(root *Root, name string) *Table {
	self.root = root
	self.name = name
	self.records = make(map[RecordId]Offset)
	root.AddTable(self)
	return self
}

func (self *Table) AddColumn(column Column) {
	self.columns = append(self.columns, column)
}

func (self *Table) NewColumn(name string) Column {
	c := new(BasicColumn).Init(name)
	self.AddColumn(c)
	return c
}

func (self *Table) Open() error {
	var err error
	keyPath := path.Join(self.root.path, fmt.Sprintf("%v.key", self.name))
	
	if self.keyFile, err = os.OpenFile(keyPath, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		return err
	}

	decoder := gob.NewDecoder(self.keyFile)
	
	for {
		var id RecordId
		
		if err := decoder.Decode(&id); err != nil {
			if err == io.EOF {
				break
			}
			
			return fmt.Errorf("Failed decoding id: %v", err)
		}
		
		if id > self.nextRecordId {
			self.nextRecordId = id
		}
		
		var offset Offset

		if err := decoder.Decode(&offset); err != nil {
			return fmt.Errorf("Failed decoding offset: %v", err)
		}

		self.records[id] = offset
	}

	self.keyEncoder = gob.NewEncoder(self.keyFile)
	dataPath := path.Join(self.root.path, fmt.Sprintf("%v.data", self.name))

	if self.dataFile, err = os.OpenFile(dataPath, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (self *Table) Close() error {
	if err := self.keyFile.Close(); err != nil {
		return err
	}

	if err := self.dataFile.Close(); err != nil {
		return err
	}

	return nil
}

func (self *Table) NextId() RecordId {
	self.nextRecordId++
	return self.nextRecordId
}

func (self *Table) StoreKey(id RecordId, offset Offset) error {
	if err := self.keyEncoder.Encode(id); err != nil {
		return err
	}

	if err := self.keyEncoder.Encode(offset); err != nil {
		return err
	}

	return nil
}

func (self *Table) StoreData(id RecordId, record interface{}) (Offset, error) {
	data := make(Record)

	for _, c := range self.columns {
		data[c.Name()] = c.Get(record)
	}

	offset, err := self.dataFile.Seek(0, io.SeekEnd)
	
	if err != nil {
		return -1, err
	}

	self.records[id] = offset
	encoder := gob.NewEncoder(self.dataFile)
	
	if err := encoder.Encode(data); err != nil {
		return -1, err
	}

	return offset, nil
}

func (self *Table) Store(id RecordId, record interface{}) error {
	offset, err := self.StoreData(id, record)

	if err != nil {
		return err
	}

	return self.StoreKey(id, offset)
}

func (self *Table) Load(id RecordId, record interface{}) error {
	offset, ok := self.records[id]

	if !ok {
		return fmt.Errorf("Record not found: %v/%v", self.name, id)
	}

	if _, err := self.dataFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	data := make(Record)
	decoder := gob.NewDecoder(self.dataFile)
	
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("Failed decoding record: %v", err)
	}

	for _, c := range self.columns {
		c.Set(record, data[c.Name()])
	}

	return nil
}
