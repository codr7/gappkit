package dom

type TableBody struct {
	Node
}

type TableRow struct {
	Node
}

func (self *Node) Table(id string) *TableBody {
	t := self.NewNode("table").Set("id", id)
	b := new(TableBody)
	t.AppendNode(b.Init("tbody"))
	return b
}

func (self *TableBody) Row() *TableRow {
	n := new(TableRow)
	self.AppendNode(n.Init("tr"))
	return n
}

func (self *TableRow) Data() *Node {
	return self.NewNode("td")
}
