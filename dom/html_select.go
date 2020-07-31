package dom

type Select struct {
	Node
}

func (self *Node) Select(id string) *Select {
	n := new(Select)
	self.AppendNode(n.Init("select").Set("id", id))
	n.Append("")
	return n
}

func (self *Select) Option(value, caption string) *Node {
	n := self.NewNode("option")
	n.Set("value", value)
	n.Append(caption)
	return n
}
