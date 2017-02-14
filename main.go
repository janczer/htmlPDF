package htmlPDF

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io"
	"io/ioutil"
	"strings"
)

//global pointer to pdf
var pdf *gofpdf.Fpdf

func parse(source string) {
	r := bytes.NewReader([]byte(source))
	d := xml.NewDecoder(r)
	n := NewNode()
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
			fmt.Printf("start %s\n", start.Name.Local)
			n = n.Start(start.Name.Local)
		case xml.EndElement:
			n = n.Stop()
			//end := token.(xml.EndElement)
		case xml.CharData:
			text := string(token.(xml.CharData))
			if len(strings.TrimSpace(text)) > 0 {
				n.AddText(text)
			}
		}
	}

	fmt.Println(n)
	i := 0
	n.print(i)
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

func Generate(in string, out string) {
	fmt.Println(in, out)
	xmlFile, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	xmlstring := strings.Replace(string(xmlFile), "\n", "", -1)
	xmlstring = strings.Replace(string(xmlstring), "\r", "", -1)

	pdf = gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "B", 16)
	pdf.AddPage()

	parse(xmlstring)

	//Generate PDF Start
	err = pdf.OutputFileAndClose(out)
	if err != nil {
		fmt.Println("Error with generate pdf", err)
	}
	////Generate PDF End
}
