package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	tt := map[string]struct {
		input string
		want  []string
	}{
		"simple":               {input: "  Hello   worlD  ", want: []string{"hello", "world"}},
		"mixed whitespacetype": {input: "\t hello \n world \t", want: []string{"hello", "world"}},
		"whitespace only":      {input: "    ", want: []string{}},
		"empty input":          {input: "", want: []string{}},
		"one word":             {input: "Test", want: []string{"test"}},
		"emojis":               {input: "🍎 🍌 🍊", want: []string{"🍎", "🍌", "🍊"}},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := cleanInput(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
