// xmlselect выводит текст выбранных элеменьлв XML-документа.
// Поддерживает селекторы: div, div[class=title], div[class=main][id=top]
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

type Selector struct {
	Name string            // имя тега
	Attr map[string]string // aттрибуты
}

func parseselector(s string) (Selector, error) {
	sel := Selector{
		Attr: make(map[string]string),
	}

	bracket := strings.Index(s, "[")
	if bracket == -1 {
		sel.Name = s
		return sel, nil
	}

	sel.Name = s[:bracket]

	rest := s[bracket:]

	for len(rest) > 0 {
		if rest[0] != '[' {
			return sel, fmt.Errorf("ожидалось '[', получено %q", rest[0])
		}
		end := strings.Index(rest, "]")
		if end == -1 {
			return sel, fmt.Errorf("незакрыьая скобка в %q", rest)
		}
		attrPart := rest[1:end]
		eq := strings.Index(attrPart, "=")
		if eq == -1 {
			return sel, fmt.Errorf("ожидался '=', получено %q", attrPart)
		}
		key := attrPart[:eq]
		value := attrPart[eq+1:]
		sel.Attr[key] = value
		rest = rest[end+1:]
	}
	return sel, nil
}

func matches(start xml.StartElement, sel Selector) bool {
	//	fmt.Fprintf(os.Stderr, "DEBUG: mathes: element=%q, selector=%q, attrs=%v\n", start.Name.Local, sel.Name, sel.Attr)
	if start.Name.Local != sel.Name {
		//	fmt.Fprintf(os.Stderr, "DEBUG matches: имя не совпалает (%q vs %q)\n", start.Name.Local, sel.Name)
		return false
	}

	for key, wantVal := range sel.Attr {
		//		fmt.Fprintf(os.Stderr, "DEBUG matches: проверяем атрибут %q = %q\n", key, wantVal)
		found := false
		for _, attr := range start.Attr {
			if attr.Name.Local == key {
				if attr.Value == wantVal {
					found = true
					//					fmt.Fprintf(os.Stderr, "DEBUG matches: атрибут %q найден со значением %q\n", key, attr.Value)
					break
				}
			}
		}
		if !found {
			//			fmt.Fprintf(os.Stderr, "DEBUG matches: атрибут %q не найден\n", key)
			return false
		}
	}
	//	fmt.Fprintf(os.Stderr, "DEBUG mawtches: УСПЕХ\n")
	return true
}

func containsAll(stack []string, selectors []Selector) bool {
	x := stack
	y := make([]string, len(selectors))
	for i, sel := range selectors {
		y[i] = sel.Name
	}
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

func matchesAll(stack []string, elemStack []xml.StartElement, selectors []Selector) bool {
	if len(selectors) > len(stack) {
		return false
	}
	selIdx := 0
	for i := 0; i < len(stack) && selIdx < len(selectors); i++ {
		if stack[i] == selectors[selIdx].Name {
			if !matches(elemStack[i], selectors[selIdx]) {
				return false
			}
			selIdx++
		}
	}
	return selIdx == len(selectors)
}

func main() {
	selectors := make([]Selector, 0, len(os.Args)-1)
	for _, arg := range os.Args[1:] {
		sel, err := parseselector(arg)
		if err != nil {
			//		fmt.Fprintf(os.Stderr, "xmlselect: ошибка в селекторе %q: %v\n", arg, err)
			os.Exit(1)
		}
		selectors = append(selectors, sel)
	}

	//	fmt.Fprintf(os.Stderr, "DEBUG: ceлекторы: %+v\n", selectors)

	dec := xml.NewDecoder(os.Stdin)
	var stack []string
	var elemStack []xml.StartElement

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok.Name.Local)
			elemStack = append(elemStack, tok)
			//		fmt.Fprintf(os.Stderr, "DEBUG: start %s, stack=%v\n", tok.Name.Local, stack)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
			elemStack = elemStack[:len(elemStack)-1]
		case xml.CharData:
			//		fmt.Fprintf(os.Stderr, "DEBUG: char data %q, stack=%v\n", tok, stack)
			if matchesAll(stack, elemStack, selectors) && strings.TrimSpace(string(tok)) != "" {
				fmt.Printf("%s: %s\n", strings.Join(stack, " "), tok)
			}
		}
	}
}
