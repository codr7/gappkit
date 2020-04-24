package db

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
)

type RecordId = uint64
type Fields = map[string]interface{}
type Offset = int64

type Table interface {
	AddColumn(Column)
	AddIndex(*Index)
	Close() error
	Exists(RecordId) bool
	Load(RecordId) (Record, error)
	LoadRecord(RecordId, Record) error
	NextId() RecordId
	Open() error
	Store(Record) error
}

type BasicTable struct {
	columns []Column
	dataFile *os.File
	keyEncoder *gob.Encoder
	keyFile *os.File
	name string
	nextRecordId RecordId
	root *Root
	records map[RecordId]Offset
	indexes []*Index
}

func (self *BasicTable) Init(root *Root, name string) *BasicTable {
	self.root = root
	self.name = name
	self.records = make(map[RecordId]Offset)
	return self
}

func (self *BasicTable) AddColumn(column Column) {
	self.columns = append(self.columns, column)
}

func (self *BasicTable) AddIndex(index *Index) {
	self.indexes = append(self.indexes, index)
}

func (self *BasicTable) Open() error {
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

	for _, index := range self.indexes {
		if err = index.Open(); err != nil {
			return err
		}
	}
	
	return nil
}

func (self *BasicTable) Close() error {
	if err := self.keyFile.Close(); err != nil {
		return err
	}

	if err := self.dataFile.Close(); err != nil {
		return err
	}

	for _, index := range self.indexes {
		if err := index.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (self *BasicTable) NextId() RecordId {
	self.nextRecordId++
	return self.nextRecordId
}

func (self *BasicTable) storeKey(id RecordId, offset Offset) error {
	if err := self.keyEncoder.Encode(id); err != nil {
		return err
	}

	if err := self.keyEncoder.Encode(offset); err != nil {
		return err
	}

	return nil
}

func (self *BasicTable) storeData(record Record) (Offset, error) {
	fields := make(Fields)

	for _, c := range self.columns {
		fields[c.Name()] = c.Get(record)
	}

	offset, err := self.dataFile.Seek(0, io.SeekEnd)
	
	if err != nil {
		return -1, err
	}

	id := record.Id()

	if prev, err := self.loadFields(id); err != nil {
		return -1, err
	} else if prev != nil {
		for _, ix := range self.indexes {
			ix.Remove(prev, id)
		}
	}
	
	self.records[id] = offset
	encoder := gob.NewEncoder(self.dataFile)
	
	if err := encoder.Encode(fields); err != nil {
		return -1, err
	}

	for _, ix := range self.indexes {
		ix.Add(fields, id)
	}
	
	return offset, nil
}

func (self *BasicTable) Store(record Record) error {
	offset, err := self.storeData(record)

	if err != nil {
		return err
	}

	if err = self.storeKey(record.Id(), offset); err != nil {
		return err
	}

	return nil
}

func (self *BasicTable) Exists(id RecordId) bool {
	_, ok := self.records[id]
	return ok
}

func (self *BasicTable) loadFields(id RecordId) (Fields, error) {
	offset, ok := self.records[id]

	if !ok {
		return nil, nil
	}
	
	if _, err := self.dataFile.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	
	decoder := gob.NewDecoder(self.dataFile)
	fields := make(Fields)
	
	if err := decoder.Decode(&fields); err != nil {
		return nil, fmt.Errorf("Failed decoding record: %v", err)
	}

	return fields, nil
}

func (self *BasicTable) LoadRecord(id RecordId, record Record) error {
	fields, err := self.loadFields(id)

	if err != nil {
		return err
	}

	if fields == nil {
		return fmt.Errorf("Record not found: %v/%v", self.name, id)
	}

	for _, c := range self.columns {
		c.Set(record, fields[c.Name()])
	}

	return nil
}
