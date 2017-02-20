package htmlPDF

import (
	"github.com/jung-kurt/gofpdf"
)

type DisplayCommand struct {
	command interface{}
}

type SolidColor struct {
	color Color
	rect  Rect
}

func (d DisplayCommand) draw(pdf *gofpdf.Fpdf) {
	switch command := d.command.(type) {
	case SolidColor:
		r := command.rect
		c := command.color
		pdf.SetFillColor(int(c.r), int(c.g), int(c.b))
		pdf.Rect(r.x, r.y, r.width, r.height, "F")
	}
}

func buildDisplayList(layoutRoot *LayoutBox) map[int]DisplayCommand {

}
