package htmlPDF

type Dimensions struct {
	content Rect

	padding EdgeSizes
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
	style StyleNode
}

type InlineNode struct {
	style StyleNode
}

type AnonymousBlock struct {
	style StyleNode
}

func (s StyleNode) value(name string) Value {
	val, ok := s.specified_values[name]
	if ok {
		return val
	}
	return Value{}
}

func (s StyleNode) display() string {
	val, ok := s.specified_values["display"]
	if ok {
		return val.keyword
	}
	return "inline"
}

func NewLayoutBox(boxType interface{}) LayoutBox {
	return LayoutBox{
		box_type: boxType,
	}
}

func buildLayoutTree(styleNode StyleNode) LayoutBox {
	display := styleNode.display()
	var boxType interface{}
	switch display {
	case "block":
		boxType = BlockNode{styleNode}
	case "inline":
		boxType = InlineNode{styleNode}
	default:
		panic("Root node has display: none.")
	}
	root := NewLayoutBox(boxType)
	return root
}
