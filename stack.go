package slide

import "sync"

// Stack is a fixed-capacity LIFO (Last In, First Out) stack backed by a Slide.
// When full, pushing discards the oldest item.
//
// The zero value is an empty Stack with the default capacity of 16, ready to use.
// Stack is safe for concurrent use.
type Stack[T any] struct {
	mu    sync.Mutex
	slide Slide[T]
}

// NewStack returns an empty Stack that holds up to capacity items.
// A capacity <= 0 selects the default of 16.
func NewStack[T any](capacity int) *Stack[T] {
	var s Stack[T]
	s.slide.init(capacity)
	return &s
}

// Push adds item to the stack (tail). If the stack is full, the oldest item is discarded.
func (s *Stack[T]) Push(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slide.Push(item)
}

// Pop removes and returns the top (tail) item. It reports false if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.slide.PopTail()
}

// Peek returns the top (tail) item without removing it. It reports false if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.slide.PeekTail()
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.slide.Len()
}

// Cap returns the maximum number of items the stack can hold.
func (s *Stack[T]) Cap() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.slide.Cap()
}
