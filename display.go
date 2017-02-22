package htmlPDF

import (
	"fmt"
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
	return renderLayoutBox(layoutRoot)
}

func renderLayoutBox(layoutBox *LayoutBox) map[int]DisplayCommand {
	list := map[int]DisplayCommand{}
	fmt.Println("===")
	//renderBackground
	colorBackrgound := getColor(layoutBox, "background")
	backgroundCommand := DisplayCommand{
		command: SolidColor{
			color: colorBackrgound,
			rect:  layoutBox.dimensions.borderBox(),
		},
	}
	list[len(list)] = backgroundCommand

	//TODO renderBorders

	//TODO renderText

	return list
}

//Return the specified color for CSS property name
func getColor(layoutBox *LayoutBox, name string) Color {
	switch layoutBox.box_type.(type) {
	case BlockNode, InlineNode:
		return layoutBox.style.value(name).color
	case AnonymousBlock:
		return Color{255, 255, 255, 0}
	default:
		return Color{255, 255, 255, 0}
	}
}
