package main

import (
	"fmt"
)

type Node struct {
	name   string
	text   string //tado change string to map[int64]string 0 first text, 1 text after first child, 2 ...
	parent *Node
	child  []*Node //todo change []*Node to map[int64]*Node
}

func (n *Node) Start(name string) *Node {
	tmp := new(Node)
	tmp.parent = n
	tmp.name = name

	return tmp
}

func (n *Node) AddText(text string) {
	n.text = text
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
