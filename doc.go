// Package slide provides generic fixed-capacity containers that discard the
// oldest item when full.
//
//   - Slide is a ring buffer that can be pushed at the tail and popped or
//     peeked from either end. It is not safe for concurrent use; callers must
//     synchronize access.
//   - Queue is a FIFO (First In, First Out) queue backed by a Slide.
//   - Stack is a LIFO (Last In, First Out) stack backed by a Slide.
//
// Queue and Stack are safe for concurrent use and must not be copied after
// first use.
package slide
