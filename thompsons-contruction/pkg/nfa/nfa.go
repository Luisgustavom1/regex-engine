package nfa

import (
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/ds"
	"github.com/Luisgustavom1/regex-engine/thompsons-construction/pkg/parser"
)

type State struct {
	IsEnd              bool
	Transitions        map[byte]*State
	EpsilonTransitions []*State
}

type Nfa struct {
	Start *State
	End   *State
}

func NewState(isEnd bool) *State {
	return &State{
		IsEnd:              isEnd,
		Transitions:        map[byte]*State{},
		EpsilonTransitions: []*State{},
	}
}

func NewNfa(s *State, e *State) Nfa {
	return Nfa{
		Start: s,
		End:   e,
	}
}

func addEpsilonTransition(from *State, to *State) {
	from.EpsilonTransitions = append(from.EpsilonTransitions, to)
}

func addTransition(from *State, to *State, symbol byte) {
	from.Transitions[symbol] = to
}

func FromEpsilon() Nfa {
	start := NewState(false)
	end := NewState(true)
	addEpsilonTransition(start, end)

	return NewNfa(start, end)
}

func FromSymbol(symbol byte) Nfa {
	start := NewState(false)
	end := NewState(true)
	addTransition(start, end, symbol)

	return NewNfa(start, end)
}

func concat(first *Nfa, second *Nfa) Nfa {
	addEpsilonTransition(first.End, second.Start)
	first.End.IsEnd = false

	return NewNfa(first.Start, second.End)
}

func union(first Nfa, second Nfa) Nfa {
	s := NewState(false)
	e := NewState(true)

	addEpsilonTransition(s, first.Start)
	addEpsilonTransition(s, second.Start)

	addEpsilonTransition(first.End, e)
	addEpsilonTransition(second.End, e)
	first.End.IsEnd = false
	second.End.IsEnd = false

	return NewNfa(s, e)
}

func closure(nfa Nfa) Nfa {
	s := NewState(false)
	e := NewState(true)

	addEpsilonTransition(nfa.End, e)
	addEpsilonTransition(nfa.End, nfa.Start)
	nfa.End.IsEnd = false

	addEpsilonTransition(s, e)
	addEpsilonTransition(s, nfa.Start)

	return NewNfa(s, e)
}

func ToNfa(postfixExp string) Nfa {
	if postfixExp == "" {
		return FromEpsilon()
	}

	stack := ds.NewStack[Nfa]()

	for _, c := range postfixExp {
		if parser.Operators(c) == parser.UNION {
			s := stack.Pop()
			e := stack.Pop()
			stack.Push(union(s, e))
		} else if parser.Operators(c) == parser.CLOSURE {
			nfa := stack.Pop()
			stack.Push(closure(nfa))
		} else if parser.Operators(c) == parser.CONCAT {
			s := stack.Pop()
			e := stack.Pop()
			stack.Push(union(s, e))
		} else {
			stack.Push(FromSymbol(byte(c)))
		}
	}

	return stack.Pop()
}
