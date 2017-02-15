package htmlPDF

import (
	"fmt"
	"regexp"
	"strconv"
)

type Stylesheet struct {
	rules map[int]Rule
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

func NewParser(source string) *Parser {
	return &Parser{
		pos:   0,
		input: source,
	}
}

func (p *Parser) parseRules() Stylesheet {
	var css Stylesheet

	rules := map[int]Rule{}

	for i := 0; i < 3; i++ {
		p.consumeWhitespace()
		if p.eof() {
			break
		}
		rules[len(rules)] = p.parseRule()
	}
	fmt.Println(rules)
	return css
}

func (p *Parser) parseRule() Rule {
	return Rule{
		selectors:   p.parseSelectors(),
		declaration: p.parseDeclarations(),
	}
}

func (p *Parser) parseDeclarations() map[int]Declaration {
	decl := map[int]Declaration{}
	for {
		p.consumeWhitespace()
		if p.nextChar() == '}' {
			p.consumeChar()
			break
		}
		decl[len(decl)] = p.parseDeclaration()
	}

	return decl
}

func (p *Parser) parseDeclaration() Declaration {
	name := p.parseIdentifier()
	p.consumeWhitespace()
	for p.consumeChar() != ':' {
	}
	p.consumeWhitespace()
	value := p.parseValue()
	p.consumeWhitespace()
	for p.consumeChar() != ';' {
	}

	return Declaration{name, value}
}

//add parse number
func (p *Parser) parseValue() Value {
	var valid = regexp.MustCompile("[0-9]")

	switch {
	case p.nextChar() == '#':
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

func (p *Parser) parseFloat() Length {
	var result string
	for !p.eof() && validLengthChar(p.nextChar()) {
		result += string(p.consumeChar())
	}
	r, e := strconv.ParseFloat(result, 64)
	if e != nil {
		fmt.Println(e)
		return Length{0, "px"}
	}
	return Length{r, "px"}
}

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

func (p *Parser) parseHexPair() uint {
	s := p.input[p.pos : p.pos+2]
	p.pos += 2
	r, e := strconv.ParseUint(s, 16, 64)
	if e != nil {
		fmt.Println(e)
		return 0
	}
	return uint(r)
}

func (p *Parser) parseSelectors() map[int]SimpleSelector {
	s := map[int]SimpleSelector{}
Loopsels:
	for {
		s[len(s)] = p.parseSelector()
		p.consumeWhitespace()
		switch p.nextChar() {
		case ',':
			p.consumeChar()
			p.consumeWhitespace()
		case '{':
			break Loopsels
		default:
			fmt.Printf("char %v\n", string(p.nextChar()))
			panic("Unexpected character")
		}
	}

	return s
}

func (p *Parser) parseSelector() SimpleSelector {
	m := SimpleSelector{class: map[int]string{}}
Loopsel:
	for !p.eof() {
		switch p.nextChar() {
		case '#':
			p.consumeChar()
			m.id = p.parseIdentifier()
		case '.':
			p.consumeChar()
			m.class[len(m.class)] = p.parseIdentifier()
		case '*':
			// universal selector
			p.consumeChar()
		default:
			break Loopsel
		}
		p.consumeWhitespace()
	}
	return m
}

//Return the current character without consuming it
func (p *Parser) nextChar() byte {
	return p.input[p.pos]
}

func (p *Parser) parseIdentifier() string {
	return p.consumeWhile()
}

func (p *Parser) consumeChar() byte {
	r := p.nextChar()
	p.pos++
	return r
}

func (p *Parser) consumeWhile() string {
	var result string
	for !p.eof() && validIdentifierChar(p.nextChar()) {
		result += string(p.consumeChar())
	}
	return result
}

// Consumed whitespaces
func (p *Parser) consumeWhitespace() {
	var valid = regexp.MustCompile("\\s")
	for !p.eof() && valid.MatchString(string(p.nextChar())) {
		p.consumeChar()
	}
}

// Return true if all input is consumed.
func (p *Parser) eof() bool {
	return p.pos >= len(p.input)
}

func main() {
	fmt.Println("css engine")

	//	css := "#test .first .third, #secid .second, #thirdid { margin: auto; color: #cc0000; }"
	css2 := `
	#test .class .class1, .class2 {
		color: #cc00bb;
		margin: auto;
		padding-top: 100.14px;
	}
	#test2 .class {
		color: #cc00bb;
	}
	`
	fmt.Println(css2)
	p := NewParser(css2)
	p.parseRules()
}

func validIdentifierChar(c byte) bool {
	var valid = regexp.MustCompile("[a-zA-Z0-9-_]")
	return valid.MatchString(string(c))
}

func validLengthChar(c byte) bool {
	var valid = regexp.MustCompile("[0-9.]")
	return valid.MatchString(string(c))
}
