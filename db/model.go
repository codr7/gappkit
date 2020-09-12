package db

import (
	"fmt"
	"gappkit/compare"
	"github.com/pkg/errors"
	"log"
	"reflect"
)

type Model interface {
	Id() RecordId
	Store() error
	Table() *Table
}

type BasicModel struct {
	id RecordId
}

func (self *BasicModel) Init(id RecordId) {
	self.id = id
}

func (self *BasicModel) Id() RecordId {
	return self.id
}

func Get(source Model, fieldName string) interface{} {
	s := reflect.ValueOf(source)

	if !s.IsValid() {
		log.Fatal("Failed reflecting source")
	}

	f := reflect.Indirect(s).FieldByName(fieldName)

	if !f.IsValid() {
		log.Fatalf("Field not found: %v", fieldName)
	}

	return f.Interface()
}

func Set(target Model, fieldName string, val interface{}) {
	t := reflect.ValueOf(target)

	if !t.IsValid() {
		log.Fatal("Failed reflecting target")
	}

	f := reflect.Indirect(t).FieldByName(fieldName)

	if !f.IsValid() {
		log.Fatalf(fmt.Sprintf("Field not found: %v", fieldName))
	}

	if !f.CanSet() {
		log.Fatalf(fmt.Sprintf("Field not settable: %v", fieldName))
	}

	v := reflect.ValueOf(val)

	if !v.IsValid() {
		log.Fatalf(fmt.Sprintf("Failed reflecting value: %v", val))
	}

	f.Set(v)
}

func Load(out Model) error {
	t := out.Table()
	id := out.Id()
	
	in, err := t.Load(id)
	
	if err != nil || in == nil {
		return errors.Wrapf(err, "Failed loading model: %v/%v", t.Name(), id)
	}

	t.CopyToModel(*in, out)
	return nil
}

func Modified(in Model) bool {
	t := in.Table()
	prev, err := t.Load(in.Id())

	if err != nil {
		log.Fatal(err)
	}

	if prev == nil {
		return true
	}
	
	for _, c := range t.columns {
		if cmp := c.Compare(Get(in, c.Name()), prev.Get(c)); cmp != compare.Eq {
			return true
		}
	}

	return false
}

func Store(in Model) (Record, error) {
	out := ToRecord(in)

	if err := in.Table().Store(out); err != nil {
		return out, err
	}

	return out, nil
}

func ToRecord(in Model) Record {
	var out Record
	out.Init(in.Id())
	in.Table().CopyFromModel(in, &out)
	return out
}
