package htmlPDF

import (
	"fmt"
	"strconv"
	"strings"
)

var pageFontSize float64 = 12
var pageFontFamily string = "Helvetica"
var pageFontColor string = "black"
var pageFontStyle string = "" //B - bold or I - italic or U - underscore
var newline bool

var beforElement string
var actualElement string

var lastMarginTop float64
var marginTop float64

var marginBottop float64

var t *Table
var countList int = 0
var ol bool = false
var ul bool = false

func startElement(name string) {
	beforElement = actualElement
	actualElement = name
	switch name {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		marginTop = 12
		marginBottop = 12
	case "p":
		marginTop = 12
		marginBottop = 12
	case "a":
	case "table":
		t = NewTable()
	case "tr":
		t.startTr()
	case "b", "i", "span":
		addSpace()
	case "br":
		setNewLine(true)
	case "ol":
		ol = true
		countList = 0
	case "ul":
		ul = true
		countList = 0
	case "li":
		countList++
	}
}

func printText(text string) {
	if len(text) <= 0 || actualElement == "div" || actualElement == "page" {
		return
	}

	switch actualElement {
	case "b", "i":
		pageFontStyle = strings.ToUpper(actualElement)
	case "li":
		if ul {
			text = "* " + text
		}
		if ol {
			text = strconv.Itoa(countList) + ". " + text
		}
	case "h1", "h2", "h3", "h4", "h5", "h6":
		setStyleHX()
	case "td":
		t.addTd(text)
		return
	case "br", "ul", "ol":
		return
	}

	drawText(text)
}

func endElement(name string) {
	marginTop = 0
	actualElement = beforElement
	switch name {
	case "div":
		setNewLine(false)
	case "p":
		setNewLine(false)
	case "a":
		setNewLine(false)
	case "h1", "h2", "h3", "h4", "h5", "h6", "li":
		setNewLine(false)
	case "table":
		t.printSelf()
	case "tr":
		t.endTr()
	case "b", "i", "span":
		addSpace()
	case "ol", "ul":
		ol = false
		ul = false
	}
	pageFontStyle = ""
}

func addSpace() {
	if !newline {
		x := pdf.GetX()
		pdf.SetX(x + pdf.GetStringWidth(" "))
	}
}

func drawText(text string) {
	pdf.SetXY(pdf.GetX(), pdf.GetY()+marginTop)

	pdf.SetFont(pageFontFamily, pageFontStyle, pageFontSize)

	pdf.Text(pdf.GetX(), pdf.GetY(), text)
	pdf.SetX(pdf.GetX() + pdf.GetStringWidth(text))
	newline = false
}

func setStyleHX() {
	fmt.Println("drawHX")
	pageFontStyle = "B"
	switch actualElement {
	case "h1":
		pageFontSize = 18
	case "h2":
		pageFontSize = 16
	case "h3":
		pageFontSize = 14
	case "h4":
		pageFontSize = 12
	case "h5", "h6":
		pageFontSize = 10
	}
}

func setNewLine(force bool) {
	if !newline || force {
		y := pdf.GetY()
		_, fontSize := pdf.GetFontSize()
		pdf.SetY(y + fontSize + marginTop)
		newline = true
	}
}
