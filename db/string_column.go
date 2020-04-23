package db

type StringColumn struct {
	BasicColumn
}

func (self *StringColumn) Init(name string) *StringColumn {
	self.BasicColumn.Init(name)
	return self
}

func (self *StringColumn) Compare(x, y interface{}) Order {
	return CompareString(x.(string), y.(string))
}

