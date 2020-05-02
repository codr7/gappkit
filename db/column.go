package db

import (
	"bufio"
	"gappkit/compare"
	"fmt"
	"io"
	"reflect"
)

type Column interface {
	Compare(x, y interface{}) compare.Order
	Decode(in *bufio.Reader) (interface{}, error)
	Encode(val interface{}, out io.Writer) error
	Get(record Record) interface{}
	Name() string
	Set(record Record, val interface{})
}

type BasicColumn struct {
	name string
}

func (self *BasicColumn) Init(name string) *BasicColumn {
	self.name = name
	return self
}

func (self *BasicColumn) Name() string {
	return self.name
}

func (self *BasicColumn) Get(in Record) interface{} {
	s := reflect.ValueOf(in)

	if !s.IsValid() {
		panic(fmt.Sprintf("Failed getting column '%v' from nil", self.name))
	}

	f := reflect.Indirect(s).FieldByName(self.name)

	if !f.IsValid() {
		panic(fmt.Sprintf("Field '%v' not found in %v", self.name, in))
	}

	return f.Interface()
}

func (self *BasicColumn) Set(out Record, val interface{}) {
	s := reflect.ValueOf(out)

	if !s.IsValid() {
		panic(fmt.Sprintf("Failed getting column '%v' from nil", self.name))
	}

	f := reflect.Indirect(s).FieldByName(self.name)

	if !f.IsValid() {
		panic(fmt.Sprintf("Field '%v' not found in %v", self.name, out))
	}

	if !f.CanSet() {
		panic(fmt.Sprintf("Field '%v' not settable in %v", self.name, out))
	}

	v := reflect.ValueOf(val)

	if !v.IsValid() {
		panic(fmt.Sprintf("Failed reflecting value for column '%v': %v", self.name, val))
	}

	f.Set(v)
}
