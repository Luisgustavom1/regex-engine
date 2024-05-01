package parser

import "github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/ds"

type Operators byte

const (
	CONCAT      Operators = '*'
	DOT         Operators = '.'
	ZERO_OR_ONE Operators = '?'
	ONE_OR_MORE Operators = '+'
	UNION       Operators = '|'
	LEFT_PAREN  Operators = '('
	RIGHT_PAREN Operators = ')'
)

type OperatorPrecedence int

const (
	Union         OperatorPrecedence = 0
	Dot                              = 1
	ZeroOrOne                        = 2
	Concatenation                    = 2
	OneOrMore                        = 2
)

var symbolPrecedence = map[Operators]OperatorPrecedence{
	UNION:       Union,
	DOT:         Dot,
	ZERO_OR_ONE: ZeroOrOne,
	CONCAT:      Concatenation,
	ONE_OR_MORE: OneOrMore,
}

func isSomeOperator(c byte) bool {
	s := Operators(c)
	return s == CONCAT || s == DOT || s == ZERO_OR_ONE || s == ONE_OR_MORE || s == UNION || s == LEFT_PAREN || s == RIGHT_PAREN
}

// abc -> a.b.c
// (a|b)c -> (a|b).c
func InsertConcatOperator(exp string) string {
	expParsed := ""

	for i := 0; i < len(exp)-1; i++ {
		c := Operators(exp[i])
		expParsed += string(c)

		if c == LEFT_PAREN || c == UNION || c == ZERO_OR_ONE || c == ONE_OR_MORE {
			continue
		}

		next := exp[i+1]
		if isSomeOperator(next) {
			continue
		}

		expParsed += "."
	}
	expParsed += string(exp[len(exp)-1])

	return expParsed
}

func ShuntingYardExp(exp string) string {
	result := ""
	operators := ds.NewStack[Operators]()

	for i := 0; i < len(exp); i++ {
		c := exp[i]
		co := Operators(c)

		if co == LEFT_PAREN || co == RIGHT_PAREN {
			if co == RIGHT_PAREN {
				for operators.Peek() != LEFT_PAREN {
					result += string(operators.Pop())
				}
				// to remove )
				operators.Pop()
				continue
			}

			operators.Push(co)
			continue
		}

		if isSomeOperator(c) {
			for operators.Len() > 0 && symbolPrecedence[co] <= symbolPrecedence[operators.Peek()] && operators.Peek() != LEFT_PAREN {
				result += string(operators.Pop())
			}

			operators.Push(co)
			continue
		}

		result += string(c)
	}

	for (len(operators.Values())) > 0 {
		result += string(operators.Pop())
	}

	return result
}
