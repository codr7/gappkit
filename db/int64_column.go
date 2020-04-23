package db

type Int64Column struct {
	BasicColumn
}

func (self *Int64Column) Init(name string) *Int64Column {
	self.BasicColumn.Init(name)
	return self
}

func (self *Int64Column) Compare(x, y interface{}) Order {
	return CompareInt64(x.(int64), y.(int64))
}

