package dom

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"sync"
)

var (
	nodePool sync.Pool
)

type Node struct {
	id string
	attributes map[string]interface{}
	content []interface{}
	script bytes.Buffer
}

func finalizeNode(node *Node) {
	nodePool.Put(node)
}

func NewNode(id string) *Node {
	n := nodePool.Get()

	if n == nil {
		n = new(Node)
	} else {
		n := n.(*Node)
		n.attributes = nil
		n.content = nil
	}

	runtime.SetFinalizer(n, finalizeNode)
	return n.(*Node).Init(id)
}

func (self *Node) Init(id string) *Node {
	self.id = id
	return self
}

func (self *Node) Append(spec string, args...interface{}) *Node {
	self.content = append(self.content, fmt.Sprintf(spec, args...))
	return self
}

func (self *Node) AppendNode(node *Node) *Node {
	self.content = append(self.content, node)
	return self
}

func (self *Node) NewNode(id string) *Node {
	n := NewNode(id)
	self.AppendNode(n)
	return n
}

func (self *Node) Get(key string) interface{} {
	if self.attributes == nil {
		return nil
	}
	
	return self.attributes[key]
}

func (self *Node) Set(key string, val interface{}) *Node {
	if self.attributes == nil {
		self.attributes = make(map[string]interface{})
	}

	self.attributes[key] = val
	return self
}

func (self *Node) Write(out io.Writer) error {
	fmt.Fprintf(out, "<%v", self.id)

	if self.attributes != nil {
		for k, v := range self.attributes {
			if v != nil {
				fmt.Fprintf(out, " %v=\"%v\"", k, v)
			}
		}

		for k, v := range self.attributes {
			if v == nil {
				fmt.Fprintf(out, " %v", k)
			}
		}
	}
	
	if len(self.content) == 0 {
		io.WriteString(out, "/>\n")
	} else {
		io.WriteString(out, ">\n")
		
		for _, v := range self.content {
			switch v := v.(type) {
			case string:
				io.WriteString(out, v)
				io.WriteString(out, "\n")
			case *Node:
				v.Write(out)
			default:
				return fmt.Errorf("Invalid node: %v", v)
			}
		}

		fmt.Fprintf(out, "</%v>\n", self.id)
	}

	return nil
}

func (self *Node) WriteJS(out io.Writer) {
	io.Copy(out, &self.script)
	self.script.Reset()

	for _, v := range self.content {
		switch v := v.(type) {
		case *Node:
			v.WriteJS(out)
		default:
			break
		}
	}
}
