package slide

import "sync"

// Queue is a fixed-capacity FIFO (First In, First Out) queue backed by a Slide.
// When full, pushing discards the oldest item.
//
// The zero value is an empty Queue with the default capacity of 16, ready to use.
// Queue is safe for concurrent use.
type Queue[T any] struct {
	mu    sync.Mutex
	slide Slide[T]
}

// NewQueue returns an empty Queue that holds up to capacity items.
// A capacity <= 0 selects the default of 16.
func NewQueue[T any](capacity int) *Queue[T] {
	var q Queue[T]
	q.slide.init(capacity)
	return &q
}

// Push adds item to the queue (tail). If the queue is full, the oldest item is discarded.
func (q *Queue[T]) Push(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.slide.Push(item)
}

// Pop removes and returns the oldest (head) item. It reports false if the queue is empty.
func (q *Queue[T]) Pop() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.slide.PopHead()
}

// Peek returns the oldest (head) item without removing it. It reports false if the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.slide.PeekHead()
}

// Len returns the number of items in the queue.
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.slide.Len()
}

// Cap returns the maximum number of items the queue can hold.
func (q *Queue[T]) Cap() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.slide.Cap()
}
