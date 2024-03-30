package main

type state struct {
	start       bool
	terminal    bool
	transitions stateTransition
}

type stateTransition = map[uint8][]*state

const epsilonChar uint8 = 0 // empty character

func toNfa(ctx *parseContext) *state {
	startState, endState := tokenToNfa(&ctx.tokens[0])

	for _, t := range ctx.tokens[1:] {
		startNext, endNext := tokenToNfa(&t)
		endState.transitions[epsilonChar] = append(
			endState.transitions[epsilonChar],
			startNext,
		)
		endState = endNext
	}

	start := &state{
		start: true,
		transitions: stateTransition{
			epsilonChar: []*state{startState},
		},
	}

	end := &state{
		terminal:    true,
		transitions: stateTransition{},
	}

	endState.transitions[epsilonChar] = append(
		endState.transitions[epsilonChar],
		end,
	)

	return start
}

func tokenToNfa(t *token) (start *state, end *state) {
	start = &state{
		transitions: stateTransition{},
	}
	end = &state{
		transitions: stateTransition{},
	}

	switch t.tokenType {
	case literal:
		ch := t.value.(uint8)
		start.transitions[ch] = []*state{end}
	case or:
		values := t.value.([]token)
		left := values[0]
		right := values[1]

		startLeft, endLeft := tokenToNfa(&left)
		startRight, endRight := tokenToNfa(&right)

		start.transitions[epsilonChar] = []*state{startLeft, startRight}
		endLeft.transitions[epsilonChar] = []*state{end}
		endRight.transitions[epsilonChar] = []*state{end}
	case bracket:
		literalsSet := t.value.(map[uint8]bool)
		for ch := range literalsSet {
			start.transitions[ch] = []*state{end}
		}
	case group, groupUncaptured:
		tokens := t.value.([]token)
		sStart, tEnd := tokenToNfa(&tokens[0])
		for _, t := range tokens[1:] {
			startNext, endNext := tokenToNfa(&t)
			tEnd.transitions[epsilonChar] = append(tEnd.transitions[epsilonChar], startNext)
			tEnd = endNext
		}
		// TODO: verify this
		start.transitions[epsilonChar] = []*state{sStart}
	case repeat:
		payload := t.value.(repeatPayload)

		from, to := tokenToNfa(&payload.token)
		start.transitions[epsilonChar] = append(start.transitions[epsilonChar], from)

		if payload.min == 0 {
			start.transitions[epsilonChar] = append(start.transitions[epsilonChar], end)
		}

		// TODO verify
		maxCopy := payload.max
		if payload.max == repeatInfinity {
			maxCopy = 0
		}

		for i := 2; i < maxCopy; i++ {
			s, e := tokenToNfa(&payload.token)

			to.transitions[epsilonChar] = append(to.transitions[epsilonChar], s)

			from = s
			to = e

			if i > payload.min {
				s.transitions[epsilonChar] = append(s.transitions[epsilonChar], end)
			}
		}

		to.transitions[epsilonChar] = append(to.transitions[epsilonChar], end)

		if payload.max == repeatInfinity {
			end.transitions[epsilonChar] = append(end.transitions[epsilonChar], from)
		}
	default:
		panic("unknown type of token")
	}

	return start, end
}
