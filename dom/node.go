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
	nextNodeId uint64
)

type Node struct {
	tag string
	attributes map[string]interface{}
	content []interface{}
	script bytes.Buffer
}

func finalizeNode(node *Node) {
	nodePool.Put(node)
}

func NewNode(tag string) *Node {
	n := nodePool.Get()

	if n == nil {
		n = new(Node)
	} else {
		n := n.(*Node)
		n.attributes = nil
		n.content = nil
	}

	runtime.SetFinalizer(n, finalizeNode)
	return n.(*Node).Init(tag)
}

func (self *Node) Init(tag string) *Node {
	self.tag = tag
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

func (self *Node) Id() interface{} {
	id := self.Get("id")

	if id == nil {
		id = nextNodeId
		nextNodeId++
		self.Set("id", id)
	}

	return id
}

func (self *Node) NewNode(tag string) *Node {
	n := NewNode(tag)
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
	fmt.Fprintf(out, "<%v", self.tag)

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
				if v != "" {
					io.WriteString(out, v)
					io.WriteString(out, "\n")
				}
			case *Node:
				v.Write(out)
			default:
				return fmt.Errorf("Invalid node: %v", v)
			}
		}

		fmt.Fprintf(out, "</%v>\n", self.tag)
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
