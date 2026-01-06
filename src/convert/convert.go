package convert

import (
	"fmt"
	"net/url"
)

func Convert(markdown string) string {
	return markdown
}

func FormatLink(linkTarget *url.URL, linkText string) string {
	return fmt.Sprintf("=> %v %v\n", linkTarget.String(), linkText)
}
