package main

import (
	"github.com/jung-kurt/gofpdf"
)

var pageFontSize float64 = 12
var pageFontFamily string = "Helvetica"
var pageFontColor string = "black"
var pageFontStyle string = "" //B - bold or I - italic or U - underscore
var newline bool
var actualElement string
var end bool
var start bool

func printElement(pdf *gofpdf.Fpdf, name string) {
	actualElement = name
	switch name {
	case "div":
		setNewLine(pdf)
	case "p":
		setNewLine(pdf)
	case "h1", "h2", "h3", "h4", "h5", "h6":
		setNewLine(pdf)
	}
}

func printText(pdf *gofpdf.Fpdf, text string) {
	if len(text) <= 0 || actualElement == "div" || actualElement == "page" {
		return
	}
	switch actualElement {
	case "b":
		pageFontStyle = "B"
	case "i":
		pageFontStyle = "I"
	case "h1", "h2", "h3", "h4", "h5", "h6":
		drawHX(pdf, text)
		return
	default:
		pageFontStyle = ""
	}
	drawText(pdf, text)
}

func printEndElement(pdf *gofpdf.Fpdf, name string) {
	actualElement = ""
	switch name {
	case "div":
		setNewLine(pdf)
	case "p":
		setNewLine(pdf)
	case "h1", "h2", "h3", "h4", "h5", "h6":
		setNewLine(pdf)
	}
}

func drawText(pdf *gofpdf.Fpdf, text string) {
	pdf.SetFont(pageFontFamily, pageFontStyle, pageFontSize)
	pdf.Text(pdf.GetX(), pdf.GetY(), text)
	x := pdf.GetX()
	pdf.SetX(x + pdf.GetStringWidth(text))
	newline = false
}

func drawHX(pdf *gofpdf.Fpdf, text string) {
	var fontSize float64
	pageFontStyle = "B"
	switch actualElement {
	case "h1":
		fontSize = 18
	case "h2":
		fontSize = 16
	case "h3":
		fontSize = 14
	case "h4":
		fontSize = 12
	case "h5", "h6":
		fontSize = 10
	}
	pdf.SetFont(pageFontFamily, pageFontStyle, fontSize)
	pdf.Text(pdf.GetX(), pdf.GetY(), text)
	newline = false
}

func setNewLine(pdf *gofpdf.Fpdf) {
	if !newline {
		y := pdf.GetY()
		_, fontSize := pdf.GetFontSize()
		pdf.SetY(y + fontSize)
		newline = true
	}

}

//------------------
//------table-------
//------------------

//func drawTable(pdf *gofpdf.Fpdf, n *Node) {
//	t := make(map[int]float64)
//	for i := 0; i < len(n.child); i++ {
//		tmp := n.child[i]
//		for j := 0; j < len(tmp.child); j++ {
//			stringSize := pdf.GetStringWidth(tmp.child[j].text[0])
//			if t[i] < stringSize {
//				t[i] = stringSize
//			}
//		}
//	}
//	fmt.Println(t)
//	for i := 0; i < len(n.child); i++ {
//		drawTr(pdf, n.child[i], t)
//	}
//	y := pdf.GetY()
//	_, fontSize := pdf.GetFontSize()
//	pdf.SetY(y + fontSize)
//}
//
//func drawTr(pdf *gofpdf.Fpdf, n *Node, t map[int]float64) {
//	y := pdf.GetY()
//	_, fontSize := pdf.GetFontSize()
//	pdf.SetY(y + fontSize)
//
//	for i := 0; i < len(n.child); i++ {
//		pdf.Text(pdf.GetX(), pdf.GetY(), n.child[i].text[0])
//		x := pdf.GetX()
//		pdf.SetX(x + t[i] + 1)
//	}
//}
