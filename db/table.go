package db

import (
	"encoding/gob"
	"fmt"
	"gappkit/util"
	"io"
	"os"
	"path"
)

type RecordId = uint64
type Fields = map[string]interface{}
type Offset = int64

type Table struct {
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

func (self *Table) Init(root *Root, name string) *Table {
	self.root = root
	self.name = name
	self.records = make(map[RecordId]Offset)
	return self
}

func (self *Table) AddColumn(column Column) {
	self.columns = append(self.columns, column)
}

func (self *Table) AddIndex(index *Index) {
	self.indexes = append(self.indexes, index)
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
	dataPath := path.Join(self.root.path, fmt.Sprintf("%v.dat", self.name))

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

func (self *Table) Close() error {
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

func (self *Table) NextId() RecordId {
	self.nextRecordId++
	return self.nextRecordId
}

func (self *Table) storeKey(id RecordId, offset Offset) error {
	if err := self.keyEncoder.Encode(id); err != nil {
		return err
	}

	if err := self.keyEncoder.Encode(offset); err != nil {
		return err
	}

	return nil
}

func (self *Table) storeData(id RecordId, record Record) (Offset, error) {
	offset, err := self.dataFile.Seek(0, io.SeekEnd)
	
	if err != nil {
		return -1, err
	}

	if prev, err := self.Load(id); err != nil {
		return -1, err
	} else if prev != nil {
		for _, ix := range self.indexes {
			ix.Remove(id, *prev)
		}
	}
	
	self.records[id] = offset
	encoder := gob.NewEncoder(self.dataFile)
	
	if err := encoder.Encode(record); err != nil {
		return -1, err
	}

	for _, ix := range self.indexes {
		ix.Add(id, record)
	}
	
	return offset, nil
}

func (self *Table) Store(id RecordId, record Record) error {
	offset, err := self.storeData(id, record)

	if err != nil {
		return err
	}

	if err = self.storeKey(id, offset); err != nil {
		return err
	}

	return nil
}

func (self *Table) Exists(id RecordId) bool {
	_, ok := self.records[id]
	return ok
}

func (self *Table) Load(id RecordId) (*Record, error) {
	offset, ok := self.records[id]

	if !ok {
		return nil, nil
	}
	
	if _, err := self.dataFile.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	
	decoder := gob.NewDecoder(self.dataFile)
	record := new(Record)
	
	if err := decoder.Decode(&record); err != nil {
		return nil, fmt.Errorf("Failed decoding record: %v", err)
	}

	return record, nil
}

func CompareRecordId(x, y RecordId) util.Order {
	if x < y {
		return util.Lt
	}

	if x > y {
		return util.Gt
	}
	
	return util.Eq
}
