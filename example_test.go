package slide_test

import (
	"fmt"

	slide "github.com/luismascotto/go-slided-queue"
)

func ExampleSlide() {
	s := slide.New[int](3)
	for i := 1; i <= 5; i++ {
		s.Push(i)
	}
	fmt.Println(s.Len(), s.SnapshotQueue(nil))

	head, _ := s.PopHead()
	tail, _ := s.PopTail()
	fmt.Println(head, tail, s.SnapshotQueue(nil))
	// Output:
	// 3 [3 4 5]
	// 3 5 [4]
}

func ExampleQueue() {
	q := slide.NewQueue[string](2)
	q.Push("a")
	q.Push("b")
	q.Push("c") // discards "a"

	for {
		item, ok := q.Pop()
		if !ok {
			break
		}
		fmt.Println(item)
	}
	// Output:
	// b
	// c
}

func ExampleStack() {
	s := slide.NewStack[string](2)
	s.Push("a")
	s.Push("b")
	s.Push("c") // discards "a"

	for {
		item, ok := s.Pop()
		if !ok {
			break
		}
		fmt.Println(item)
	}
	// Output:
	// c
	// b
}
