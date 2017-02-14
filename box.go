package htmlPDF

type Dimensions struct {
	content Rect

	padding EngeSizes
	border  EdgeSizes
	margin  EdgeSizes
}

type Rect struct {
	x float64
	y float64

	width  float64
	height float64
}

type EdgeSizes struct {
	left   float64
	right  float64
	top    float64
	bottom float64
}

type LayoutBox struct {
	dimensions Dimensions
	box_type   interface{} //box_type can be a block node, an inline node, or an anonymous block box
	childre    map[int]LayoutBox
}

type BlockNode struct {
}

type InlineNode struct {
}

type AnonymousBlock struct {
}
