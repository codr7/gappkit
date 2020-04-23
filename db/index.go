package db

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
)

type Index struct {
	root *Root
	name string
	unique bool
	compare IndexCompare
	file *os.File
	records *IndexNode
	len uint64
}

type IndexAction = func(key IndexKey, value RecordId) error
type IndexCompare = func(x, y IndexKey) Order
type IndexKey []interface{}
type IndexValue []RecordId

type IndexNode struct {
	left, right *IndexNode
	key IndexKey
	value IndexValue
	red bool
}

type IndexIter struct {
	node *IndexNode
	stack []*IndexNode
}

func (self *Index) Init(root *Root, name string, unique bool, keyColumns...Column) {
	self.root = root
	self.name = name
	self.unique = unique
	
	self.compare = func(x, y IndexKey) Order {
		for i, c := range keyColumns {
			if result := c.Compare(x[i], y[i]); result != Eq {
				return result
			}
		}

		return Eq
	}
}

func (self *Index) Open() error {
	var err error
	path := path.Join(self.root.path, fmt.Sprintf("%v.index", self.name))
	
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

func (self *Index) AddNode(key IndexKey, value RecordId) bool {
	var ok bool
	self.records, ok = self.addNode(self.records, key, value)
	self.records.red = false
	return ok
}

func (self *Index) Do(action IndexAction) error {
	return self.records.Do(action)
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

func (self *Index) RemoveNode(key IndexKey) IndexValue {
	var out IndexValue
	self.records, out = self.removeNode(self.records, key)
	self.records.red = false
	return out
}

func (self *Index) removeNode(node *IndexNode, key IndexKey) (*IndexNode, IndexValue) {
	var out IndexValue

	if self.compare(key, node.key) == Lt {
		if !node.left.isRed() && !node.left.left.isRed() {
			node = node.moveRedLeft()
		}

		node.left, out = self.removeNode(node.left, key)
	} else {
		if node.left.isRed() {
			node = node.rotr()
		}

		if self.compare(key, node.key) == Eq && node.right == nil {
			out = node.value
			self.len--
			return nil, out
		}

		if !node.right.isRed() && !node.right.left.isRed() {
			node.flip()

			if node.left.left.isRed() {
				node = node.rotr()
				node.flip()
			}
		}

		if self.compare(key, node.key) == Eq {
			l := node.left
			var r *IndexNode
			r, node = node.right.removeMin()
			node.left = l
			node.right = r
			self.len--
		} else {
			node.right, out = self.removeNode(node.right, key)
		}
	}
	
	return node.fix(), out
}

func (self *Index) addNode(node *IndexNode, key IndexKey, value RecordId) (*IndexNode, bool) {
	if node == nil {
		node = &IndexNode{key: key, value: IndexValue{value}, red: true}
		self.len++
		return node, true
	}

	var ok bool

	switch self.compare(key, node.key) {
	case Lt:
		node.left, ok = self.addNode(node.left, key, value)
	case Gt:
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
		case Lt:
			node = node.left
		case Gt:
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
		case Lt:
			node = node.left
		default:
			return node
		}
	}

	return nil
}

func (self *IndexNode) Do(action IndexAction) error {
	if self != nil {
		if err := self.left.Do(action); err != nil {
			return err
		}
		
		for _, v := range self.value {
			if err := action(self.key, v); err != nil {
				return err
			}
		}
		
		if err := self.right.Do(action); err != nil {
			return err
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
	self := new(IndexIter)

	if node !=nil {
		self.Append(node)
	}

	return self
}

func (self *IndexIter) Append(node *IndexNode) {
	if node.right != nil {
		self.Append(node.right)
	}

	self.stack = append(self.stack, node)

	if node.left != nil {
		self.Append(node.left)
	}
}

func (self *IndexIter) Key() IndexKey {
	return self.node.key
}

func (self *IndexIter) Next() bool {
	i := len(self.stack)-1

	if i == -1 {
		return false
	}

	self.node, self.stack = self.stack[i], self.stack[:i]
	return true
}

func (self *IndexIter) Value() IndexValue {
	return self.node.value
}
