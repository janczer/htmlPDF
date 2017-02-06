package main

import (
	"github.com/jung-kurt/gofpdf"
)

func drawInLine(pdf *gofpdf.Fpdf, n *Node) {
	switch n.name {
	case "b":
		pdf.SetFont(pageFontFamily, "B", pageFontSize)
		pdf.Text(pdf.GetX(), pdf.GetY(), n.text[0])
		x := pdf.GetX()
		pdf.SetX(x + pdf.GetStringWidth(n.text[0]) + 1)
		newline = false
	case "i":
		pdf.SetFont(pageFontFamily, "I", pageFontSize)
		pdf.Text(pdf.GetX(), pdf.GetY(), n.text[0])
		x := pdf.GetX()
		pdf.SetX(x + pdf.GetStringWidth(n.text[0]) + 1)
		newline = false
	}
}
