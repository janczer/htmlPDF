package htmlPDF

import (
	"fmt"
	"regexp"
)

// Parse an HTML document and return the *Node
func ParseHtml(source string) *Node {
	p := new(Parser)
	p.input = source
	p.pos = 0
	nodes := p.parseNodes()
	if len(nodes) > 1 {
		panic("Not one root tag")
	}

	//return first node
	return nodes[0]
}

func (p *Parser) parseNodes() map[int]*Node {
	nodes := map[int]*Node{}
	for {
		p.consumeWhitespace()
		if p.eof() || p.startWith("</") {
			break
		}
		nodes[len(nodes)] = p.parseNode()
	}
	return nodes
}

//Parse a single node
func (p *Parser) parseNode() *Node {
	if p.nextChar() == "<" {
		return p.parseElement()
	} else {
		return p.parseText()
	}
}

//Parse a single element, including contents(and childrens if exist)
func (p *Parser) parseElement() *Node {
	//Opening tag
	start := p.consumeChar()
	if start != "<" {
		panic(fmt.Sprintf("%v was not an openig tag <", start))
	}

	tagName := p.parseTagName()
	attrs := p.parseAttributes()

	end := p.consumeChar()
	if end != ">" {
		panic(fmt.Sprintf("%v was not a closing tag <", end))
	}

	//Parse children
	children := p.parseNodes()

	//Closing tag
	start = p.consumeChar()
	if start != "<" {
		panic(fmt.Sprintf("%v was not an openig tag <", start))
	}

	slash := p.consumeChar()
	if slash != "/" {
		panic(fmt.Sprintf("%v was not a tag /", slash))
	}

	closeName := p.parseTagName()
	if closeName != tagName {
		panic(fmt.Sprintf("open tag %v and close tag %v don't equal ", tagName, closeName))
	}

	end = p.consumeChar()
	if end != ">" {
		panic(fmt.Sprintf("%v was not a closing tag <", end))
	}

	return elem(tagName, attrs, children)
}

//Parse a text node
func (p *Parser) parseText() *Node {
	return text(p.consumeWhile(func(char string) bool {
		return char != "<"
	}))
}

//Parse a tag or attribute name
func (p *Parser) parseTagName() string {
	reg := regexp.MustCompile("[a-zA-Z0-9]")
	f := func(char string) bool {
		return reg.MatchString(char)
	}
	return p.consumeWhile(f)
}

//Parse a list of name="value" pairs
func (p *Parser) parseAttributes() map[string]string {
	attr := map[string]string{}

	for {
		p.consumeWhitespace()
		if p.nextChar() == ">" {
			break
		}
		name, value := p.parseAttribute()
		attr[name] = value
	}

	return attr
}

//Parse a single name="value" pair
func (p *Parser) parseAttribute() (string, string) {
	name := p.parseTagName()
	delimiter := p.consumeChar()

	if delimiter != "=" {
		panic(fmt.Sprintf("%v was not =", delimiter))
	}
	value := p.parseAttributeValue()
	return name, value
}

//Parse a quoted value
func (p *Parser) parseAttributeValue() string {
	q := p.consumeChar()
	if q != "\"" && q != "'" {
		panic(fmt.Sprintf("%v was not \" or '", q))
	}

	value := p.consumeWhile(func(char string) bool {
		return char != q
	})

	cq := p.consumeChar()
	if cq != q {
		panic(fmt.Sprintf("%v was not %v", cq, q))
	}

	return value
}

//Return true if current input start with the given string
func (p *Parser) startWith(test string) bool {
	start := true
	for i := 0; i < len(test); i++ {
		if p.input[p.pos+i] != test[i] {
			start = false
		}
	}

	return start
}

//Consume characters until function 'test' returns false
func (p *Parser) consumeWhile(test func(char string) bool) string {
	var result string
	for {
		if p.eof() || !test(p.nextChar()) {
			break
		}
		result += p.consumeChar()
	}

	return result
}

//Consume and discard zero or more whitespace characters
func (p *Parser) consumeWhitespace() {
	reg := regexp.MustCompile("[\\s]")
	f := func(char string) bool {
		return reg.MatchString(char)
	}

	p.consumeWhile(f)
}

//Return the current character with consuming it
func (p *Parser) consumeChar() string {
	char := p.input[p.pos]
	p.pos++
	return string(char)
}

//Read the current character without consuming it
func (p *Parser) nextChar() string {
	return string(p.input[p.pos])
}

//Return true if all input is consumed
func (p *Parser) eof() bool {
	return p.pos >= len(p.input)
}
