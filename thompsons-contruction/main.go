package main

import (
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/nfa"
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/parser"
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/search"
)

func main() {
	regex := "a*b"
	word := "lklkl"

	postfixExp := parser.ToPostFixExp(parser.InsertConcatOperator(regex))
	n := nfa.ToNfa(postfixExp)

	match := search.Search(n, word)
	if match {
		println("Match")
	} else {
		println("No match")
	}
}
