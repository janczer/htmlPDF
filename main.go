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
	n = n.child[0]

	n.Print(0)

	//Generate PDF Start
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf = n.PrintSelf(pdf)
	err = pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Println("Error pdf", err)
	}
	////Generate PDF End
}

func tab(i int) {
	for j := 0; j < i; j++ {
		fmt.Printf("    ")
	}
}
