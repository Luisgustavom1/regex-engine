package main

type state struct {
	start       bool
	terminal    bool
	transitions stateTransition
}

type stateTransition = map[uint8][]*state

const epsilonChar uint8 = 0 // empty character

func toNfa(ctx *parseContext) *state {
	s := &state{
		transitions: map[uint8][]*state{},
	}
	startState, endState := tokenToNfa(&ctx.tokens[0], s)

	for _, t := range ctx.tokens[1:] {
		_, endNext := tokenToNfa(&t, endState)
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

func tokenToNfa(t *token, start *state) (*state, *state) {
	switch t.tokenType {
	case literal:
		ch := t.value.(uint8)
		end := &state{
			transitions: stateTransition{},
		}
		start.transitions[ch] = []*state{end}
		return start, end
	case or:
		values := t.value.([]token)
		left := values[0]
		right := values[1]

		startLeft, endLeft := tokenToNfa(&left, start)
		startRight, endRight := tokenToNfa(&right, start)

		start.transitions[epsilonChar] = []*state{startLeft, startRight}
		end := &state{
			transitions: stateTransition{},
		}
		endLeft.transitions[epsilonChar] = []*state{end}
		endRight.transitions[epsilonChar] = []*state{end}
		return start, end
	case bracket:
		literalsSet := t.value.(map[uint8]bool)
		end := &state{
			transitions: stateTransition{},
		}
		for ch := range literalsSet {
			start.transitions[ch] = []*state{end}
		}
		return start, end
	case group, groupUncaptured:
		tokens := t.value.([]token)
		tStart, tEnd := tokenToNfa(&tokens[0], start)
		for _, t := range tokens[1:] {
			_, endNext := tokenToNfa(&t, tEnd)
			tEnd = endNext
		}
		start.transitions[epsilonChar] = append(start.transitions[epsilonChar], tStart)
		return start, tEnd
	case repeat:
		payload := t.value.(repeatPayload)

		end := &state{
			transitions: stateTransition{},
		}
		if payload.min == 0 {
			start.transitions[epsilonChar] = append(start.transitions[epsilonChar], end)
		}

		maxCopy := payload.max
		if payload.max == repeatInfinity {
			if payload.min == 0 {
				maxCopy = 1
			} else {
				maxCopy = payload.min
			}
		}

		from, to := tokenToNfa(&payload.token, start)
		start.transitions[epsilonChar] = append(start.transitions[epsilonChar], from)

		for i := 2; i < maxCopy; i++ {
			tmpStart := &state{
				transitions: map[uint8][]*state{},
			}
			s, e := tokenToNfa(&payload.token, tmpStart)

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
		return start, end
	default:
		panic("unknown type of token")
	}
}
