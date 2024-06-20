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

func (buf *Ringbuffer[T]) Enqueue(value T) {
	buf.write = (buf.write + 1) % buf.capacity
	buf.buffer[buf.write] = value

	if buf.read == buf.write && len(buf.buffer) > 0 {
		buf.read = (buf.read + 1) % buf.capacity
	}
}

func (buf *Ringbuffer[T]) Dequeue() (T, error) {
	t := reflect.TypeOf((*T)(nil))
	zero := reflect.Zero(t).Interface()

	if len(buf.buffer) == 0 {
		return zero.(T), fmt.Errorf("empty buffer")
	}

	out := buf.buffer[buf.read]
	buf.read = (buf.read + 1) % buf.capacity

	return out, nil
}
