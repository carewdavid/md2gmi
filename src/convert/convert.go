package convert

import (
	"bufio"
	"fmt"
	"net/url"
	"strings"
)

type linkParser struct {
	input *bufio.Scanner
	seen  strings.Builder
	links []string
}

func (p *linkParser) Scan() bool {
	return p.input.Scan()
}

func (p *linkParser) Next() string {
	return ""
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
	return markdown
}

// Format linkTarget and linkText as a gemtext link. The format is extremely simple:
// => gemini://example.com/ Example Text
func FormatLink(linkTarget *url.URL, linkText string) string {
	return fmt.Sprintf("=> %v %v\n", linkTarget.String(), linkText)
}
