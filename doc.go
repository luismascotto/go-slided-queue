// Package slide provides Slide, a generic fixed-capacity ring buffer that
// overwrites the oldest item when full.
//
// Items can be pushed at the tail and popped or peeked from either end,
// so a Slide can be consumed as a FIFO queue or as a LIFO stack.
//
// Slide is not safe for concurrent use; callers must synchronize access.
package slide
