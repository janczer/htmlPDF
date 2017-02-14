package htmlPDF

import "fmt"

type Node struct {
	parent    *Node
	children  map[int]interface{}
	node_type NodeType
}

type NodeType struct {
	element ElementData
	text    string
}

type ElementData struct {
	tag_name string
	attrs    map[string]string
}

func (n *Node) Start(name string) *Node {
	t := new(Node)
	t.parent = n
	t.children = make(map[int]interface{})
	t.node_type = NodeType{element: ElementData{tag_name: name}}
	return t
}

func (n *Node) Stop() *Node {
	p := n.parent
	i := len(p.children)
	p.children[i] = n
	return p
}

func (n *Node) AddText(text string) {
	if len(text) > 0 {
		i := len(n.children)
		n.children[i] = text
	}
}

func createNodeText(data string) *Node {
	return &Node{
		children:  make(map[int]interface{}),
		node_type: NodeType{text: data},
	}
}

func NewNode() *Node {
	return &Node{
		children: make(map[int]interface{}),
		node_type: NodeType{
			element: ElementData{
				attrs: make(map[string]string),
			},
		},
	}
}

func (e ElementData) id() string {
	return e.attrs["id"]
}

//change: return map with class or struct
func (e ElementData) classes() string {
	return e.attrs["class"]
}

func tab(i int) {
	for j := 0; j < i; j++ {
		fmt.Printf("  ")
	}
}

func (n *Node) print(l int) {
	tab(l)
	l++
	fmt.Printf("%s text: %s\n", n.node_type.element.tag_name, n.node_type.text)
	for i := 0; i < len(n.children); i++ {
		switch str := n.children[i].(type) {
		case *Node:
			str.print(l + 1)
		case string:
			tab(l)
			fmt.Printf("text: %s\n", str)
		}
	}
}
