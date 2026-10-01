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
	var q Slide[T]
	q.init(capacity)
	return &q
}

func (q *Slide[T]) init(capacity int) {
	if capacity <= 0 {
		capacity = _defaultCapacity
	}
	q.buff = make([]T, capacity+1)
	q.head = 0
	q.tail = 0
}

// Empty reports whether q holds no items.
func (q *Slide[T]) Empty() bool {
	return q.head == q.tail
}

// Len returns the number of items in q.
func (q *Slide[T]) Len() int {
	if q.head <= q.tail {
		return q.tail - q.head
	}
	return len(q.buff) - q.head + q.tail
}

// Cap returns the maximum number of items q can hold.
func (q *Slide[T]) Cap() int {
	if len(q.buff) == 0 {
		return _defaultCapacity
	}
	return len(q.buff) - 1
}

// tailInfo returns the internal buffer index of the newest item.
// q must not be empty.
func (q *Slide[T]) tailInfo() int {
	return (len(q.buff) + q.tail - 1) % len(q.buff)
}

// Clear removes all items from q, keeping its capacity.
func (q *Slide[T]) Clear() {
	clear(q.buff)
	q.head = 0
	q.tail = 0
}

// Push appends x at the tail. If q is full, the oldest item is discarded.
func (q *Slide[T]) Push(x T) {
	if len(q.buff) == 0 {
		q.init(_defaultCapacity)
	}

	if (q.tail+1)%len(q.buff) == q.head {
		q.head = (q.head + 1) % len(q.buff)
	}

	q.buff[q.tail] = x
	q.tail = (q.tail + 1) % len(q.buff)
}

// TryPopHead removes and returns the oldest item.
// It reports false if q is empty.
func (q *Slide[T]) TryPopHead() (x T, ok bool) {
	if q.Empty() {
		return x, false
	}

	x = q.buff[q.head]
	var zero T
	q.buff[q.head] = zero
	q.head = (q.head + 1) % len(q.buff)
	return x, true
}

// TryPopTail removes and returns the newest item.
// It reports false if q is empty.
func (q *Slide[T]) TryPopTail() (x T, ok bool) {
	if q.Empty() {
		return x, false
	}
	exTail := q.tailInfo()
	x = q.buff[exTail]
	var zero T
	q.buff[exTail] = zero
	q.tail = exTail
	return x, true
}

// TryPeekTail returns the newest item without removing it.
// It reports false if q is empty.
func (q *Slide[T]) TryPeekTail() (x T, ok bool) {
	if q.Empty() {
		return x, false
	}
	return q.buff[q.tailInfo()], true
}

// TryPeekHead returns the oldest item without removing it.
// It reports false if q is empty.
func (q *Slide[T]) TryPeekHead() (x T, ok bool) {
	if q.Empty() {
		return x, false
	}
	return q.buff[q.head], true
}

// SnapshotQueue appends the items of q to dst, oldest first, and returns the
// extended slice. Pass a dst with spare capacity to avoid allocation, or nil.
func (q *Slide[T]) SnapshotQueue(dst []T) []T {
	if q.head <= q.tail {
		return append(dst, q.buff[q.head:q.tail]...)
	}

	dst = append(dst, q.buff[q.head:]...)
	return append(dst, q.buff[:q.tail]...)
}

// SnapshotStack appends the items of q to dst, newest first, and returns the
// extended slice. Pass a dst with spare capacity to avoid allocation, or nil.
func (q *Slide[T]) SnapshotStack(dst []T) []T {
	n := len(dst)
	dst = q.SnapshotQueue(dst)
	slices.Reverse(dst[n:])
	return dst
}
