package slide

import "slices"

const _defaultCapacity = 16

// Slide is a fixed-capacity ring buffer that overwrites its oldest item when full.
// Items are pushed at the tail and can be removed from either end.
//
// The zero value is an empty Slide with the default capacity of 16, ready to use.
type Slide[T any] struct {
	// buff holds one more slot than the capacity so that head == tail
	// unambiguously means empty.
	buff []T
	head int // index of the oldest item
	tail int // index of the next free slot
}

// New returns an empty Slide that holds up to capacity items.
// A capacity <= 0 selects the default of 16.
func New[T any](capacity int) *Slide[T] {
	var s Slide[T]
	s.init(capacity)
	return &s
}

func (s *Slide[T]) init(capacity int) {
	if capacity <= 0 {
		capacity = _defaultCapacity
	}
	s.buff = make([]T, capacity+1)
	s.head = 0
	s.tail = 0
}

// Empty reports whether the buffer holds no items.
func (s *Slide[T]) Empty() bool {
	return s.head == s.tail
}

// Len returns the number of items in the buffer.
func (s *Slide[T]) Len() int {
	if s.head <= s.tail {
		return s.tail - s.head
	}
	return len(s.buff) - s.head + s.tail
}

// Cap returns the maximum number of items the buffer can hold.
func (s *Slide[T]) Cap() int {
	if len(s.buff) == 0 {
		return _defaultCapacity
	}
	return len(s.buff) - 1
}

// tailInfo returns the internal buffer index of the newest item.
// The buffer must not be empty.
func (s *Slide[T]) tailInfo() int {
	return (len(s.buff) + s.tail - 1) % len(s.buff)
}

// Clear removes all items from the buffer, keeping its capacity.
func (s *Slide[T]) Clear() {
	clear(s.buff)
	s.head = 0
	s.tail = 0
}

// Push adds item to the buffer (tail). If the buffer is full, the oldest item is discarded.
func (s *Slide[T]) Push(item T) {
	if len(s.buff) == 0 {
		s.init(_defaultCapacity)
	}

	if (s.tail+1)%len(s.buff) == s.head {
		_, _ = s.PopHead()
	}

	s.buff[s.tail] = item
	s.tail = (s.tail + 1) % len(s.buff)
}

// PopHead removes and returns the oldest (head) item. It reports false if the buffer is empty.
func (s *Slide[T]) PopHead() (x T, ok bool) {
	if s.Empty() {
		return x, false
	}

	x = s.buff[s.head]
	var zero T
	s.buff[s.head] = zero
	s.head = (s.head + 1) % len(s.buff)
	return x, true
}

// PopTail removes and returns the newest (tail) item. It reports false if the buffer is empty.
func (s *Slide[T]) PopTail() (x T, ok bool) {
	if s.Empty() {
		return x, false
	}
	exTail := s.tailInfo()
	x = s.buff[exTail]
	var zero T
	s.buff[exTail] = zero
	s.tail = exTail
	return x, true
}

// PeekTail returns the newest (tail) item without removing it. It reports false if the buffer is empty.
func (s *Slide[T]) PeekTail() (x T, ok bool) {
	if s.Empty() {
		return x, false
	}
	return s.buff[s.tailInfo()], true
}

// PeekHead returns the oldest (head) item without removing it. It reports false if the buffer is empty.
func (s *Slide[T]) PeekHead() (x T, ok bool) {
	if s.Empty() {
		return x, false
	}
	return s.buff[s.head], true
}

// SnapshotQueue appends the items of the buffer to dst, oldest first, and returns the
// extended slice. Pass a dst with spare capacity to avoid allocation, or nil.
func (s *Slide[T]) SnapshotQueue(dst []T) []T {
	if s.Empty() {
		return dst
	}
	if s.head < s.tail {
		return append(dst, s.buff[s.head:s.tail]...)
	}

	dst = append(dst, s.buff[s.head:]...)
	return append(dst, s.buff[:s.tail]...)
}

// SnapshotStack appends the items of the buffer to dst, newest first, and returns the
// extended slice. Pass a dst with spare capacity to avoid allocation, or nil.
func (s *Slide[T]) SnapshotStack(dst []T) []T {
	n := len(dst)
	dst = s.SnapshotQueue(dst)
	slices.Reverse(dst[n:])
	return dst
}
