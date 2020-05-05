package db

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

func Get(source Model, fieldName string) (interface{}, error) {
	s := reflect.ValueOf(source)

	if !s.IsValid() {
		return nil, errors.New("Failed reflecting source")
	}

	f := reflect.Indirect(s).FieldByName(fieldName)

	if !f.IsValid() {
		return nil, errors.New(fmt.Sprintf("Field not found: %v", fieldName))
	}

	return f.Interface(), nil
}

func Set(target Model, fieldName string, val interface{}) error {
	t := reflect.ValueOf(target)

	if !t.IsValid() {
		return errors.New("Failed reflecting target")
	}

	f := reflect.Indirect(t).FieldByName(fieldName)

	if !f.IsValid() {
		return errors.New(fmt.Sprintf("Field not found: %v", fieldName))
	}

	if !f.CanSet() {
		return errors.New(fmt.Sprintf("Field not settable: %v", fieldName))
	}

	v := reflect.ValueOf(val)

	if !v.IsValid() {
		return errors.New(fmt.Sprintf("Failed reflecting value: %v", val))
	}

	f.Set(v)
	return nil
}
