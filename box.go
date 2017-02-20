package htmlPDF

import "fmt"

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
	children   map[int]*LayoutBox
	style      StyleNode
}

func (l LayoutBox) getStyleNode() StyleNode {
	switch l.box_type.(type) {
	case BlockNode, InlineNode:
		return l.style
	case AnonymousBlock:
		panic("Anonymous block box has no style node")
	default:
		panic("Type must be BlockNode or InlineNode")
	}
}

type BlockNode struct {
}

type InlineNode struct {
}

type AnonymousBlock struct {
}

func (s StyleNode) value(name string) Value {
	val, ok := s.specified_values[name]
	if ok {
		return val
	}
	return Value{}
}

func (s StyleNode) lookup(first string, second string, end Length) Length {
	f, ok := s.specified_values[first]
	if ok {
		return f.length
	}
	f, ok = s.specified_values[second]
	if ok {
		return f.length
	}
	return end
}

func (s StyleNode) display() string {
	val, ok := s.specified_values["display"]
	if ok {
		return val.keyword
	}
	return "inline"
}

func NewLayoutBox(boxType interface{}, style StyleNode) *LayoutBox {
	return &LayoutBox{
		dimensions: Dimensions{},
		box_type:   boxType,
		children:   map[int]*LayoutBox{},
		style:      style,
	}
}

func buildLayoutTree(styleNode StyleNode) *LayoutBox {
	display := styleNode.display()
	var boxType interface{}
	switch display {
	case "block":
		boxType = BlockNode{}
	case "inline":
		boxType = InlineNode{}
	default:
		panic("Root node has display: none.")
	}
	root := NewLayoutBox(boxType, styleNode)

	for _, child := range styleNode.children {
		display = child.display()
		switch display {
		case "block":
			root.children[len(root.children)] = buildLayoutTree(child)
		case "inline":
			inline := root.getInlineContainer()
			inline.children[len(inline.children)] = buildLayoutTree(child)
		default:
		}
	}

	return root
}

func (l *LayoutBox) getInlineContainer() *LayoutBox {
	boxT := l.box_type
	switch boxT.(type) {
	case AnonymousBlock:
		return l
	case BlockNode:
		return NewLayoutBox(BlockNode{}, StyleNode{})
		//switch childBoxType := value.(l.children[len(l.children)-1].box_type) {
		//}
	default:
		return l
	}
}

func layoutTree(node StyleNode, containBlock Dimensions) *LayoutBox {
	containBlock.content.height = 0

	rootBox := buildLayoutTree(node)
	rootBox.layout(containBlock)

	return rootBox
}

func (l *LayoutBox) layout(containBlock Dimensions) {
	switch typ := l.box_type.(type) {
	case BlockNode:
		l.layoutBox(containBlock)
	case InlineNode:
		//TODO
	case AnonymousBlock:
		//TODO
	default:
		fmt.Println(typ)
	}
}

func (l *LayoutBox) layoutBox(containBlock Dimensions) {
	//Child width can depend on parent width, so we need to calculate
	//this box's width before laying out its children.
	l.calculateBlockWidth(containBlock)

	//Determine where the box is located within its container.
	l.calculateBlockPosition(containBlock)

	//Recursibely layout the children of this box
	l.layoutBlockChildren()

	//Parent height can depend on child height, so calculateHeight
	//must be called after the children are laid out.
	l.calculateBlockHeight(containBlock)
}

func (l *LayoutBox) layoutBlockChildren() {
	//TODO
}

//Calculate the width of a block-level non-replaced element in normal flow
//http://www.w3.org/TR/CSS2/visudet.html#blockwidth
//Sets the horizontal margin/padding/border dimesions, and the 'width'
func (l *LayoutBox) calculateBlockWidth(containBlock Dimensions) {
	style := l.getStyleNode()

	//width has initial value auto
	width := style.value("width")

	//margin, border, and padding have initial value 0
	zero := Length{0.0, "px"}

	marginLeft := style.lookup("margin-left", "margin", zero)
	marginRight := style.lookup("margin-right", "margin", zero)

	borderLeft := style.lookup("border-left-width", "border-width", zero)
	borderRight := style.lookup("border-rigth-width", "border-width", zero)

	paddingLeft := style.lookup("padding-left", "padding", zero)
	paddingRight := style.lookup("padding-right", "padding", zero)

	total := GetTotalFrom(marginLeft, marginRight, borderLeft, borderRight, paddingLeft, paddingRight)

	underflow := containBlock.content.width - total

	widthAuto := width.keyword == "auto"
	marginLeftAuto := style.value("margin-left").keyword == "auto"
	marginRightAuto := style.value("margin-right").keyword == "auto"
	widthLength := width.length

	//If the values are overconstrained, calculate margin_rigth
	if !widthAuto && !marginLeftAuto && !marginRightAuto {
		marginRight = Length{value: marginRight.value + underflow}
	}

	//If execly one size is auto, its used value fallows from the equality
	if !widthAuto && !marginLeftAuto && marginRightAuto {
		marginRight = Length{value: underflow}
	}

	if !widthAuto && marginLeftAuto && !marginRightAuto {
		marginLeft = Length{value: underflow}
	}

	if widthAuto {
		if marginLeftAuto {
			marginLeft = Length{}
		}
		if marginRightAuto {
			marginRight = Length{}
		}

		if underflow >= 0 {
			//Expand width to fill the underflow
			widthLength = Length{value: underflow}
		} else {
			//Width can't be negative.Adjust the right margin instead
			widthLength = Length{}
			marginRight = Length{value: marginRight.value + underflow}
		}
	}
	//If margin-left and margin-right are both auto, their used values are equal
	if !widthAuto && marginLeftAuto && marginRightAuto {
		marginLeft = Length{value: underflow / 2}
		marginRight = Length{value: underflow / 2}
	}
	l.dimensions.content.width = widthLength.value

	l.dimensions.padding.left = paddingLeft.value
	l.dimensions.padding.right = paddingRight.value

	l.dimensions.border.left = borderLeft.value
	l.dimensions.border.right = borderRight.value

	l.dimensions.margin.left = marginLeft.value
	l.dimensions.margin.right = marginRight.value
}

//Finish calculating the block's edge sizes, and position it within its containing block
// http://www.w3.org/TR/CSS2/visudet.html#normal-block
//Sets the vertical margin/padding/border dimensions, and the 'x', 'y' values
func (l *LayoutBox) calculateBlockPosition(containBlock Dimensions) {
	style := l.getStyleNode()

	zero := Length{0.0, "Px"}

	l.dimensions.margin.top = style.lookup("margin-top", "margin", zero).value
	l.dimensions.margin.bottom = style.lookup("margin-bottom", "margin", zero).value

	l.dimensions.border.top = style.lookup("border-top-width", "border-width", zero).value
	l.dimensions.border.bottom = style.lookup("border-bottom-width", "border-width", zero).value

	l.dimensions.padding.top = style.lookup("padding-top", "padding", zero).value
	l.dimensions.padding.bottom = style.lookup("padding-bottom", "padding", zero).value

	l.dimensions.content.x = containBlock.content.x + l.dimensions.margin.left + l.dimensions.border.left + l.dimensions.padding.left

	l.dimensions.content.y = containBlock.content.height + containBlock.content.y + l.dimensions.margin.top + l.dimensions.border.top + l.dimensions.padding.top

}

//Height of a block-level non-replaced element in normal flow with overflow visible
func (l *LayoutBox) calculateBlockHeight(containBlock Dimensions) {
	//If the height is set to an explicit length, use that exact lenght
	//Otherwise, jsut keep the value set by 'layoutBlockChildre'
	height := l.getStyleNode().value("height")
	if height.length.value != 0 {
		l.dimensions.content.height = height.length.value
	}
}

func GetTotalFrom(ml, mr, bl, br, pl, pr Length) float64 {
	return ml.value + mr.value + bl.value + br.value + pl.value + pr.value
}
