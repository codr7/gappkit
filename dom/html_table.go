package dom

type Table struct {
	Node
}

type TableRow struct {
	Node
}

func (self *Node) Table(id string) *Table {
	n := new(Table)
	self.AppendNode(n.Init("table").Set("id", id))
	return n
}

func (self *Table) Row() *TableRow {
	n := new(TableRow)
	self.AppendNode(n.Init("tr"))
	return n
}

func (self *TableRow) Data() *Node {
	return self.NewNode("td")
}
