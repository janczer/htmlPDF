package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	//"github.com/jung-kurt/gofpdf"
	"io"
	"io/ioutil"
	"strings"
)

type Node struct {
	name string
	text string
	parent *Node
	child []*Node
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

func main() {
	xmlFile, err := ioutil.ReadFile("test.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	xmlstring := strings.Replace(string(xmlFile), "\n", "", -1)
	xmlstring = strings.Replace(string(xmlstring), "\r", "", -1)
	r := bytes.NewReader([]byte(xmlstring))
	d := xml.NewDecoder(r)

	n := new(Node)

	var t int = 0
	for {
		token, err := d.Token()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			panic(err)
		}
		switch token.(type) {
		case xml.StartElement:
			start := token.(xml.StartElement)
			tab(t)
			t++
			fmt.Printf("Name: %s, Attr: %v\n", start.Name, start.Attr)
			n = n.Start(start.Name.Local)
		case xml.EndElement:
			t--
			tab(t)
			end := token.(xml.EndElement)
			fmt.Println(end)
			n = n.Stop()
		case xml.CharData:
			tab(t)
			text := token.(xml.CharData)
			fmt.Printf("%s\n", text)
			n.AddText(string(text))
		}
	}

	n.Print(0)


	//switch n.Name {
	//case "h1":
	//	n.Opt.fontSize = 32
	//	n.Opt.m_left = 20
	//	n.Opt.m_top = 10
	//}
	////Generate PDF Start
	//pdf := gofpdf.New("P", "mm", "A4", "")
	//pdf.AddPage()
	//pdf.SetFont("Arial", "B", 16)
	//pdf.SetFontSize(n.Opt.fontSize)
	//x, y := pdf.GetXY()
	//pdf.Text(x + n.Opt.m_left, y + n.Opt.m_top, n.Text)
	//err = pdf.OutputFileAndClose("hello.pdf")
	//if err != nil {
	//	fmt.Println("Problem z tworzeniem pdf", err)
	//}
	////Generate PDF End
}

func tab(i int) {
	for j := 0; j < i; j++ {
		fmt.Printf("    ")
	}
}
