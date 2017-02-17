package htmlPDF

type StyleNode struct {
	node             *Node
	specified_values map[string]string
	children         map[int]*StyleNode
}

type MatchedRule struct {
}

func matchesSelector(elem ElementData, selector SimpleSelector) bool {

	//Check type selector
	if selector.tag_name != elem.tag_name {
		return false
	}

	//Check id
	if selector.id != elem.id() {
		return false
	}

	// Check class selectors
	//Check this loops
	//for _, selectorClass := range selector.class {
	//	for _, elemClass := range elem.classes() {
	//		if selectorClass == elemClass {
	//			return false
	//		}
	//	}
	//}

	return true
}
