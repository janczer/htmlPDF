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
	children   map[int]LayoutBox
}

func (l LayoutBox) getStyleNode() interface{} {
	boxType := l.box_type
	switch style := boxType.(type) {
	case BlockNode, InlineNode:
		return style
	case AnonymousBlock:
		panic("Anonymous block box has no style node")
	default:
		panic("Type must be BlockNode or InlineNode")
	}
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
		dimensions: Dimensions{},
		box_type:   boxType,
		children:   map[int]LayoutBox{},
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

func (l LayoutBox) getInlineContainer() LayoutBox {
	boxT := l.box_type
	switch boxT.(type) {
	case AnonymousBlock:
		return l
	case BlockNode:
		return NewLayoutBox(BlockNode{})
		//switch childBoxType := value.(l.children[len(l.children)-1].box_type) {
		//}
	default:
		return l
	}
}

func layoutTree(node StyleNode, containBlock Dimensions) LayoutBox {
	containBlock.content.height = 0

	rootBox := buildLayoutTree(node)
	rootBox.layout(containBlock)

	return rootBox
}

func (l LayoutBox) layout(containBlock Dimensions) {
	switch l.box_type.(type) {
	case BlockNode:
		l.layoutBox(containBlock)
	case InlineNode:
		//TODO
	case AnonymousBlock:
		//TODO
	}
}

func (l LayoutBox) layoutBox(containBlock Dimensions) {
	//Child width can depend on parent width, so we need to calculate
	//this box's width before laying out its children.
	l.calculateBlockWidth(containBlock)

	//Determine where the box is located within its container.
	l.calculateBlockPosition(containBlock)

	//Parent height can depend on child height, so calculateHeight
	//must be called after the children are laid out.
	l.calculateBlockHeight(containBlock)
}

func (l LayoutBox) calculateBlockWidth(containBlock Dimensions) {
	//style := l.getStyleNode()

	//width has initial value auto

	//margin, border, and padding have initial value 0
	//zero := Length{0.0, "px"}
}

func (l LayoutBox) calculateBlockPosition(containBlock Dimensions) {
}

func (l LayoutBox) calculateBlockHeight(containBlock Dimensions) {
}
