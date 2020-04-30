package db

type StoredField struct {
	Column string
	Value interface{}
}

type StoredRecord struct {
	Fields []StoredField
}

func (self StoredRecord) Load(table *Table) *Record {
	out := new(Record)
	out.Fields = make([]Field, len(self.Fields))
	
	for i, f := range self.Fields {
		of := &out.Fields[i]
		of.Column = table.FindColumn(f.Column)
		of.Value = f.Value
	}

	return out
}

func (self *StoredRecord) Set(column Column, value interface{}) {
	self.Fields = append(self.Fields, StoredField{Column: column.Name(), Value: value})
}
