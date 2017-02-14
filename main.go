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

func parse(source string) *Node {
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

	//draw tree
	n.print(0)

	return n
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

	//parse xml to node tree
	//n *Node
	n := parse(xmlstring)

	//Generate PDF Start
	err = pdf.OutputFileAndClose(out)
	if err != nil {
		fmt.Println("Error with generate pdf", err)
	}
	////Generate PDF End
}
