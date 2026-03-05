package scraper

import (
	"testing"
)

func TestStateName(t *testing.T) {
	tests := []struct {
		code string
		want string
	}{
		{"AL", "Alabama"},
		{"KS", "Kansas"},
		{"MO", "Missouri"},
		{"TX", "Texas"},
		{"CA", "California"},
		{"XX", ""},
	}
	for _, tt := range tests {
		got := StateName(tt.code)
		if got != tt.want {
			t.Errorf("StateName(%q) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input string
		want  float64
		ok    bool
	}{
		{"43,830", 43830, true},
		{"2,091,480", 2091480, true},
		{"120000", 120000, true},
		{"0", 0, true},
		{"*", 0, false},
		{"#", 0, false},
		{"**", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseNumber(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("ParseNumber(%q) = (%v, %v), want (%v, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input string
		want  int
		ok    bool
	}{
		{"2,091,480", 2091480, true},
		{"830", 830, true},
		{"*", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseInt(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("ParseInt(%q) = (%v, %v), want (%v, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestShouldSkipOccCode(t *testing.T) {
	tests := []struct {
		code string
		skip bool
	}{
		{"00-0000", true},
		{"11-0000", true},
		{"15-0000", true},
		{"15-1254", false},
		{"11-1011", false},
		{"11-1021", false},
	}
	for _, tt := range tests {
		got := ShouldSkipOccCode(tt.code)
		if got != tt.skip {
			t.Errorf("ShouldSkipOccCode(%q) = %v, want %v", tt.code, got, tt.skip)
		}
	}
}