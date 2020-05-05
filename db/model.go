package db

type Model interface {
	Id() RecordId
	Store() error
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
