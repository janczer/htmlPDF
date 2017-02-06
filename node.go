package main

import (
	"fmt"
)

type Node struct {
	name   string
	text   map[int64]string
	number int64
	parent *Node
	child  []*Node
}

func (n *Node) Start(name string) *Node {
	tmp := new(Node)
	tmp.parent = n
	tmp.name = name
	tmp.number = 0
	tmp.text = make(map[int64]string)

	return tmp
}

func (n *Node) AddText(text string) {
	n.text[n.number] = text
	n.number++
}

func (n *Node) Stop() *Node {
	n.parent.child = append(n.parent.child, n)
	return n.parent
}

func (n *Node) Print(level int) {
	tab(level)
	level++
	fmt.Printf("<%s>\n", n.name)
	if len(n.text) > 0 {
		tab(level + 1)
		fmt.Printf("%s\n", n.text)
	}
	if len(n.child) > 0 {
		for i := 0; i < len(n.child); i++ {
			n.child[i].Print(level)
		}
	}
	level--
	tab(level)
	fmt.Printf("</%s>\n", n.name)
}
