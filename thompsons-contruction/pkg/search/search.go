package search

import "github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/nfa"

func Search(n nfa.Nfa, word string) (match bool) {
	currentStates := &[]nfa.State{}
	addNextState(n.Start, currentStates, &[]*nfa.State{})

	for _, c := range word {
		nextStates := &[]nfa.State{}

		for _, state := range *currentStates {
			nextState := state.Transitions[byte(c)]
			if nextState != nil {
				addNextState(nextState, nextStates, &[]*nfa.State{})
			}
		}

		currentStates = nextStates
	}

	for _, state := range *currentStates {
		if state.IsEnd {
			match = true
			return match
		}
	}

	return match
}

func addNextState(state *nfa.State, nextState *[]nfa.State, visited *[]*nfa.State) {
	if len(state.EpsilonTransitions) > 0 {
		for _, s := range state.EpsilonTransitions {
			alreadyVisited := false

			for _, v := range *visited {
				if s == v {
					alreadyVisited = true
					break
				}
			}

			if !alreadyVisited {
				(*visited) = append(*visited, s)
				addNextState(s, nextState, visited)
			}
		}
	} else {
		(*nextState) = append(*nextState, *state)
	}
}
