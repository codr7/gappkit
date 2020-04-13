package db

import (
	"fmt"
	"reflect"
)

type Column interface {
	Name() string
	Clone(val interface{}) interface{}
	
	Get(record interface{}) interface{}
	Set(record interface{}, val interface{})
}

type ColumnBase struct {
	name string
}

func (self *ColumnBase) Init(name string) *ColumnBase {
	self.name = name
	return self
}

func (self *ColumnBase) Name() string {
	return self.name
}

func (_ *ColumnBase) Clone(val interface{}) interface{} {
	return val
}

func (self *ColumnBase) Get(in interface{}) interface{} {
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

func (self *ColumnBase) Set(out interface{}, val interface{}) {
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
