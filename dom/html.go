package dom

import (
	"fmt"
)

var br Node

func init() {
	br.Init("br")
}

func (self *Node) Autofocus() *Node {
	return self.Set("autofocus", nil)
}

func (self *Node) Br() *Node {
	self.AppendNode(&br)
	return self
}

func (self *Node) Button(id string, caption string) *Node {
	return self.NewNode("button").Set("id", id).Append(caption)
}

func (self *Node) Div(id string) *Node {
	return self.NewNode("div").Set("id", id)
}

func (self *Node) FieldSet(id string) *Node {
	return self.NewNode("fieldset").Set("id", id)
}

func (self *Node) H1(caption string) *Node {
	self.NewNode("h1").Append(caption)
	return self
}

func (self *Node) Input(id string, inputType string) *Node {
	return self.NewNode("input").
		Set("id", id).
		Set("type", inputType)
}

func (self *Node) Label(id string, caption string) *Node {
	return self.NewNode("label").Set("id", id).Append(caption)
}

func (self *Node) OnClick(spec string, args...interface{}) *Node {
	fmt.Fprintf(&self.script,
		"document.getElementById('%v').addEventListener('click', (event) => {\n%v\n});",
		self.Id(),
		fmt.Sprintf(spec, args...))

	return self
}
