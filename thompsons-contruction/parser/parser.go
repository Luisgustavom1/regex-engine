package parser

type Symbols byte

const (
	CONCAT      Symbols = '*'
	ZERO_OR_ONE Symbols = '?'
	ONE_OR_MORE Symbols = '+'
	UNION       Symbols = '|'
	LEFT_PAREN  Symbols = '('
	RIGHT_PAREN Symbols = ')'
)

func isSomeSymbol(c byte) bool {
	s := Symbols(c)
	return s == CONCAT || s == ZERO_OR_ONE || s == ONE_OR_MORE || s == UNION || s == LEFT_PAREN || s == RIGHT_PAREN
}

// abc -> a.b.c
// (a|b)c -> (a|b).c
func InsertConcatOperator(exp string) string {
	expParsed := ""

	for i := 0; i < len(exp)-1; i++ {
		c := Symbols(exp[i])
		expParsed += string(c)

		if c == LEFT_PAREN || c == UNION || c == ZERO_OR_ONE || c == ONE_OR_MORE {
			continue
		}

		next := exp[i+1]
		if isSomeSymbol(next) {
			continue
		}

		expParsed += "."
	}
	expParsed += string(exp[len(exp)-1])

	return expParsed
}
