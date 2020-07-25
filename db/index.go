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

type Index struct {
	file *os.File
	keyColumns []Column
	len int64
	name string
	items []IndexItem
	root *Root
	unique bool
}

type IndexItem struct {
	key IndexKey
	value RecordId
}

type IndexKey []interface{}

type IndexIter struct {
	index *Index
	i int
}

func (self *Index) Init(root *Root, name string, unique bool, keyColumns...Column) {
	self.root = root
	self.name = name
	self.unique = unique
	self.keyColumns = keyColumns	
}

func (self *Index) loadKey(in *bufio.Reader) (IndexKey, error) {
	k := make(IndexKey, len(self.keyColumns))
	
	for i, c := range self.keyColumns {
		v, err := c.Decode(in)

		if err != nil {
			if err == io.EOF {
				return nil, err
			}

			return nil, errors.Wrap(err, "Failed decoding key")
		}

		k[i] = v
	}

	return k, nil
}

func (self *Index) loadValue(in *bufio.Reader) (RecordId, error) {
	v, err := DecodeRecordId(in)
	
	if err != nil {
		return -1, errors.Wrap(err, "Failed decoding value")
	}

	return v, nil
}

func (self *Index) Open(maxTime time.Time) error {
	var err error
	path := path.Join(self.root.path, fmt.Sprintf("%v.idx", self.name))
	
	if self.file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		return err
	}

	reader := bufio.NewReader(self.file)
	
	for {
		var t time.Time

		if t, err = DecodeTime(reader); err != nil {
			if err == io.EOF {
				break
			}
			
			return errors.Wrap(err, "Failed decoding timestamp")
		}

		if !t.Before(maxTime) {
			self.file.Seek(0, 2)
			break;
		}
		
		var key IndexKey		
		key, err := self.loadKey(reader)

		if err != nil {
			if err == io.EOF {
				break
			}
			
			return err
		}

		var val RecordId
		
		if val, err = self.loadValue(reader); err != nil {
			return err
		}

		if val < 0 {
			self.remove(key, -val)
		} else {
			self.add(key, val)
		}
	}

	return nil
}

func (self *Index) Close() error {
	if err := self.file.Close(); err != nil {
		return errors.Wrap(err, "Failed closing file")
	}

	return nil
}

func (self *Index) Key(rec Record) IndexKey {
	k := make(IndexKey, len(self.keyColumns))
	
	for i, c := range self.keyColumns {
		k[i] = rec.Get(c)
	}

	return k
}

func (self *Index) Compare(x, y IndexKey) compare.Order {
	for i, xv := range x {
		if result := self.keyColumns[i].Compare(xv, y[i]); result != compare.Eq {
			return result
		}
	}
	
	return compare.Eq
}

func (self *Index) storeKey(key IndexKey) error {
	for i, c := range self.keyColumns {
		if err := c.Encode(key[i], self.file); err != nil {
			return errors.Wrap(err, "Failed encoding key")
		}
	}

	return nil
}

func (self *Index) storeValue(val RecordId) error {
	if err := EncodeRecordId(val, self.file); err != nil {
		return errors.Wrap(err, "Failed encoding value")
	}

	return nil
}

func (self *Index) add(key IndexKey, value RecordId) bool {
	i, ok := self.find(key)

	if ok && self.unique {
		return false
	}

	it := IndexItem{key: key, value: value}
	len := self.Len()
	
	if i == len {
		self.items = append(self.items, it)
	} else if i == len-1 {
		self.items = append(self.items[:i], it, self.items[i])
	} else {
		self.items = append(self.items, it)
		copy(self.items[i+1:], self.items[i:])
		self.items[i] = it
	}

	return true
}

func (self *Index) Add(rec Record) (bool, error) {
	key := self.Key(rec)
	
	if !self.add(key, rec.id) {
		return false, nil
	}

	if err := EncodeTime(time.Now(), self.file); err != nil {
		return false, errors.Wrap(err, "Failed encoding timestamp")
	}
	
	if err := self.storeKey(key); err != nil {
		return false, err
	}
	
	if err := self.storeValue(rec.id); err != nil {
		return false, err
	}
	
	return true, nil
}

func (self Index) find(key IndexKey) (int, bool) {
	min, max := 0, self.Len()

	for min < max {
		i := (min+max) / 2
		
		switch self.Compare(key, self.items[i].key) {		
			case compare.Lt:
				max = i
			case compare.Gt:
				min = i+1
			default:
				return i, true
		}
	}

	return max, false
}

func (self *Index) Find(key...interface{}) RecordId {
	i, ok := self.find(key)

	if !ok {
		return -1
	}

	return self.items[i].value
}

func (self *Index) First(key...interface{}) *IndexIter {
	return self.NewIter(0)
}

func (self *Index) FindLower(key...interface{}) *IndexIter {
	if len(self.items) == 0 {
		return self.NewIter(0)
	}
	
	i, _ := self.find(key)

	if i > 0 && i == len(self.items) {
		i--
	}

	for i > 0 && self.Compare(key, self.items[i].key) == compare.Lt {
		i--
	}

	return self.NewIter(i)
}

func (self *Index) Len() int {
	return len(self.items)
}

func (self *Index) Remove(record Record) (bool, error) {
	key := self.Key(record)
	
	if !self.remove(key, record.id) {
		return false, nil
	}
	
	if err := self.storeKey(key); err != nil {
		return false, err
	}
	
	if err := self.storeValue(-record.id); err != nil {
		return false, err
	}
	
	return true, nil
}

func (self *Index) remove(key IndexKey, value RecordId) bool {
	i, ok := self.find(key)	

	if !ok {
		return false
	}

	j := i
	
	for i > 0 && self.Compare(key, self.items[i-1].key) == compare.Eq {
		i--
	}

	len := self.Len()

	for j+1 < len && self.Compare(key, self.items[j+1].key) == compare.Eq {
		j++
	}

	for k := i; k <= j; k++ {
		if self.items[k].value == value {
			copy(self.items[k:], self.items[k+1:])
			self.items = self.items[:len-1]	
			return true
		}
	}

	return false
}

func (self *Index) NewIter(i int) *IndexIter {
	return &IndexIter{index: self, i: i}
}

func (self *IndexIter) Key(j int) interface{} {
	return self.index.items[self.i].key[j]
}

func (self *IndexIter) Prev() {
	self.i--
}

func (self *IndexIter) Next() {
	self.i++
}

func (self *IndexIter) Valid() bool {
	return self.i < len(self.index.items)
}

func (self *IndexIter) Value() RecordId {
	return self.index.items[self.i].value
}
