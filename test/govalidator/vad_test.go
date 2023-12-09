package govalidator

import (
	"github.com/asaskevich/govalidator"
	"testing"
)

type regexTest struct {
	a string
	b bool
}

func TestVad(t *testing.T) {
	pattern := `^[0-9]+$`
	tests := []regexTest{
		{"", false},
		{"11111", true},
		{"aaaaa", false},
		{"11111a", false},
		{"a11111", false},
		{"a11111a", false},
		{"aaaaaa", false},
		{"abc123abca1sad1sd13412", false},
	}
	for _, test := range tests {
		if got := govalidator.StringMatches(test.a, pattern); got != test.b {
			t.Errorf("govalidator.StringMatches(%v, %v) = %v, but expected = %v", test.a, pattern, got, test.b)
		}
	}
}
