package slide

import (
	"sync"
	"testing"
)

func TestQueueZeroValue(t *testing.T) {
	var q Queue[int]

	if _, ok := q.Pop(); ok {
		t.Fatal("Pop on empty queue returned ok")
	}
	if _, ok := q.Peek(); ok {
		t.Fatal("Peek on empty queue returned ok")
	}
	if got := q.Cap(); got != _defaultCapacity {
		t.Fatalf("Cap() = %d, want %d", got, _defaultCapacity)
	}

	q.Push(1)
	if got, ok := q.Pop(); !ok || got != 1 {
		t.Fatalf("Pop() = %d, %v, want 1, true", got, ok)
	}
}

func TestQueueFIFO(t *testing.T) {
	q := NewQueue[int](3)
	q.Push(1)
	q.Push(2)
	q.Push(3)

	if got, ok := q.Peek(); !ok || got != 1 {
		t.Fatalf("Peek() = %d, %v, want 1, true", got, ok)
	}
	if got := q.Len(); got != 3 {
		t.Fatalf("Len() after Peek = %d, want 3", got)
	}

	for _, want := range []int{1, 2, 3} {
		if got, ok := q.Pop(); !ok || got != want {
			t.Fatalf("Pop() = %d, %v, want %d, true", got, ok, want)
		}
	}
	if _, ok := q.Pop(); ok {
		t.Fatal("Pop on drained queue returned ok")
	}
}

func TestQueueDiscardsOldestWhenFull(t *testing.T) {
	q := NewQueue[int](2)
	q.Push(1)
	q.Push(2)
	q.Push(3)

	if got := q.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	for _, want := range []int{2, 3} {
		if got, ok := q.Pop(); !ok || got != want {
			t.Fatalf("Pop() = %d, %v, want %d, true", got, ok, want)
		}
	}
}

func TestQueueConcurrentPush(t *testing.T) {
	const goroutines, perGoroutine = 8, 100
	q := NewQueue[int](goroutines * perGoroutine)

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				q.Push(i)
			}
		}()
	}
	wg.Wait()

	if got, want := q.Len(), goroutines*perGoroutine; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}
}
