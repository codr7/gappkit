package db

type RecordSet struct {
	Items []RecordId
}

func (self *RecordSet) Add(id RecordId) {
	i, _ := self.Find(id)
	self.Insert(i, id)
}

func (self *RecordSet) Find(id RecordId) (int, bool) {
	min, max := 0, len(self.Items)

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
	l := len(self.Items)
	
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

func (self *RecordSet) Remove(id RecordId) bool {
	i, ok := self.Find(id)	

	if !ok {
		return false
	}
	
	l := len(self.Items)
	copy(self.Items[i:], self.Items[i+1:])
	self.Items = self.Items[:l-1]	
	return true
}
