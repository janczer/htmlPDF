package htmlPDF

import (
	"fmt"
	"strings"
)

//Struct for Node tree
type Node struct {
	children  map[int]*Node
	node_type NodeType
}

type NodeType struct {
	element ElementData
	text    string
}

type ElementData struct {
	tag_name string
	attr     map[string]string
}

func (e ElementData) id() string {
	return e.attr["id"]
}

func (e ElementData) classes() []string {
	class, ok := e.attr["class"]

	if ok {
		return strings.Split(class, " ")
	}
	return []string{}
}

func tab(i int) {
	for j := 0; j < i; j++ {
		fmt.Printf("  ")
	}
}

func (n *Node) print(l int) {
	tab(l)
	l++
	fmt.Printf("%s\n", n.node_type.element.tag_name)
	for i := 0; i < len(n.children); i++ {
		tab(l)
		fmt.Printf("%s text: %s\n", n.children[i].node_type.element.tag_name, n.children[i].node_type.text)
		n.children[i].print(l + 1)
	}
}

func text(data string) *Node {
	return &Node{
		children: map[int]*Node{},
		node_type: NodeType{
			element: ElementData{attr: map[string]string{}},
			text:    data,
		},
	}
}

func elem(name string, attrs map[string]string, children map[int]*Node) *Node {
	return &Node{
		children: children,
		node_type: NodeType{
			element: ElementData{
				tag_name: name,
				attr:     attrs,
			},
		},
	}
}
