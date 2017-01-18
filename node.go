package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

type Node struct {
	name string
	text string
	margin_left int
	margin_top int
	parent *Node
	child []*Node
}

func (n *Node) Start(name string) *Node {
	tmp := new(Node)
	tmp.parent = n
	tmp.name = name
	switch name {
	case "div":
		tmp.margin_top = 5 + n.margin_top
		tmp.margin_left = 10 + n.margin_left
	case "h1":
		tmp.margin_top = 5 + n.margin_top
		tmp.margin_left = 5 + n.margin_left
		n.margin_top += 5
	}
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
		tab(level+1)
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

func (n *Node) PrintSelf(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetFontSize(16)
	n.pdf(pdf)

	return pdf
}

func (n *Node) pdf(pdf *gofpdf.Fpdf) {
	switch n.name {
	case "div":
		x, y := pdf.GetXY()
		pdf.SetXY(x + float64(n.margin_left), y + float64(n.margin_top))
	case "h1":
		pdf.Text(float64(n.margin_left), float64(n.margin_top), n.text)
	}
	for i:=0; i < len(n.child); i++ {
		n.child[i].pdf(pdf)
	}
}
