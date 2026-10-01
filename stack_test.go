package slide

import (
	"sync"
	"testing"
)

func TestStackZeroValue(t *testing.T) {
	var s Stack[int]

	if _, ok := s.Pop(); ok {
		t.Fatal("Pop on empty stack returned ok")
	}
	if _, ok := s.Peek(); ok {
		t.Fatal("Peek on empty stack returned ok")
	}
	if got := s.Cap(); got != _defaultCapacity {
		t.Fatalf("Cap() = %d, want %d", got, _defaultCapacity)
	}

	s.Push(1)
	if got, ok := s.Pop(); !ok || got != 1 {
		t.Fatalf("Pop() = %d, %v, want 1, true", got, ok)
	}
}

func TestStackLIFO(t *testing.T) {
	s := NewStack[int](3)
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if got, ok := s.Peek(); !ok || got != 3 {
		t.Fatalf("Peek() = %d, %v, want 3, true", got, ok)
	}
	if got := s.Len(); got != 3 {
		t.Fatalf("Len() after Peek = %d, want 3", got)
	}

	for _, want := range []int{3, 2, 1} {
		if got, ok := s.Pop(); !ok || got != want {
			t.Fatalf("Pop() = %d, %v, want %d, true", got, ok, want)
		}
	}
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop on drained stack returned ok")
	}
}

func TestStackDiscardsOldestWhenFull(t *testing.T) {
	s := NewStack[int](2)
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if got := s.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	for _, want := range []int{3, 2} {
		if got, ok := s.Pop(); !ok || got != want {
			t.Fatalf("Pop() = %d, %v, want %d, true", got, ok, want)
		}
	}
}

func TestStackConcurrentPush(t *testing.T) {
	const goroutines, perGoroutine = 8, 100
	s := NewStack[int](goroutines * perGoroutine)

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				s.Push(i)
			}
		}()
	}
	wg.Wait()

	if got, want := s.Len(), goroutines*perGoroutine; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}
}
