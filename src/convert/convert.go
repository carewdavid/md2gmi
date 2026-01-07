package convert

import (
	"bufio"
	"fmt"
	"strings"
)

type linkParser struct {
	input *bufio.Scanner
	seen  strings.Builder
	links []link
}

type link struct {
	url  string
	text string
}

func newParser(input string) *linkParser {
	parser := linkParser{input: bufio.NewScanner(strings.NewReader(input)),
		seen:  strings.Builder{},
		links: make([]link, 0)}
	parser.input.Split(bufio.ScanRunes)
	return &parser
}

func (l link) String() string {
	return FormatLink(l.url, l.text)
}

func (p *linkParser) Scan() bool {
	return p.input.Scan()
}

func (p *linkParser) Next() {
	c := p.input.Text()
	switch c {
	case "<":
		p.scanAngleLink()
	case "[":
		p.scanLinkOrFootnote()
	//TODO: also handle image links starting ![
	default:
		p.seen.WriteString(c)
	}
}

// Parse a bare link like <https://example.com>
func (p *linkParser) scanAngleLink() {
	url := strings.Builder{}
	for p.input.Scan() {
		c := p.input.Text()
		if c == ">" {
			break
		} else {
			url.WriteString(c)
		}
	}
	l := link{url: url.String(), text: ""}
	p.links = append(p.links, l)
}

func (p *linkParser) scanLinkOrFootnote() {
	url := strings.Builder{}
	text := strings.Builder{}
	//Scan link text
	//TODO: handle footnotes
	for p.input.Scan() {
		c := p.input.Text()
		//TODO: look up whether markdown allows nested [] in link text
		if c == "]" {
			break
		} else {
			text.WriteString(c)
		}
	}
	//Scan url

	//First make sure we actually do have a url coming
	if !p.input.Scan() {
		p.seen.WriteString(text.String())
		return
	}

	c := p.input.Text()
	if c != "(" {
		p.seen.WriteString(text.String())
		p.seen.WriteString(c)
		return
	}
	parenDepth := 0
	for p.input.Scan() {
		c := p.input.Text()
		//This would be so much simpler if urls couldn't have parentheses in them...
		if c == "(" {
			parenDepth += 1
			url.WriteString(c)
		} else if c == ")" {
			if parenDepth == 0 {
				break
			} else {
				parenDepth -= 1
				url.WriteString(c)
			}
		} else {
			url.WriteString(c)
		}
	}
	l := link{url: url.String(), text: text.String()}
	p.links = append(p.links, l)

}

func Convert(markdown string) string {
	scanner := bufio.NewScanner(strings.NewReader(markdown))
	output := strings.Builder{}
	//Assume one paragraph per line
	//TODO: more sophisticated chunking of input text
	for scanner.Scan() {
		output.WriteString(ConvertLine(scanner.Text()))
	}
	return output.String()
}

func ConvertLine(markdown string) string {
	parser := newParser(markdown)
	for parser.Scan() {
		parser.Next()
	}
	output := strings.Builder{}
	output.WriteString(parser.seen.String())
	for _, link := range parser.links {
		output.WriteString(link.String() + "\n")
	}
	return output.String()
}

// Format linkTarget and linkText as a gemtext link. The format is extremely simple:
// => gemini://example.com/ Example Text
func FormatLink(linkTarget string, linkText string) string {
	link := fmt.Sprintf("=> %v %v", linkTarget, linkText)
	return strings.TrimSpace(link)
}
