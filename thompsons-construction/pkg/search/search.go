package search

import "github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/nfa"

func Search(n nfa.Nfa, word string) (match bool) {
	currentStates := []nfa.State{}
	addNextState(n.Start, &currentStates, &[]*nfa.State{})

	for _, symbol := range word {
		nextStates := []nfa.State{}

		for _, state := range currentStates {
			nextState := state.Transitions[byte(symbol)]
			if nextState != nil {
				addNextState(nextState, &nextStates, &[]*nfa.State{})
			}
		}

		currentStates = nextStates
	}

	for _, state := range currentStates {
		if state.IsEnd {
			match = true
			return match
		}
	}

	return match
}

func addNextState(state *nfa.State, nextStates *[]nfa.State, visited *[]*nfa.State) {
	if len(state.EpsilonTransitions) > 0 {
		for _, st := range state.EpsilonTransitions {
			alreadyVisited := false

			for _, v := range *visited {
				if st == v {
					alreadyVisited = true
					break
				}
			}

			if !alreadyVisited {
				(*visited) = append(*visited, st)
				addNextState(st, nextStates, visited)
			}
		}
	} else {
		(*nextStates) = append(*nextStates, *state)
	}
}
