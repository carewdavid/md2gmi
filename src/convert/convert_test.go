package convert

import "testing"

func TestConvertLine(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"The quick brown fox jumps over the lazy dog.", "The quick brown fox jumps over the lazy dog."},
		{"<https://example.com/>", "=> https://example.com/"},
		{"[This is a test link.](https://example.com)", "=> https://example.com This is a test link."},
		{"[Exemplar on Wikipedia](https://en.wikipedia.org/wiki/Exemplar_(disambiguation)", "=> https://en.wikipedia.org/wiki/Exemplar_(disambiguation) Exemplar on Wikipedia"},
		{"Test with [square brackets] but no link.", "Test with [square brackets] but no link."},
		{"This is a paragraph containing a [link](https://example.com) to somewhere else.", "This is a paragraph containing a link[*] to somewhere else.\n=> https://example.com link"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			result := ConvertLine(tt.in)
			if result != tt.out {
				t.Errorf("Got %q, expected %q\n", result, tt.out)
			}
		})
	}
}
