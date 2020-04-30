package db

import (
	"gappkit/util"
)

type Field struct {
	Column Column
	Value interface{}
}

type Record struct {
	Fields []Field
}

func (self Record) Compare(other Record) util.Order {
	max := other.Len()-1
	var i int
	var f Field

	for i, f = range self.Fields {
		if i > max {
			return util.Gt
		}

		otherValue := other.Get(f.Column)
		
		if otherValue == nil {
			return util.Gt
		}
		
		if result := f.Column.Compare(f.Value, otherValue); result != util.Eq {
			return result
		}
	}

	if max > i {
		return util.Lt
	}

	return util.Eq
}

func (self Record) Find(column Column) (int, bool) {	
	min, max := 0, self.Len()

	for min < max {
		i := (min+max) / 2

		switch util.ComparePointer(column.Pointer(), self.Fields[i].Column.Pointer()) {		
			case util.Lt:
				max = i
			case util.Gt:
				min = i+1
			default:
				return i, true
		}
	}

	return min, false
}

func (self Record) Get(column Column) interface{} {
	i, ok := self.Find(column)

	if !ok {
		return nil
	}

	return self.Fields[i].Value
}

func (self *Record) Len() int {
	return len(self.Fields)
}

func (self *Record) Set(column Column, value interface{}) {
	i, ok := self.Find(column)
	
	if ok {
		self.Fields[i].Value = value
	} else {
		f := Field{Column: column, Value: value}
		l := self.Len()
		
		if i == l {
			self.Fields = append(self.Fields, f)
		} else if i == l-1 {
			self.Fields = append(self.Fields[:i], f, self.Fields[i])
		} else {
			self.Fields = append(self.Fields, f)
			copy(self.Fields[i+1:], self.Fields[i:])
			self.Fields[i] = f
		}
	}
}
