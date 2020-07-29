package db

import (
	"github.com/pkg/errors"
	"os"
	"time"
)

type Root struct {
	path string
	tables []*Table
	indexes []*Index
}

func NewRoot(path string) *Root {
	return new(Root).Init(path)
}

func (self *Root) Init(path string) *Root {
	self.path = path
	return self
}

func (self *Root) Drop() error {
	if err := os.RemoveAll(self.path); err != nil {
		return errors.Wrap(err, "Failed dropping database") 
	}

	return nil
}

func (self *Root) Open(maxTime time.Time) error {
	if _, err := os.Stat(self.path); os.IsNotExist(err) {
		if err = os.Mkdir(self.path, os.ModePerm); err != nil {
			return err
		}
	}

	for _, i := range self.indexes {
		if err := i.Open(maxTime); err != nil {
			return err
		}
	}

	for _, t := range self.tables {
		if err := t.Open(maxTime); err != nil {
			return err
		}
	}

	return nil
}

func (self *Root) Close() error {
	for _, t := range self.tables {
		if err := t.Close(); err != nil {
			return err
		}
	}

	for _, i := range self.indexes {
		if err := i.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (self *Root) addTable(table *Table) {
	self.tables = append(self.tables, table)
}

func (self *Root) addIndex(index *Index) {
	self.indexes = append(self.indexes, index)
}
