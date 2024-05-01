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
			Name:  "a",
			Input: "a",
			Want:  "a",
		},
		{
			Name:  "abc",
			Input: "abc",
			Want:  "a.b.c",
		},
		{
			Name:  "(a|b)c",
			Input: "(a|b)c",
			Want:  "(a|b).c",
		},
		{
			Name:  "(a|b|c)?cc",
			Input: "(a|b|c)?cc",
			Want:  "(a|b|c)?c.c",
		},
		{
			Name:  "(a|b|c)+cc",
			Input: "(a|b|c)+cc",
			Want:  "(a|b|c)+c.c",
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
