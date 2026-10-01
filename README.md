# go-slided-queue

[![Go Reference](https://pkg.go.dev/badge/github.com/luismascotto/go-slided-queue.svg)](https://pkg.go.dev/github.com/luismascotto/go-slided-queue)
[![CI](https://github.com/luismascotto/go-slided-queue/actions/workflows/ci.yml/badge.svg)](https://github.com/luismascotto/go-slided-queue/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/luismascotto/go-slided-queue)](https://goreportcard.com/report/github.com/luismascotto/go-slided-queue)

A generic, fixed-capacity ring buffer for Go that overwrites the oldest item when full.
Pop and peek from either end to use it as a sliding FIFO queue or LIFO stack.

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
	q := slide.New[int](3)
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	fmt.Println(q.SnapshotQueue(nil)) // [3 4 5]

	head, _ := q.TryPopHead() // 3
	tail, _ := q.TryPopTail() // 5
	fmt.Println(head, tail)
}
```

The zero value `slide.Slide[T]{}` is ready to use.

`Slide` is not safe for concurrent use; guard it with a mutex if shared across goroutines.

## License

[MIT](LICENSE)
