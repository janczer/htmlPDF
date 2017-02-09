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

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "B", 16)
	pdf.AddPage()

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
			printElement(pdf, start.Name.Local)
		case xml.EndElement:
			end := token.(xml.EndElement)
			printEndElement(pdf, end.Name.Local)
		case xml.CharData:
			text := string(token.(xml.CharData))
			printText(pdf, strings.TrimSpace(text))
		}
	}

	//Generate PDF Start
	err = pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Println("Error with generate pdf", err)
	}
	////Generate PDF End
}

func tab(i int) {
	for j := 0; j < i; j++ {
		fmt.Printf("    ")
	}
}
