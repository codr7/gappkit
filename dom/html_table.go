package dom

type Table struct {
	Node
}

type TableRow struct {
	Node
}

func (self *Node) Table(id string) *Table {
	t := new(Table)
	t.Append("")
	self.AppendNode(t.Init("table").Set("id", id))
	return t
}

func (self *Table) Tr() *TableRow {
	n := new(TableRow)
	self.AppendNode(n.Init("tr"))
	return n
}

func (self *TableRow) Th() *Node {
	return self.NewNode("th")
}

func (self *TableRow) Td() *Node {
	return self.NewNode("td")
}
