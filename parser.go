package main

import (
	"fmt"
	"strconv"
	"strings"
)

type tokenType uint8

const (
	group           tokenType = iota
	bracket         tokenType = iota
	or              tokenType = iota
	repeat          tokenType = iota
	literal         tokenType = iota
	groupUncaptured tokenType = iota
)

type token struct {
	tokenType tokenType
	value     interface{}
}

type repeatPayload struct {
	min   int
	max   int
	token token
}

const repeatInfinity = -1

type parseContext struct {
	pos    int
	tokens []token
}

func parse(regex string) *parseContext {
	ctx := &parseContext{
		pos:    0,
		tokens: []token{},
	}

	for ctx.pos < len(regex) {
		process(regex, ctx)
		ctx.pos++
	}

	return ctx
}

func process(regex string, ctx *parseContext) {
	ch := regex[ctx.pos]

	switch ch {
	case '(':
		groupCtx := &parseContext{
			pos:    ctx.pos,
			tokens: []token{},
		}
		parseGroup(regex, groupCtx)
		ctx.tokens = append(
			ctx.tokens,
			token{
				tokenType: group,
				value:     groupCtx.tokens,
			},
		)
	case '[':
		parseBracket(regex, ctx)
	case '|':
		parseOr(regex, ctx)
	case '*':
	case '+':
	case '?':
		parseRepeat(regex, ctx)
	case '{':
		parseRepeatSpecified(regex, ctx)
	default:
		parseLiteral(regex, ctx)
	}
}

func parseGroup(regex string, ctx *parseContext) {
	ctx.pos++
	for regex[ctx.pos] != ')' {
		process(regex, ctx)
		ctx.pos++
	}
}

func parseBracket(regex string, ctx *parseContext) {
	ctx.pos++
	literals := []string{}

	for regex[ctx.pos] != ']' {
		ch := regex[ctx.pos]
		if ch == '-' {
			next := regex[ctx.pos+1]
			prev := literals[len(literals)-1][0]
			literals = append(literals, fmt.Sprintf("%c%c", prev, next))
		} else {
			literals = append(literals, fmt.Sprintf("%c", ch))
		}
		ctx.pos++
	}

	literalsSet := map[uint8]bool{}

	for _, l := range literals {
		for i := l[0]; i <= l[len(l)-1]; i++ {
			literalsSet[i] = true
		}
	}

	ctx.tokens = append(ctx.tokens, token{
		tokenType: bracket,
		value:     literalsSet,
	})
}

func parseOr(regex string, ctx *parseContext) {
	// right hand side
	rhsCtx := &parseContext{
		pos:    ctx.pos + 1,
		tokens: []token{},
	}

	for rhsCtx.pos < len(regex) && regex[rhsCtx.pos] != ')' {
		process(regex, rhsCtx)
		rhsCtx.pos++
	}

	left := token{
		tokenType: groupUncaptured,
		value:     ctx.tokens,
	}

	right := token{
		tokenType: groupUncaptured,
		value:     rhsCtx.tokens,
	}

	ctx.pos = rhsCtx.pos
	ctx.tokens = []token{{
		tokenType: or,
		value:     []token{left, right},
	}}
}

func parseLiteral(regex string, ctx *parseContext) {
	t := token{
		tokenType: literal,
		value:     regex[ctx.pos],
	}
	ctx.tokens = append(ctx.tokens, t)
}

func parseRepeat(regex string, ctx *parseContext) {
	ch := regex[ctx.pos]
	var min, max int

	switch ch {
	case '*':
		min = 0
		max = repeatInfinity
	case '+':
		min = 1
		max = repeatInfinity
	case '?':
		min = 0
		max = 1
	}

	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens[len(ctx.tokens)-1] = token{
		tokenType: repeat,
		value: repeatPayload{
			min:   min,
			max:   max,
			token: lastToken,
		},
	}
}

func parseRepeatSpecified(regex string, ctx *parseContext) {
	start := ctx.pos + 1
	for regex[ctx.pos] != '}' {
		ctx.pos++
	}

	specifiedStr := regex[start:ctx.pos]
	limits := strings.Split(specifiedStr, ",")
	if len(limits) > 2 {
		panic("repeat specified should have at most 2 limits")
	}

	var min, max int
	if len(limits) == 1 {
		value, err := strconv.Atoi(limits[0])
		if err != nil {
			panic(err.Error())
		}
		min = value
		max = value
	}

	if len(limits) == 2 {
		firstLimit, err := strconv.Atoi(limits[0])
		if err != nil {
			panic(err.Error())
		}
		min = firstLimit

		if limits[1] == "" {
			max = repeatInfinity
		} else {
			secondLimit, err := strconv.Atoi(limits[1])
			if err != nil {
				panic(err.Error())
			}
			max = secondLimit
		}
	}

	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens[len(ctx.tokens)-1] = token{
		tokenType: repeat,
		value: repeatPayload{
			min:   min,
			max:   max,
			token: lastToken,
		},
	}
}
