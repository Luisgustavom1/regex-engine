package parser_test

import (
	"testing"

	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/parser"
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
			Name:  "a*b",
			Input: "a*b",
			Want:  "a*.b",
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

func TestShutingYardExp(t *testing.T) {
	tests := []struct {
		Input string
		Want  string
	}{
		{
			Input: "3+4",
			Want:  "34+",
		},
		{
			Input: "a*b",
			Want:  "ab*",
		},
		{
			Input: "a.b.c",
			Want:  "ab.c.",
		},
		{
			Input: "(a|b).c",
			Want:  "ab|c.",
		},
		{
			Input: "(a|b*c).c",
			Want:  "abc*|c.",
		},
		{
			Input: "(c*a|b).c",
			Want:  "ca*b|c.",
		},
		{
			Input: "(a|b|c)?c.c",
			Want:  "ab|c|c?c.",
		},
		{
			Input: "(a|b|c)+c.c",
			Want:  "ab|c|c+c.",
		},
		{
			Input: "(a|b+c)*c.c",
			Want:  "abc+|c*c.",
		},
		{
			Input: "(a|(b*c|d)*e)+c.c",
			Want:  "abc*d|e*|c+c.",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			r := parser.ToPostFixExp(tc.Input)
			if r != tc.Want {
				t.Fatalf("want %s, got %s", tc.Want, r)
			}
		})
	}
}
