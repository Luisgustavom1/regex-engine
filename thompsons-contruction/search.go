package main

import "github.com/Luisgustavom1/regex-engine/thompsons-construction/nfa"

func search(n nfa.Nfa, word string) {
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
	}

	for _, state := range *currentStates {
		if state.IsEnd {
			println("Match")
			return
		}
	}
}

func addNextState(state *nfa.State, nextState *[]nfa.State, visited *[]*nfa.State) {
	if len(state.EpsilonTransitions) > 0 {
		for _, s := range state.EpsilonTransitions {
			notVisited := true

			for _, v := range *visited {
				if s == v {
					notVisited = false
					break
				}
			}

			if notVisited {
				(*visited) = append(*visited, s)
				addNextState(s, nextState, visited)
			}
		}
	} else {
		(*nextState) = append(*nextState, *state)
	}
}
