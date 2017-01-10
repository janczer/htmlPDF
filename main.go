package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
)

func main() {
	xmlFile, err := ioutil.ReadFile("test.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	xmlPure := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, string(xmlFile))
	r := bytes.NewReader([]byte(xmlPure))
	d := xml.NewDecoder(r)
	var t int = 0
	for {
		token, err := d.Token()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			panic(err)
		}
		switch token.(type) {
		case xml.StartElement:
			start := token.(xml.StartElement)
			tab(t)
			t++
			fmt.Printf("Name: %s, Attr: %v\n", start.Name, start.Attr)
		case xml.EndElement:
			t--
			tab(t)
			end := token.(xml.EndElement)
			fmt.Println(end)
		case xml.CharData:
			tab(t)
			text := token.(xml.CharData)
			fmt.Printf("%s\n", text)
		}

	}
}

func tab(i int) {
	for j := 0; j < i; j++ {
		fmt.Printf("\t")
	}
}
