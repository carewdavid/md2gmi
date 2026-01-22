package convert

import (
	"bufio"
	"fmt"
	"strings"
)

type linkParser struct {
	input           *bufio.Scanner
	seen            strings.Builder
	links           []link
	lineOnlyHasLink bool
}

type link struct {
	url  string
	text string
}

func newParser(input string) *linkParser {
	parser := linkParser{input: bufio.NewScanner(strings.NewReader(input)),
		seen:            strings.Builder{},
		links:           make([]link, 0),
		lineOnlyHasLink: true}
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
	//TODO: Handle code fences
	case "`":
		p.seen.WriteString(c)
		p.scanQuote()
	case "<":
		p.scanAngleLink()
	case "[":
		p.scanLinkOrFootnote()
	//TODO: also handle image links starting ![
	default:
		p.lineOnlyHasLink = false
		p.seen.WriteString(c)
	}
}

// Parse an inline quote in backticks
func (p *linkParser) scanQuote() {
	//All we really do is advance the parser to the next "`". Take care to include the backticks themselves in the output
	for p.input.Scan() {
		c := p.input.Text()
		p.seen.WriteString(c)
		if c == "`" {
			break
		}
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
		//Make sure to replace square brackets that might have been eaten by the parser
		p.seen.WriteString("[" + text.String() + "]")
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
	//Make sure we write the link text back to the source paragraph instead of snipping it out.
	//There are two areas for improvement:
	//1. We don't want to do this if the link is the only thing in the paragraph.
	//2. It could use a nicer indicator that there is a link attached to a part of the text. I.e. numbered footnotes.
	//Either or both of those may need some restructuring
	p.seen.WriteString(l.text + "[*]")
	p.links = append(p.links, l)

}

func Convert(markdown string) string {
	scanner := bufio.NewScanner(strings.NewReader(markdown))
	output := strings.Builder{}
	//Assume one paragraph per line
	//TODO: more sophisticated chunking of input text
	for scanner.Scan() {
		output.WriteString(ConvertLine(scanner.Text()) + "\n")
	}
	return output.String()
}

func ConvertLine(markdown string) string {
	parser := newParser(markdown)
	for parser.Scan() {
		parser.Next()
	}
	output := strings.Builder{}
	seenAsString := parser.seen.String()

	if !parser.lineOnlyHasLink {
		output.WriteString(seenAsString)
	}

	for i, link := range parser.links {
		if i == 0 && !parser.lineOnlyHasLink {
			output.WriteString("\n")
		}

		output.WriteString(link.String())
		if i == len(parser.links)-1 {
			continue
		} else {
			output.WriteString("\n")
		}
	}
	return output.String()
}

// Format linkTarget and linkText as a gemtext link. The format is extremely simple:
// => gemini://example.com/ Example Text
func FormatLink(linkTarget string, linkText string) string {
	link := fmt.Sprintf("=> %v %v", linkTarget, linkText)
	return strings.TrimSpace(link)
}
