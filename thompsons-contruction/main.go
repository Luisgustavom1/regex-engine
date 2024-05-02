package main

import (
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/nfa"
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/parser"
)

func main() {
	regex := "a*b"
	word := "lklkl"

	postfixExp := parser.ToPostFixExp(parser.InsertConcatOperator(regex))
	n := nfa.ToNfa(postfixExp)

	search(n, word)
}
