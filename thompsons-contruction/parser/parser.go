package parser

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

func pop(list *[]Operators) Operators {
	l := len(*list)
	last := (*list)[l-1]
	*list = (*list)[:l-1]
	return last
}

func peek(list []Operators) Operators {
	if len(list) == 0 {
		return 0
	}
	return list[len(list)-1]
}

func ShuntingYardExp(exp string) string {
	result := ""
	operators := make([]Operators, 0)

	for i := 0; i < len(exp); i++ {
		c := exp[i]
		co := Operators(c)

		if co == LEFT_PAREN || co == RIGHT_PAREN {
			if co == RIGHT_PAREN {
				for peek(operators) != LEFT_PAREN {
					result += string(pop(&operators))
				}
				// to remove )
				pop(&operators)
				continue
			}

			operators = append(operators, co)
			continue
		}

		if isSomeOperator(c) {
			for len(operators) > 0 && symbolPrecedence[co] <= symbolPrecedence[peek(operators)] && peek(operators) != LEFT_PAREN {
				result += string(pop(&operators))
			}

			operators = append(operators, co)
			continue
		}

		result += string(c)
	}

	for (len(operators)) > 0 {
		result += string(pop(&operators))
	}

	return result
}
