package main

import (
	"fmt"
)

var pageFontSize float64 = 12
var pageFontFamily string = "Helvetica"
var pageFontColor string = "black"
var pageFontStyle string = "" //B - bold or I - italic or U - underscore
var newline bool
var actualElement string
var end bool
var start bool
var table bool = false
var tr bool = false
var td bool = false
var t *Table

func printElement(name string) {
	actualElement = name
	switch name {
	case "div":
		setNewLine()
	case "p":
		setNewLine()
	case "h1", "h2", "h3", "h4", "h5", "h6":
		setNewLine()
	case "table":
		table = true
		t = NewTable()
	case "tr":
		t.startTr()
		tr = true
	case "td":
		td = true
	}
}

func printText(text string) {
	if len(text) <= 0 || actualElement == "div" || actualElement == "page" {
		return
	}
	if table || tr {
		if td {
			t.addTd(text)
		}
		return
	}
	switch actualElement {
	case "b":
		pageFontStyle = "B"
	case "i":
		pageFontStyle = "I"
	case "h1", "h2", "h3", "h4", "h5", "h6":
		drawHX(text)
		return
	}
	drawText(text)
}

func printEndElement(name string) {
	actualElement = ""
	switch name {
	case "div":
		setNewLine()
	case "p":
		setNewLine()
	case "h1", "h2", "h3", "h4", "h5", "h6":
		setNewLine()
	case "table":
		table = false
		t.printSelf()
	case "tr":
		t.endTr()
		tr = false
	case "td":
		td = false
	case "b", "i":
		pageFontStyle = ""
	}
}

func drawText(text string) {
	pdf.SetFont(pageFontFamily, pageFontStyle, pageFontSize)
	pdf.Text(pdf.GetX(), pdf.GetY(), text)
	x := pdf.GetX()
	pdf.SetX(x + pdf.GetStringWidth(text))
	newline = false
}

func drawHX(text string) {
	fmt.Println("drawHX")
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

func setNewLine() {
	if !newline {
		y := pdf.GetY()
		_, fontSize := pdf.GetFontSize()
		pdf.SetY(y + fontSize)
		newline = true
	}
}
