package db

import (
	"gappkit/compare"
)

type RecordSet struct {
	Items []RecordId
}

func (self *RecordSet) Add(id RecordId) bool {
	i, ok := self.Find(id)

	if ok {
		return false
	}

	self.Insert(i, id)
	return true
}

func (self RecordSet) Compare(other RecordSet) compare.Order {
	max := other.Len()-1
	var i int
	var id RecordId
	
	for i, id = range self.Items {
		if i > max {
			return compare.Gt
		}

		if result := CompareRecordId(id, other.Items[i]); result != compare.Eq {
			return result
		}
	}

	if max > i {
		return compare.Lt
	}

	return compare.Eq
}

func (self RecordSet) Find(id RecordId) (int, bool) {
	min, max := 0, self.Len()

	for min < max {
		i := (min+max) / 2
		item := self.Items[i]

		if id < item {
			max = i
		} else if id > item {
			min = i+1
		} else {
			return i, true
		}
	}

	return min, false
}

func (self *RecordSet) Insert(i int, id RecordId) {
	l := self.Len()
	
	if i == l {
		self.Items = append(self.Items, id)
	} else if i == l-1 {
		self.Items = append(self.Items[:i], id, self.Items[i])
	} else {
		self.Items = append(self.Items, 0)
		copy(self.Items[i+1:], self.Items[i:])
		self.Items[i] = id
	}
}

func (self RecordSet) Len() int {
	return len(self.Items)
}

func (self *RecordSet) Remove(id RecordId) bool {
	i, ok := self.Find(id)	

	if !ok {
		return false
	}
	
	l := self.Len()
	copy(self.Items[i:], self.Items[i+1:])
	self.Items = self.Items[:l-1]	
	return true
}
