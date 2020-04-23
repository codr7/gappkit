package db

type RecordSetColumn struct {
	BasicColumn
}

func (self *RecordSetColumn) Init(name string) *RecordSetColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *RecordSetColumn) Compare(x, y interface{}) Order {
	return x.(RecordSet).Compare(y.(RecordSet))
}

