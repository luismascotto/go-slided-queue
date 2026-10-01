# go-slided-queue

[![Go Reference](https://pkg.go.dev/badge/github.com/luismascotto/go-slided-queue.svg)](https://pkg.go.dev/github.com/luismascotto/go-slided-queue)
[![CI](https://github.com/luismascotto/go-slided-queue/actions/workflows/ci.yml/badge.svg)](https://github.com/luismascotto/go-slided-queue/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/luismascotto/go-slided-queue)](https://goreportcard.com/report/github.com/luismascotto/go-slided-queue)

Generic, fixed-capacity containers for Go that discard the oldest item when full:

- `Slide`: a ring buffer you can pop and peek from either end.
- `Queue`: a FIFO (First In, First Out) queue backed by a `Slide`.
- `Stack`: a LIFO (Last In, First Out) stack backed by a `Slide`.

## Install

```sh
go get github.com/luismascotto/go-slided-queue
```

Requires Go 1.21+.

## Usage

```go
package main

import (
	"fmt"

	slide "github.com/luismascotto/go-slided-queue"
)

func main() {
	// Slide: ring buffer, both ends.
	s := slide.New[int](3)
	for i := 1; i <= 5; i++ {
		s.Push(i)
	}
	fmt.Println(s.SnapshotQueue(nil)) // [3 4 5]

	head, _ := s.PopHead() // 3
	tail, _ := s.PopTail() // 5
	fmt.Println(head, tail)

	// Queue: oldest first.
	q := slide.NewQueue[string](2)
	q.Push("a")
	q.Push("b")
	q.Push("c") // discards "a"
	first, _ := q.Pop() // "b"
	fmt.Println(first)

	// Stack: newest first.
	st := slide.NewStack[string](2)
	st.Push("a")
	st.Push("b")
	st.Push("c") // discards "a"
	top, _ := st.Pop() // "c"
	fmt.Println(top)
}
```

The zero value of each type is ready to use, with a default capacity of 16.

## Concurrency

- `Slide` is not safe for concurrent use; guard it with a mutex if shared across goroutines.
- `Queue` and `Stack` lock internally and are safe for concurrent use. They must not be copied after first use.

## License

[MIT](LICENSE)
