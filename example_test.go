package slide_test

import (
	"fmt"

	slide "github.com/luismascotto/go-slided-queue"
)

func ExampleSlide() {
	q := slide.New[int](3)
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	fmt.Println(q.Len(), q.SnapshotQueue(nil))

	head, _ := q.TryPopHead()
	tail, _ := q.TryPopTail()
	fmt.Println(head, tail, q.SnapshotQueue(nil))
	// Output:
	// 3 [3 4 5]
	// 3 5 [4]
}
