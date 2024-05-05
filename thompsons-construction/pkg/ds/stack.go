package ds

type Stack[T any] struct {
	v []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		v: []T{},
	}
}

func (s *Stack[T]) Push(v T) {
	s.v = append(s.v, v)
}

func (s *Stack[T]) Pop() T {
	l := len(s.v) - 1
	last := s.v[l]
	s.v = s.v[:l]
	return last
}

func (s *Stack[T]) Values() []T {
	return s.v
}

func (s *Stack[T]) Len() int {
	return len(s.v)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.v) == 0
}

func (s *Stack[T]) Peek() T {
	return s.v[len(s.v)-1]
}
