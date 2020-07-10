package dom

import (
	"fmt"
	"github.com/google/uuid"
)

var br Node

func init() {
	br.Init("br")
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

func (self *Node) Id() interface{} {
	id := self.Get("id")

	if id == nil {
		id = uuid.New().String()
		self.Set("id", id)
	}

	return id
}

func (self *Node) Input(id string, inputType string) *Node {
	return self.NewNode("input").
		Set("id", id).
		Set("type", inputType).
		Set("size", -1) 
}

func (self *Node) Label(id string, caption string) *Node {
	return self.NewNode("label").Set("id", id).Append(caption)
}

func (self *Node) OnClick(spec string, args...interface{}) *Node {
	id := self.Get("id")

	if id == nil {
		self.Set("id", uuid.New().String())
	}

	fmt.Fprintf(&self.script,
		"document.getElementById('%v').addEventListener('click', (event) => {\n%v\n});",
		id,
		fmt.Sprintf(spec, args...))

	return self
}
