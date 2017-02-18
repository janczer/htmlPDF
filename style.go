package htmlPDF

import (
	"sort"
)

type StyleNode struct {
	node             *Node
	specified_values map[string]Value
	children         map[int]StyleNode
}

type MatchedRule struct {
	spec Specificity
	rule *Rule
}

//If rule match elem, return a MatchedRule
func matchRule(elem *ElementData, rule *Rule) MatchedRule {
	for _, selector := range rule.selectors {
		if matchesSelector(elem, selector) {
			//Fine the first (highest-specificity) matching selector
			mr := MatchedRule{
				selector.specificity(),
				rule,
			}
			return mr
		}
	}

	return MatchedRule{}
}

func matchesSelector(elem *ElementData, selector SimpleSelector) bool {
	//Check type selector
	if selector.tag_name != "" && selector.tag_name != elem.tag_name {
		return false
	}

	//Check id
	if selector.id != "" && selector.id != elem.id() {
		return false
	}

	// Check class selectors
	if !elem.classContains(selector.class) {
		return false
	}

	return true
}

//Find all CSS rules that match the given element
func matchingRules(elem *ElementData, stylesheet *Stylesheet) map[int]MatchedRule {
	matched := map[int]MatchedRule{}

	for i, rule := range stylesheet.rules {
		mr := matchRule(elem, rule)
		if mr.rule != nil {
			matched[i] = mr
		}
	}
	return matched
}

type SortBySpec map[int]MatchedRule

func (a SortBySpec) Len() int           { return len(a) }
func (a SortBySpec) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortBySpec) Less(i, j int) bool { return a[i].spec.a < a[j].spec.a }

func specifiedValues(elem *ElementData, stylesheet *Stylesheet) map[string]Value {
	values := map[string]Value{}
	rules := matchingRules(elem, stylesheet)

	//add sort rules
	sort.Sort(SortBySpec(rules))

	for _, matchedRule := range rules {
		for _, declaration := range matchedRule.rule.declaration {
			values[declaration.name] = declaration.value
		}
	}

	return values
}

func styleTree(root *Node, stylesheet *Stylesheet) StyleNode {
	children := map[int]StyleNode{}
	for i, child := range root.children {
		children[i] = styleTree(child, stylesheet)
	}

	specifiedValue := map[string]Value{}
	if root.node_type.element.tag_name != "" {
		specifiedValue = specifiedValues(&root.node_type.element, stylesheet)
	}

	return StyleNode{
		node:             root,
		specified_values: specifiedValue,
		children:         children,
	}
}

//Return true if ElementData contain one or more class
func (e ElementData) classContains(class map[int]string) bool {
	if len(class) == 0 {
		return true
	}
	for _, class := range class {
		for _, eclass := range e.classes() {
			if class == eclass {
				return true
			}
		}
	}

	return false
}
