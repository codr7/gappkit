package dom

import (
	"bytes"
	"fmt"
	"io"
)

type Doc struct {
	Node
	Head Node
	Title Node
	Body Node
}

func NewDoc(title string) *Doc {
	return new(Doc).Init(title)
}

func (self *Doc) Init(title string) *Doc {
	self.Node.Init("html")
	self.AppendNode(self.Head.Init("head"))
	self.Head.AppendNode(self.Title.Init("title"))
	self.Title.Append(title)
	self.AppendNode(self.Body.Init("body"))
	return self
}

func (self *Doc) CSS(href string) *Node {
	return self.Head.NewNode("link").Set("rel", "stylesheet").Set("href", href)
}

func (self *Doc) Script(src string) *Node {
	return self.Head.NewNode("script").Set("src", src).Append("")
}

func (self *Doc) Write(out io.Writer) error {
	var js bytes.Buffer
	self.WriteJS(&js)
	
	if js.Len() > 0 {
		self.Head.NewNode("script").Append(
			fmt.Sprintf("document.addEventListener('DOMContentLoaded', (event) => {\n%v\n});",
				js.String()))
	}
	
	io.WriteString(out, "<!DOCTYPE html>\n")
	self.Node.Write(out)
	return nil
}
