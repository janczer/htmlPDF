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

type BlockNode struct{}

type InlineNode struct{}

type AnonymousBlock struct{}

func (r Rect) expandedBy(edge EdgeSizes) Rect {
	return Rect{
		x:      r.x - edge.left,
		y:      r.y - edge.top,
		width:  r.width + edge.left + edge.right,
		height: r.height + edge.top + edge.bottom,
	}
}

func (d Dimensions) paddingBox() Rect {
	return d.content.expandedBy(d.padding)
}

func (d Dimensions) borderBox() Rect {
	return d.paddingBox().expandedBy(d.border)
}

func (d Dimensions) marginBox() Rect {
	return d.borderBox().expandedBy(d.margin)
}

func (d Dimensions) textBox() Rect {
	return Rect{
		x: d.content.x + d.margin.left + d.padding.left + d.border.left,
		y: d.content.y + d.margin.top + d.padding.top + d.border.top,
	}
}

func (s StyleNode) value(name string) Value {
	val, ok := s.specified_values[name]
	if ok {
		return val
	}
	//Return white color and transparent color
	return Value{color: Color{0, 0, 0, 0}}
}

func (s StyleNode) lookup(first string, second string, end Value) Value {
	f, ok := s.specified_values[first]
	if ok {
		return f
	}
	f, ok = s.specified_values[second]
	if ok {
		return f
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

func (s LayoutBox) print(l int) {
	tab(l)
	fmt.Printf("dimensions %+v\n", s.dimensions)
	tab(l)
	fmt.Printf("box type %#v\n", s.box_type)
	//tab(l)
	//fmt.Printf("style %v\n", s.style)
	//tab(l)
	//fmt.Printf("childrens: \n")
	l++
	for i := 0; i < len(s.children); i++ {
		s.children[i].print(l + 1)
	}
}

func buildLayoutTree(styleNode StyleNode) *LayoutBox {

	display := styleNode.display()
	fmt.Println(display)
	var boxType interface{}
	switch display {
	case "block":
		boxType = BlockNode{}
	case "inline":
		boxType = InlineNode{}
	default:
		panic("Root node has display: none.")
	}
	fmt.Println(boxType)

	l := LayoutBox{
		box_type: boxType,
		children: map[int]*LayoutBox{},
		style:    styleNode,
	}

	for i := 0; i < len(styleNode.children); i++ {
		child := styleNode.children[i]
		display = child.display()
		switch display {
		case "block":
			childLayoutTree := buildLayoutTree(child)
			l.children[len(l.children)] = childLayoutTree
		case "inline":
			lastContainer := l.getLastContainer()
			boxT := lastContainer.box_type
			fmt.Printf("last container %#v\n", boxT)
			//add anonymous box
			switch boxT.(type) {
			case AnonymousBlock, InlineNode:
				childLayoutTree := buildLayoutTree(child)

				lastContainer.children[len(lastContainer.children)] = childLayoutTree
			case BlockNode:
				//create AnonymousBlock
				anonymous := LayoutBox{
					box_type: AnonymousBlock{},
					children: map[int]*LayoutBox{},
				}
				//buildLayoutTree
				childLayoutTree := buildLayoutTree(child)
				//add to AnonymousBlock
				anonymous.children[len(anonymous.children)] = childLayoutTree
				//add anonymousBox to child
				l.children[len(l.children)] = &anonymous
			}
		}
	}

	return &l
}

func (l *LayoutBox) getLastContainer() *LayoutBox {
	if len(l.children) == 0 {
		return l
	}

	if len(l.children) == 1 {
		return l.children[0]
	}

	return l.children[len(l.children)-1]
}

func (l *LayoutBox) getInlineContainer() *LayoutBox {
	boxT := l.box_type
	switch boxT.(type) {
	case AnonymousBlock, InlineNode:
		return l
	case BlockNode:
		return NewLayoutBox(AnonymousBlock{}, StyleNode{})
	default:
		return l
	}
}

func layoutTree(node StyleNode, containBlock Dimensions) *LayoutBox {
	containBlock.content.height = 0

	rootBox := buildLayoutTree(node)
	rootBox.layout(&containBlock)

	return rootBox
}

func (l *LayoutBox) layout(containBlock *Dimensions) {
	switch l.box_type.(type) {
	case BlockNode:
		l.layoutBox(containBlock)
	case InlineNode:
		fmt.Println("layout inlinenode")
		l.inlineBox(containBlock)
	case AnonymousBlock:
		fmt.Println("layout anonymous")
		l.anonymousBox(containBlock)
	default:
	}
}

func (l *LayoutBox) inlineBox(containBlock *Dimensions) {
	fmt.Printf("%+v", containBlock.content)

	l.dimensions.content.x = containBlock.content.x
	l.dimensions.content.y = containBlock.content.y
	//l.dimensions.content.y = containBlock.content.height
	d := &l.dimensions

	//calculate box width

	for _, child := range l.children {
		child.layout(d)
	}
}

func (l *LayoutBox) anonymousBox(containBlock *Dimensions) {
	fmt.Printf("%+v", containBlock.content)
	//block position is the same previous
	l.dimensions.content.x = containBlock.content.x
	l.dimensions.content.y = containBlock.content.y
	l.dimensions.content.height = containBlock.content.height
	l.dimensions.content.width = containBlock.content.width

	//Recursibely layout the children of this box
	l.layoutBlockChildren()
}

func (l *LayoutBox) layoutBox(containBlock *Dimensions) {
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
	d := &l.dimensions

	for _, child := range l.children {
		child.layout(d)
		// Track the height so each child is laid out below the previous content.
		d.content.height = d.content.height + child.dimensions.marginBox().height
	}
}

//Calculate the width of a block-level non-replaced element in normal flow
//http://www.w3.org/TR/CSS2/visudet.html#blockwidth
//Sets the horizontal margin/padding/border dimesions, and the 'width'
func (l *LayoutBox) calculateBlockWidth(containBlock *Dimensions) {
	style := l.getStyleNode()

	//width has initial value auto
	width, ok := style.specified_values["width"]
	if !ok {
		width = Value{
			keyword: "auto",
		}
	}

	//margin, border, and padding have initial value 0
	zero := Value{
		length: Length{0.0, "px"},
	}

	marginLeft := style.lookup("margin-left", "margin", zero)
	marginRight := style.lookup("margin-right", "margin", zero)

	borderLeft := style.lookup("border-left-width", "border-width", zero)
	borderRight := style.lookup("border-rigth-width", "border-width", zero)

	paddingLeft := style.lookup("padding-left", "padding", zero)
	paddingRight := style.lookup("padding-right", "padding", zero)

	total := GetTotalFrom(marginLeft, marginRight, borderLeft, borderRight, paddingLeft, paddingRight)

	if width.keyword != "auto" && total > containBlock.content.width {
		if marginLeft.keyword == "auto" {
			marginLeft = Value{length: Length{0, "Px"}}
		}
		if marginRight.keyword == "auto" {
			marginRight = Value{length: Length{0, "Px"}}
		}
	}

	underflow := containBlock.content.width - total

	widthAuto := width.keyword == "auto"
	marginLeftAuto := style.value("margin-left").keyword == "auto"
	marginRightAuto := style.value("margin-right").keyword == "auto"
	widthLength := width

	//If the values are overconstrained, calculate margin_rigth
	if !widthAuto && !marginLeftAuto && !marginRightAuto {
		marginRight = Value{length: Length{marginRight.length.value + underflow, "Px"}}
	}

	//If execly one size is auto, its used value fallows from the equality
	if !widthAuto && !marginLeftAuto && marginRightAuto {
		marginRight.length = Length{value: underflow}
	}

	if !widthAuto && marginLeftAuto && !marginRightAuto {
		marginLeft.length = Length{value: underflow}
	}

	if widthAuto {
		if marginLeftAuto {
			marginLeft = Value{}
		}
		if marginRightAuto {
			marginRight = Value{}
		}

		if underflow >= 0 {
			//Expand width to fill the underflow
			widthLength = Value{length: Length{value: underflow}}
		} else {
			//Width can't be negative.Adjust the right margin instead
			widthLength = Value{}
			marginRight = Value{length: Length{marginRight.length.value + underflow, "Px"}}
		}
	}
	//If margin-left and margin-right are both auto, their used values are equal
	if !widthAuto && marginLeftAuto && marginRightAuto {
		marginLeft.length = Length{value: underflow / 2}
		marginRight.length = Length{value: underflow / 2}
	}
	l.dimensions.content.width = widthLength.length.value

	l.dimensions.padding.left = paddingLeft.length.value
	l.dimensions.padding.right = paddingRight.length.value

	l.dimensions.border.left = borderLeft.length.value
	l.dimensions.border.right = borderRight.length.value

	l.dimensions.margin.left = marginLeft.length.value
	l.dimensions.margin.right = marginRight.length.value
}

//Finish calculating the block's edge sizes, and position it within its containing block
// http://www.w3.org/TR/CSS2/visudet.html#normal-block
//Sets the vertical margin/padding/border dimensions, and the 'x', 'y' values
func (l *LayoutBox) calculateBlockPosition(containBlock *Dimensions) {
	style := l.getStyleNode()

	zero := Value{
		length: Length{0.0, "px"},
	}

	l.dimensions.margin.top = style.lookup("margin-top", "margin", zero).length.value
	l.dimensions.margin.bottom = style.lookup("margin-bottom", "margin", zero).length.value

	l.dimensions.border.top = style.lookup("border-top-width", "border-width", zero).length.value
	l.dimensions.border.bottom = style.lookup("border-bottom-width", "border-width", zero).length.value

	l.dimensions.padding.top = style.lookup("padding-top", "padding", zero).length.value
	l.dimensions.padding.bottom = style.lookup("padding-bottom", "padding", zero).length.value

	l.dimensions.content.x = containBlock.content.x + l.dimensions.margin.left + l.dimensions.border.left + l.dimensions.padding.left

	l.dimensions.content.y = containBlock.content.height + containBlock.content.y + l.dimensions.margin.top + l.dimensions.border.top + l.dimensions.padding.top

}

//Height of a block-level non-replaced element in normal flow with overflow visible
func (l *LayoutBox) calculateBlockHeight(containBlock *Dimensions) {
	//If the height is set to an explicit length, use that exact lenght
	//Otherwise, just keep the value set by 'layoutBlockChildren'
	height := l.getStyleNode().value("height")
	if height.length.value != 0 {
		l.dimensions.content.height = height.length.value
	}
}

func GetTotalFrom(ml, mr, bl, br, pl, pr Value) float64 {
	return ml.length.value + mr.length.value + bl.length.value + br.length.value + pl.length.value + pr.length.value
}
