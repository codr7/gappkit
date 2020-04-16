package db

type Record interface {
	Id() RecordId
}

type BasicRecord struct {
	id RecordId
}

func (self *BasicRecord) Init(id RecordId) {
	self.id = id
}

func (self *BasicRecord) Id() RecordId {
	return self.id
}
