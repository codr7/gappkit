package db

import (
	"encoding/gob"
	"fmt"
	"gappkit/compare"
	"io"
	"os"
	"path"
)

type Index struct {
	compare IndexCompare
	file *os.File
	keyColumns []Column
	len uint64
	name string
	records *IndexNode
	root *Root
	unique bool
}

type IndexCompare = func(x, y IndexKey) compare.Order
type IndexKey []interface{}
type IndexValue []RecordId

type IndexNode struct {
	left, right *IndexNode
	key IndexKey
	value IndexValue
	red bool
}

type IndexIter struct {
	current *IndexNode
	next *IndexNode
	stack []*IndexNode
	valueIndex int
}

func (self *Index) Init(root *Root, name string, unique bool, keyColumns...Column) {
	self.root = root
	self.name = name
	self.unique = unique
	
	self.compare = func(x, y IndexKey) compare.Order {
		for i, c := range keyColumns {
			if result := c.Compare(x[i], y[i]); result != compare.Eq {
				return result
			}
		}

		return compare.Eq
	}
}

func (self *Index) Open() error {
	var err error
	path := path.Join(self.root.path, fmt.Sprintf("%v.idx", self.name))
	
	if self.file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		return err
	}

	decoder := gob.NewDecoder(self.file)
	
	for {
		var key IndexKey
		
		if err := decoder.Decode(&key); err != nil {
			if err == io.EOF {
				break
			}
			
			return fmt.Errorf("Failed decoding key: %v", err)
		}

		var value IndexValue
		
		if err := decoder.Decode(&value); err != nil {
			if err == io.EOF {
				break
			}
			
			return fmt.Errorf("Failed decoding id: %v", err)
		}

		for _, v := range value {
			self.AddNode(key, v)
		}
	}

	return nil
}

func (self *Index) Close() error {
	if err := self.file.Close(); err != nil {
		return err
	}

	return nil
}

func (self *Index) Key(record Record) IndexKey {
	var k IndexKey
	
	for _, c := range self.keyColumns {
		k = append(k, record.Get(c))
	}

	return k
}

func (self *Index) Add(id RecordId, record Record) bool {
	return self.AddNode(self.Key(record), id)
}

func (self *Index) AddNode(key IndexKey, value RecordId) bool {
	var ok bool
	self.records, ok = self.addNode(self.records, key, value)
	self.records.red = false
	return ok
}

func (self *Index) Find(key IndexKey) IndexValue {
	n := self.findNode(self.records, key)
	if n == nil { return nil }
	return n.value
}

func (self *Index) FindLower(key...interface{}) *IndexIter {
	n := self.findLowerNode(self.records, IndexKey(key))
	return NewIndexIter(n)
}

func (self *Index) Len() uint64 {
	return self.len
}

func (self *Index) Remove(id RecordId, record Record) bool {
	return self.RemoveNode(self.Key(record), id)
}

func (self *Index) RemoveNode(key IndexKey, value RecordId) bool {
	var ok bool
	self.records, ok = self.removeNode(self.records, key, value)
	self.records.red = false
	return ok
}

func (self *Index) removeValue(node *IndexNode, value RecordId) bool {
	self.len--
	
	for i, v := range node.value {
		if v == value {
			node.value = append(node.value[:i], node.value[i+1:]...)
			return true
		}
	}

	return false
}

func (self *Index) removeNode(node *IndexNode, key IndexKey, value RecordId) (*IndexNode, bool) {
	var ok bool

	if self.compare(key, node.key) == compare.Lt {
		if !node.left.isRed() && !node.left.left.isRed() {
			node = node.moveRedLeft()
		}

		node.left, ok = self.removeNode(node.left, key, value)
	} else {
		if node.left.isRed() {
			node = node.rotr()
		}

		if self.compare(key, node.key) == compare.Eq && node.right == nil {
			ok := self.removeValue(node, value)

			if ok && len(node.value) == 0 {
				node = nil
			}
			
			return node, ok
		}

		if !node.right.isRed() && !node.right.left.isRed() {
			node.flip()

			if node.left.left.isRed() {
				node = node.rotr()
				node.flip()
			}
		}

		if self.compare(key, node.key) == compare.Eq {
			ok := self.removeValue(node, value)
			
			if ok && len(node.value) > 0 {
				return node, true
			}
			
			l := node.left
			var r *IndexNode
			r, node = node.right.removeMin()
			node.left = l
			node.right = r
		} else {
			node.right, ok = self.removeNode(node.right, key, value)
		}
	}
	
	return node.fix(), ok
}

func (self *Index) addNode(node *IndexNode, key IndexKey, value RecordId) (*IndexNode, bool) {
	if node == nil {
		node = &IndexNode{key: key, value: IndexValue{value}, red: true}
		self.len++
		return node, true
	}

	var ok bool

	switch self.compare(key, node.key) {
	case compare.Lt:
		node.left, ok = self.addNode(node.left, key, value)
	case compare.Gt:
		node.right, ok = self.addNode(node.right, key, value)
	default:
		if self.unique {
			return node, false
		}
		
		node.value = append(node.value, value)
		self.len++
		return node, true
	}

	return node.fix(), ok
}

func (self *Index) findNode(node *IndexNode, key IndexKey) *IndexNode {
	for node != nil {
		switch self.compare(key, node.key) {
		case compare.Lt:
			node = node.left
		case compare.Gt:
			node = node.right
		default:
			return node
		}
	}

	return nil
}

func (self *Index) findLowerNode(node *IndexNode, key IndexKey) *IndexNode {
	for node != nil {
		switch self.compare(key, node.key) {
		case compare.Lt:
			node = node.left
		default:
			return node
		}
	}

	return nil
}

func (self *IndexNode) fix() *IndexNode {
	if (self.right.isRed()) {
		self = self.rotl()
	}

	if (self.left.isRed() && self.left.left.isRed()) {
		self = self.rotr()
	}

	if (self.left.isRed() && self.right.isRed()) {
		self.flip()
	}

	return self
}

func (self *IndexNode) flip() {
	self.red = !self.red
	self.left.red = !self.left.red
	self.right.red = !self.right.red
}

func (self *IndexNode) isRed() bool {
	return self != nil && self.red
}

func (self *IndexNode) moveRedLeft() *IndexNode {
	self.flip()

	if self.right.left.isRed() {
		self.right = self.right.rotr()
		self = self.rotl()
		self.flip()
	}

	return self
}

func (self *IndexNode) removeMin() (*IndexNode, *IndexNode) {
	if self.left == nil {
		return nil, self
	}

	if !self.left.isRed() && !self.left.left.isRed() {
		self = self.moveRedLeft()
	}

	var out *IndexNode
	self.left, out = self.left.removeMin()
	return self.fix(), out
}

func (self *IndexNode) rotl() *IndexNode {
	r := self.right
	self.right = r.left
	r.left = self
	r.red = self.red
	self.red = true
	return r
}

func (self *IndexNode) rotr() *IndexNode {
	l := self.left
	self.left = l.right
	l.right = self
	l.red = self.red
	self.red = true
	return l
}

func NewIndexIter(node *IndexNode) *IndexIter {
	return &IndexIter{next: node}
}

func (self *IndexIter) Key(i int) interface{} {
	return self.current.key[i]
}

func (self *IndexIter) Next() bool {
	self.valueIndex++
	
	if self.current == nil || self.valueIndex >= len(self.current.value) {
		for self.next != nil {
			self.stack = append(self.stack, self.next)
			self.next = self.next.left
		}
		
		i := len(self.stack)-1
		
		if i == -1 {
			self.current = nil
			return false
		}
		
		var n *IndexNode
		n, self.stack = self.stack[i], self.stack[:i]
		self.next = n.right
		self.current = n
		self.valueIndex = 0
	}
	
	return true
}

func (self *IndexIter) Value() RecordId {
	return self.current.value[self.valueIndex]
}
