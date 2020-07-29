package db

import (
	"bufio"
	"fmt"
	"gappkit/compare"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"time"
)

type RecordId = int64
type Fields = map[string]interface{}
type Offset = int64

type Table struct {
	columns map[string]Column
	dataFile *os.File
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
	self.columns = make(map[string]Column)
	self.records = make(map[RecordId]Offset)
	self.root.addTable(self)
	return self
}

func (self *Table) AddColumn(column Column) {
	self.columns[column.Name()] = column
}

func (self *Table) AddIndex(index *Index) {
	self.indexes = append(self.indexes, index)
}

func (self *Table) FindColumn(name string) Column {
	return self.columns[name]
}

func (self *Table) Name() string {
	return self.name
}

func (self *Table) Open(maxTime time.Time) error {
	var err error
	keyPath := path.Join(self.root.path, fmt.Sprintf("%v.key", self.name))
	
	if self.keyFile, err = os.OpenFile(keyPath, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		return errors.Wrapf(err, "Failed opening file: %v", keyPath)
	}

	keyReader := bufio.NewReader(self.keyFile)
	
	for {
		var t time.Time
		
		if t, err = DecodeTime(keyReader); err != nil {
			if err == io.EOF {
				break
			}
			
			return errors.Wrap(err, "Failed decoding timestamp")
		}

		if !t.Before(maxTime) {
			self.keyFile.Seek(0, 2)
			break;
		}
		
		var id RecordId
		var err error
		
		if id, err = DecodeInt(keyReader); err != nil {
			if err == io.EOF {
				break
			}
			
			return errors.Wrap(err, "Failed decoding id")
		}
		
		if id > self.nextRecordId {
			self.nextRecordId = id
		}
		
		var offset Offset
		
		if offset, err = DecodeInt(keyReader); err != nil {
			return errors.Wrap(err, "Failed decoding offset")
		}

		self.records[id] = offset
	}

	dataPath := path.Join(self.root.path, fmt.Sprintf("%v.dat", self.name))

	if self.dataFile, err = os.OpenFile(dataPath, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		return errors.Wrapf(err, "Failed opening file: %v", dataPath)
	}

	for _, index := range self.indexes {
		if err = index.Open(maxTime); err != nil {
			return errors.Wrapf(err, "Failed opening index: %v", index.name)
		}
	}
	
	return nil
}

func (self *Table) Close() error {
	if err := self.keyFile.Close(); err != nil {
		return errors.Wrap(err, "Failed closing key file")
	}

	if err := self.dataFile.Close(); err != nil {
		return errors.Wrap(err, "Failed closing data file")
	}

	for _, index := range self.indexes {
		if err := index.Close(); err != nil {
			return errors.Wrapf(err, "Failed closing index: %v", index.name)
		}
	}

	return nil
}

func (self *Table) Len() int {
	return len(self.records)
}

func (self *Table) NextId() RecordId {
	self.nextRecordId++
	return self.nextRecordId
}

func (self *Table) storeKey(id RecordId, offset Offset) error {
	if err := EncodeTime(time.Now(), self.keyFile); err != nil {
		return errors.Wrap(err, "Failed encoding timestamp")
	}

	if err := EncodeInt(id, self.keyFile); err != nil {
		return errors.Wrap(err, "Failed encoding id")
	}

	if err := EncodeInt(offset, self.keyFile); err != nil {
		return errors.Wrap(err, "Failed encoding offset")
	}

	return nil
}

func (self *Table) storeData(rec Record) (Offset, error) {	
	if prev, err := self.Load(rec.id); err != nil {
		return -1, errors.Wrap(err, "Failed loading record")
	} else if prev != nil {
		for _, ix := range self.indexes {
			if _, err := ix.Remove(*prev); err != nil {
				return -1, errors.Wrapf(err, "Failed removing from index: %v", ix.name)
			}
		}
	}
	
	offset, err := self.dataFile.Seek(0, io.SeekEnd)

	if err != nil {
		return -1, errors.Wrap(err, "Failed seeking data file")
	}

	self.records[rec.id] = offset
	
	if err := EncodeInt(int64(rec.Len()), self.dataFile); err != nil {
		return -1, errors.Wrap(err, "Failed encoding record length")
	}
	
	for _, f := range rec.Fields {
		if err := EncodeString(f.Column.Name(), self.dataFile); err != nil {
			return -1, errors.Wrap(err, "Failed encoding field name")
		}

		if err := f.Column.Encode(f.Value, self.dataFile); err != nil {
			return -1, errors.Wrap(err, "Failed encoding field value")
		}		
	}

	for _, ix := range self.indexes {
		if ok, err := ix.Add(rec); err != nil {
			return -1, errors.Wrapf(err, "Failed adding to index: %v", ix.name)
		} else if !ok {
			return -1, errors.New(fmt.Sprintf("Duplicate key in index: %v\n%v", ix.name, rec))
		}
	}
	
	return offset, nil
}

func (self *Table) Store(rec Record) error {
	offset, err := self.storeData(rec)

	if err != nil {
		return errors.Wrap(err, "Failed storing data")
	}


	if err = self.storeKey(rec.id, offset); err != nil {
		return errors.Wrap(err, "Falied storing key")
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
		return nil, errors.Wrap(err, "Failed seeking data file")
	}

	var n int
	var err error

	dataReader := bufio.NewReader(self.dataFile)
	
	if n, err = DecodeLen(dataReader); err != nil {
		return nil, errors.Wrap(err, "Failed decoding record")
	}

	out := NewRecord(id)
	out.Fields = make([]Field, n)

	for i := 0; i < n; i++ {
		f := &out.Fields[i]
		var c string
		
		if c, err = DecodeString(dataReader); err != nil {
			return nil, errors.Wrap(err, "Failed decoding field name")
		}
		
		f.Column = self.FindColumn(c)

		if f.Value, err = f.Column.Decode(dataReader); err != nil {
			return nil, errors.Wrapf(err, "Failed decoding field value: %v", c)			
		}
	}

	return out, nil
}

func (self *Table) CopyFromModel(source Model, target *Record) error {
	for _, c := range self.columns {
		v, err := Get(source, c.Name())

		if err != nil {
			return err
		}
		
		target.Set(c, v)
	}

	return nil
}

func (self *Table) CopyToModel(source Record, target Model) error {
	for _, c := range self.columns {
		if err := Set(target, c.Name(), source.Get(c)); err != nil {
			return err
		}
	}

	return nil
}

func CompareRecordId(x, y RecordId) compare.Order {
	if x < y {
		return compare.Lt
	}

	if x > y {
		return compare.Gt
	}
	
	return compare.Eq
}
