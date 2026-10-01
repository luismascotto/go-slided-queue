package slide

import (
	"slices"
	"testing"
)

func TestZeroValue(t *testing.T) {
	var q Slide[int]

	if !q.Empty() {
		t.Fatal("zero value should be empty")
	}
	if _, ok := q.TryPopHead(); ok {
		t.Fatal("TryPopHead on empty buffer returned ok")
	}
	if _, ok := q.TryPopTail(); ok {
		t.Fatal("TryPopTail on empty buffer returned ok")
	}

	q.Push(1)
	q.Push(2)

	if got := q.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	if got, ok := q.TryPeekHead(); !ok || got != 1 {
		t.Fatalf("TryPeekHead() = %d, %v, want 1, true", got, ok)
	}
	if got, ok := q.TryPeekTail(); !ok || got != 2 {
		t.Fatalf("TryPeekTail() = %d, %v, want 2, true", got, ok)
	}
}

func TestPushOverwritesOldest(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		pushes   int
		want     []int
	}{
		{name: "below capacity", capacity: 3, pushes: 2, want: []int{1, 2}},
		{name: "exactly capacity", capacity: 3, pushes: 3, want: []int{1, 2, 3}},
		{name: "wraps once", capacity: 3, pushes: 4, want: []int{2, 3, 4}},
		{name: "wraps many times", capacity: 3, pushes: 10, want: []int{8, 9, 10}},
		{name: "capacity one", capacity: 1, pushes: 3, want: []int{3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New[int](tt.capacity)
			for i := 1; i <= tt.pushes; i++ {
				q.Push(i)
			}

			if got := q.SnapshotQueue(nil); !slices.Equal(got, tt.want) {
				t.Fatalf("SnapshotQueue() = %v, want %v", got, tt.want)
			}
			if got := q.Len(); got != len(tt.want) {
				t.Fatalf("Len() = %d, want %d", got, len(tt.want))
			}
		})
	}
}

func TestPopBothEnds(t *testing.T) {
	q := New[int](4)
	for i := 1; i <= 4; i++ {
		q.Push(i)
	}

	if got, ok := q.TryPopHead(); !ok || got != 1 {
		t.Fatalf("TryPopHead() = %d, %v, want 1, true", got, ok)
	}
	if got, ok := q.TryPopTail(); !ok || got != 4 {
		t.Fatalf("TryPopTail() = %d, %v, want 4, true", got, ok)
	}
	if got := q.SnapshotQueue(nil); !slices.Equal(got, []int{2, 3}) {
		t.Fatalf("SnapshotQueue() = %v, want [2 3]", got)
	}
}

func TestClear(t *testing.T) {
	q := New[int](2)
	q.Push(1)
	q.Push(2)
	q.Clear()

	if !q.Empty() || q.Len() != 0 {
		t.Fatalf("after Clear: Empty() = %v, Len() = %d", q.Empty(), q.Len())
	}
}

func TestNonPositiveCapacityUsesDefault(t *testing.T) {
	for _, capacity := range []int{0, -1, -5} {
		q := New[int](capacity)
		for i := 1; i <= _defaultCapacity+4; i++ {
			q.Push(i)
		}
		if got := q.Len(); got != _defaultCapacity {
			t.Fatalf("New(%d): Len() = %d, want %d", capacity, got, _defaultCapacity)
		}
	}
}

func TestCap(t *testing.T) {
	var zero Slide[int]
	if got := zero.Cap(); got != _defaultCapacity {
		t.Fatalf("zero value Cap() = %d, want %d", got, _defaultCapacity)
	}
	zero.Push(1)
	if got := zero.Cap(); got != _defaultCapacity {
		t.Fatalf("zero value after Push Cap() = %d, want %d", got, _defaultCapacity)
	}

	tests := []struct {
		capacity int
		want     int
	}{
		{capacity: 1, want: 1},
		{capacity: 5, want: 5},
		{capacity: 0, want: _defaultCapacity},
		{capacity: -3, want: _defaultCapacity},
	}
	for _, tt := range tests {
		if got := New[int](tt.capacity).Cap(); got != tt.want {
			t.Fatalf("New(%d).Cap() = %d, want %d", tt.capacity, got, tt.want)
		}
	}
}

func TestSnapshotStack(t *testing.T) {
	tests := []struct {
		name   string
		pushes int
		dst    []int
		want   []int
	}{
		{name: "empty", pushes: 0, dst: nil, want: nil},
		{name: "no wrap", pushes: 2, dst: nil, want: []int{2, 1}},
		{name: "wrapped", pushes: 5, dst: nil, want: []int{5, 4, 3}},
		{name: "keeps dst prefix", pushes: 5, dst: []int{0}, want: []int{0, 5, 4, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New[int](3)
			for i := 1; i <= tt.pushes; i++ {
				q.Push(i)
			}

			if got := q.SnapshotStack(tt.dst); !slices.Equal(got, tt.want) {
				t.Fatalf("SnapshotStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemovalReleasesReferences(t *testing.T) {
	newFull := func() *Slide[*int] {
		q := New[*int](3)
		for i := 0; i < 3; i++ {
			q.Push(new(int))
		}
		return q
	}
	assertAllNil := func(t *testing.T, q *Slide[*int]) {
		t.Helper()
		for i, p := range q.buff {
			if p != nil {
				t.Fatalf("buff[%d] still holds a reference", i)
			}
		}
	}

	t.Run("pop", func(t *testing.T) {
		q := newFull()
		q.TryPopHead()
		q.TryPopTail()
		q.TryPopHead()
		assertAllNil(t, q)
	})

	t.Run("clear", func(t *testing.T) {
		q := newFull()
		q.Clear()
		assertAllNil(t, q)
	})
}

func TestSnapshotQueueAppendsToDst(t *testing.T) {
	q := New[int](2)
	q.Push(1)
	q.Push(2)

	got := q.SnapshotQueue([]int{0})
	if want := []int{0, 1, 2}; !slices.Equal(got, want) {
		t.Fatalf("SnapshotQueue() = %v, want %v", got, want)
	}
}
