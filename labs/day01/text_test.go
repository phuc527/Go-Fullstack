package day01

import (
	"maps"
	"slices"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", ""},
		{"abc", "cba"},
		{"Việt Nam", "maN tệiV"},
		{"Go🚀", "🚀oG"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := Reverse(tt.in); got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestWordFreq(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want map[string]int
	}{
		{"empty", "", map[string]int{}},
		{"case and punctuation", "Go go GO! Node, node.", map[string]int{"go": 3, "node": 2}},
		{"unicode", "Việt Nam, việt", map[string]int{"việt": 2, "nam": 1}},
		{"numbers", "go 1.27 go", map[string]int{"go": 2, "1": 1, "27": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WordFreq(tt.in); !maps.Equal(got, tt.want) {
				t.Errorf("WordFreq(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestTopN(t *testing.T) {
	freq := map[string]int{"go": 3, "node": 2, "rust": 2, "java": 1}
	tests := []struct {
		name string
		n    int
		want []WordCount
	}{
		{"zero", 0, []WordCount{}},
		{"negative", -1, []WordCount{}},
		{"tie sorted alphabetically", 3, []WordCount{{"go", 3}, {"node", 2}, {"rust", 2}}},
		{"n larger than len", 10, []WordCount{{"go", 3}, {"node", 2}, {"rust", 2}, {"java", 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TopN(freq, tt.n); !slices.Equal(got, tt.want) {
				t.Errorf("TopN(n=%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}
