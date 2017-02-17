package htmlPDF

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
)

//global pointer to pdf
var pdf *gofpdf.Fpdf

func Generate(in string, out string) {
	fmt.Println(in, out)
	xmlFile, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	//parse html to Node tree
	n := ParseHtml(string(xmlFile))
	n.print(0)

	cssFile, err := ioutil.ReadFile("style.css")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	cssStyle := string(cssFile)
	fmt.Println(cssStyle)
	//todo change NewParser to ParseCSS with return *Rules
	p2 := NewParser(cssStyle)
	p2.parseRules()
}
