package htmlPDF

import (
	"regexp"
	"strconv"
)

type Stylesheet struct {
	rules map[int]*Rule
}

type Rule struct {
	selectors   map[int]SimpleSelector
	declaration map[int]Declaration
}

type SimpleSelector struct {
	tag_name string
	id       string
	class    map[int]string
}

type Specificity struct {
	a int
	b int
	c int
}

//Calculate specificity
//https://www.w3.org/TR/selectors/#specificity
//cahnge algorithms
func (s SimpleSelector) specificity() Specificity {
	var a, b, c int
	if len(s.id) > 0 {
		a++
	}
	if len(s.class) > 0 {
		b = len(s.class)
	}
	if len(s.tag_name) > 0 {
		c++
	}

	return Specificity{a, b, c}
}

type Declaration struct {
	name  string
	value Value
}

type Value struct {
	keyword string
	length  Length
	color   Color
}

type Length struct {
	value float64
	unit  string //only px
}

type Color struct {
	r uint
	g uint
	b uint
	a uint
}

type Parser struct {
	pos   int
	input string
}

func validLengthChar(c string) bool {
	var valid = regexp.MustCompile("[0-9.]")
	return valid.MatchString(c)
}

//Parse a whole CSS stylesheet
func CssParser(source string) *Parser {
	return &Parser{
		pos:   0,
		input: source,
	}
}

//Parse a list of rule sets, separated by optional whitespace
func (p *Parser) parseRules() Stylesheet {
	rules := map[int]*Rule{}

	for {
		p.consumeWhitespace()
		if p.eof() {
			break
		}
		rules[len(rules)] = p.parseRule()
	}

	return Stylesheet{rules}
}

//Parse a rule: 'selectors { declarations }'
//declarations it's pair of 'property: value;'
func (p *Parser) parseRule() *Rule {
	return &Rule{
		selectors:   p.parseSelectors(),
		declaration: p.parseDeclarations(),
	}
}

//Parse a list of declarations enclosed in '{ ... }'
func (p *Parser) parseDeclarations() map[int]Declaration {
	p.consumeChar()
	decl := map[int]Declaration{}
	for {
		p.consumeWhitespace()
		if p.nextChar() == "}" {
			p.consumeChar()
			break
		}
		decl[len(decl)] = p.parseDeclaration()
	}

	return decl
}

// Parse one declaration pair: 'property: value;'
func (p *Parser) parseDeclaration() Declaration {
	name := p.parseIdentifier()
	p.consumeWhitespace()
	for p.consumeChar() != ":" {
	}
	p.consumeWhitespace()
	value := p.parseValue()
	p.consumeWhitespace()
	for p.consumeChar() != ";" {
	}

	return Declaration{name, value}
}

//Parse value
func (p *Parser) parseValue() Value {
	var valid = regexp.MustCompile("[0-9]")

	switch {
	case p.nextChar() == "#":
		p.consumeChar()
		return p.parseColor()
	case valid.MatchString(string(p.nextChar())):
		return p.parseLength()
	default:
		return Value{keyword: p.parseIdentifier()}
	}
}

func (p *Parser) parseLength() Value {
	return Value{length: p.parseFloat()}
}

//Parse value 000px, support only px
func (p *Parser) parseFloat() Length {
	var result string
	for !p.eof() && validLengthChar(p.nextChar()) {
		result += string(p.consumeChar())
	}
	r, e := strconv.ParseFloat(result, 64)
	if e != nil {
		return Length{0, "px"}
	}
	return Length{r, "px"}
}

//Parse color #000000
func (p *Parser) parseColor() Value {
	return Value{
		color: Color{
			r: p.parseHexPair(),
			g: p.parseHexPair(),
			b: p.parseHexPair(),
			a: 255,
		},
	}
}

//Parse two hexadecimal digits
func (p *Parser) parseHexPair() uint {
	s := p.input[p.pos : p.pos+2]
	p.pos += 2
	r, e := strconv.ParseUint(s, 16, 64)
	if e != nil {
		return 0
	}
	return uint(r)
}

//Parse a comma-separated list of selectors
func (p *Parser) parseSelectors() map[int]SimpleSelector {
	s := map[int]SimpleSelector{}
Loopsels:
	for {
		s[len(s)] = p.parseSelector()
		p.consumeWhitespace()
		switch p.nextChar() {
		case ",":
			p.consumeChar()
			p.consumeWhitespace()
		case "{":
			break Loopsels
		default:
			panic("Unexpected character")
		}
	}

	return s
}

//Parse one simple selector, e.g.: '#id, class1, class2, class3'
func (p *Parser) parseSelector() SimpleSelector {
	m := SimpleSelector{class: map[int]string{}}
Loopsel:
	for !p.eof() {
		c := p.nextChar()
		switch {
		case c == "#":
			p.consumeChar()
			m.id = p.parseIdentifier()
		case c == ".":
			p.consumeChar()
			m.class[len(m.class)] = p.parseIdentifier()
		case c == "*":
			// universal selector
			p.consumeChar()
		case validIdentifierChar(c):
			m.tag_name = p.parseIdentifier()
		default:
			break Loopsel
		}
		p.consumeWhitespace()
	}
	return m
}

func validIdentifierChar(c string) bool {
	var valid = regexp.MustCompile("[a-zA-Z0-9-_]")
	return valid.MatchString(c)
}

//Parse a property name or keyword
func (p *Parser) parseIdentifier() string {
	var valid = regexp.MustCompile("[a-zA-Z0-9-_]")
	return p.consumeWhile(func(char string) bool {
		return valid.MatchString(char)
	})
}
