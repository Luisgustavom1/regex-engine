package parser_test

import (
	"testing"

	"github.com/Luisgustavom1/regex-engine/thompsons-construction/parser"
)

func TestInsertConcatOperator(t *testing.T) {
	tests := []struct {
		Name  string
		Input string
		Want  string
	}{
		{
			Name:  "abc",
			Input: "abc",
			Want:  "a.b.c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			r := parser.InsertConcatOperator(tc.Input)
			if r != tc.Want {
				t.Fatalf("want %s, got %s", tc.Want, r)
			}
		})
	}
}
