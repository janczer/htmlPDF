package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io"
	"io/ioutil"
	"strings"
)

type Node struct {
	Name string
	Opt  Options
	Text string
}

type Options struct {
	fontSize float64
	m_left   float64
	m_top    float64
	m_right  float64
	m_bottom float64
}

type Nodes struct {
	n []Node
}

func (n *Nodes) Add(nn Node) []Node {
	n.n = append(n.n, nn)
	return n.n
}

func main() {
	xmlFile, err := ioutil.ReadFile("test.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	//	xmlPure := strings.Map(func(r rune) rune {
	//		if unicode.IsSpace(r) {
	//			return -1
	//		}
	//		return r
	//	}, string(xmlFile))

	xmlstring := strings.Replace(string(xmlFile), "\n", "", -1)
	xmlstring = strings.Replace(string(xmlstring), "\r", "", -1)
	r := bytes.NewReader([]byte(xmlstring))
	d := xml.NewDecoder(r)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	var n Node
	var nnn Nodes

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
			n.Name = start.Name.Local
		case xml.EndElement:
			t--
			tab(t)
			end := token.(xml.EndElement)
			fmt.Println(end)
			nnn.Add(n)
		case xml.CharData:
			tab(t)
			text := token.(xml.CharData)
			fmt.Printf("%s\n", text)
			n.Text = string(text)
		}
	}
	fmt.Println(nnn.n)

	switch n.Name {
	case "h1":
		n.Opt.fontSize = 32
		n.Opt.m_left = 20
		n.Opt.m_top = 10
	}
	pdf.SetFontSize(n.Opt.fontSize)
	x, y := pdf.GetXY()
	pdf.Text(x+n.Opt.m_left, y+n.Opt.m_top, n.Text)
	err = pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Println("Problem z tworzeniem pdf", err)
	}

	fmt.Println(n)
}

func tab(i int) {
	for j := 0; j < i; j++ {
		fmt.Printf("\t")
	}
}
