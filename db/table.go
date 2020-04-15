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
type Offset = int64

type Table struct {
	root *Root
	name string
	columns []Column
	nextRecordId RecordId
	records map[RecordId]Offset
	file *os.File
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

func (self *Table) NewColumn(name string) Column {
	c := new(ColumnBase).Init(name)
	self.AddColumn(c)
	return c
}

func (self *Table) Open() error {	
	p := path.Join(self.root.path, fmt.Sprintf("%v.table", self.name))
	var err error

	if self.file, err = os.OpenFile(p, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		return err
	}

	var offset Offset
	
	for {
		if _, err := self.file.Seek(offset, io.SeekStart); err != nil {
			return err
		}
		
		decoder := gob.NewDecoder(self.file)
		var id RecordId

		if err := decoder.Decode(&id); err != nil {
			if err == io.EOF {
				fmt.Printf("EOF\n")
				break
			}

			return fmt.Errorf("Failed decoding id: %v", err)
		}

		fmt.Printf("id: %v\n", id)
		
		if id > self.nextRecordId {
			self.nextRecordId = id
		}
		
		var size RecordSize

		if err := decoder.Decode(&size); err != nil {
			return fmt.Errorf("Failed decoding size: %v", err)
		}

		fmt.Printf("size: %v\n", size)

		self.records[id], _ = self.file.Seek(0, io.SeekCurrent)

		if offset, err := self.file.Seek(int64(size), io.SeekCurrent); err != nil {
			return err
		}
	}

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

	if _, err := self.file.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	encoder := gob.NewEncoder(self.file)
	if err := encoder.Encode(id); err != nil {
		return err
	}

	var buffer bytes.Buffer
	dataEncoder := gob.NewEncoder(&buffer)
	
	if err := dataEncoder.Encode(data); err != nil {
		return err
	}

	size := RecordSize(buffer.Len())
	
	if err := encoder.Encode(size); err != nil {
		return err
	}

	self.records[id], _ = self.file.Seek(0, io.SeekCurrent)
	
	if _, err := self.file.Write(buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (self *Table) Load(id RecordId, record interface{}) error {
	offset, ok := self.records[id]

	if !ok {
		return fmt.Errorf("Record not found: %v/%v", self.name, id)
	}

	if _, err := self.file.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	data := make(Record)
	decoder := gob.NewDecoder(self.file)
	
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("Failed decoding record: %v", err)
	}

	return nil
}
