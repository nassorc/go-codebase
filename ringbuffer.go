package gandalf

import (
	"fmt"
	"reflect"
)

func NewRingBuffer[T any](capacity int) *Ringbuffer[T] {
	return &Ringbuffer[T]{
		read:     0,
		write:    -1,
		capacity: capacity,
		buffer:   make([]T, capacity, capacity),
	}
}

type Ringbuffer[T any] struct {
	read     int
	write    int
	capacity int
	buffer   []T
}

func Enqueue[T any](rb *Ringbuffer[T], value T) {
	rb.write = (rb.write + 1) % rb.capacity
	rb.buffer[rb.write] = value

	if rb.read == rb.write && len(rb.buffer) > 0 {
		rb.read = (rb.read + 1) % rb.capacity
	}
}

func Dequeue[T any](rb *Ringbuffer[T]) (T, error) {
	t := reflect.TypeOf((*T)(nil))
	zero := reflect.Zero(t).Interface()

	if len(rb.buffer) == 0 {
		return zero.(T), fmt.Errorf("Empty Ringbuffer")
	}

	out := rb.buffer[rb.read]
	rb.read = (rb.read + 1) % rb.capacity

	return out, nil
}
