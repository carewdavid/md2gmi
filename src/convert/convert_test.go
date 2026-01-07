package convert

import "testing"

func TestConvertLine(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"The quick brown fox jumps over the lazy dog.", "The quick brown fox jumps over the lazy dog."},
		{"<https://example.com/>", "=> https://example.com\n"},
		{"[This is a test link.](https://example.com)", "=> https://example.com This is a test link.\n"},
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
