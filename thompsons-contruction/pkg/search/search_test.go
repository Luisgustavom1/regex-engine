package search_test

import (
	"testing"

	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/nfa"
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/parser"
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/search"
)

func TestSearch(t *testing.T) {
	testCases := []struct {
		regex string
		world string
		want  bool
	}{
		{
			regex: "a*b",
			world: "",
			want:  false,
		},
		{
			regex: "a*b",
			world: "lkdlk",
			want:  false,
		},
		{
			regex: "a*b",
			world: "ab",
			want:  true,
		},
		{
			regex: "a*b",
			world: "aaaaaab",
			want:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.regex+" "+tc.world, func(t *testing.T) {
			t.Parallel()
			postfixExp := parser.ToPostFixExp(parser.InsertConcatOperator(tc.regex))
			n := nfa.ToNfa(postfixExp)
			r := search.Search(n, tc.world)
			if r != tc.want {
				t.Fatalf("want %v, got %v", tc.want, r)
			}
		})
	}
}
