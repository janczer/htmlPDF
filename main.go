package htmlPDF

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
)

//global pointer to pdf
var pdf *gofpdf.Fpdf

func Generate(html string, css string, out string) {
	xmlFile, err := ioutil.ReadFile(html)
	if err != nil {
		return
	}

	//parse html to Node tree
	n := ParseHtml(string(xmlFile))
	n.print(0)

	cssFile, err := ioutil.ReadFile(css)
	if err != nil {
		return
	}
	cssStyle := string(cssFile)
	p2 := CssParser(cssStyle)
	stylesheet := p2.parseRules()

	styletree := styleTree(n, &stylesheet)

	viewport := Dimensions{}
	viewport.content.width = 210
	viewport.content.height = 600

	layoutTree := layoutTree(styletree, viewport)
	list := buildDisplayList(layoutTree)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)
	for i := 0; i < len(list); i++ {
		list[i].draw(pdf)
	}
	err = pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Println("Error pdf", err)
	}
	pdf.Close()
}
