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
	fmt.Println("\x1b[41m\x1b[1mprint Node\x1b[0m")
	n.print(0)
	fmt.Println("\x1b[41m\x1b[1mend print Node\x1b[0m\n")

	cssFile, err := ioutil.ReadFile(css)
	if err != nil {
		return
	}
	cssStyle := string(cssFile)
	p2 := CssParser(cssStyle)
	stylesheet := p2.parseRules()

	styletree := styleTree(n, &stylesheet)
	fmt.Println("\x1b[41m\x1b[1mprint StyleTree\x1b[0m")
	styletree.print(0)
	fmt.Println("\x1b[41m\x1b[1mend print StyleTree\x1b[0m\n")

	viewport := Dimensions{}
	viewport.content.width = 210
	viewport.content.height = 600

	layoutTree := layoutTree(styletree, viewport)
	fmt.Println("\n\x1b[41m\x1b[1mprint LayoutTree\x1b[0m")
	layoutTree.print(0)
	fmt.Println("\x1b[41m\x1b[1mend print LayoutTree\x1b[0m")
	list := buildDisplayList(layoutTree)
	fmt.Println(layoutTree)

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

func GenerateFromString(html string, css string, out string) {
	//parse html to Node tree
	n := ParseHtml(string(html))
	fmt.Println("\x1b[41m\x1b[1mprint Node\x1b[0m")
	n.print(0)
	fmt.Println("\x1b[41m\x1b[1mend print Node\x1b[0m\n")

	cssStyle := string(css)
	p2 := CssParser(cssStyle)
	stylesheet := p2.parseRules()

	styletree := styleTree(n, &stylesheet)
	fmt.Println("\x1b[41m\x1b[1mprint StyleTree\x1b[0m")
	styletree.print(0)
	fmt.Println("\x1b[41m\x1b[1mend print StyleTree\x1b[0m\n")

	viewport := Dimensions{}
	viewport.content.width = 210
	viewport.content.height = 600

	layoutTree := layoutTree(styletree, viewport)
	fmt.Println("\n\x1b[41m\x1b[1mprint LayoutTree\x1b[0m")
	layoutTree.print(0)
	fmt.Println("\x1b[41m\x1b[1mend print LayoutTree\x1b[0m")
	list := buildDisplayList(layoutTree)
	fmt.Println(layoutTree)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)
	for i := 0; i < len(list); i++ {
		list[i].draw(pdf)
	}
	err := pdf.OutputFileAndClose(out)
	if err != nil {
		fmt.Println("Error pdf", err)
	}
	pdf.Close()
}
