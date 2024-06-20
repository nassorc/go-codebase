package gandalf

func NewRingBuffer[T any](capacity int) *Ringbuffer[T] {
	return &Ringbuffer[T]{
		read:     0,
		write:    -1,
		capacity: capacity,
		buffer:   make([]T, capacity),
	}
}

type Ringbuffer[T any] struct {
	read     int
	write    int
	size     int
	capacity int
	buffer   []T
}

func (buf *Ringbuffer[T]) Enqueue(value T) bool {
	if buf.Full() {
		return false
	}
	buf.write = (buf.write + 1) % buf.capacity
	buf.buffer[buf.write] = value
	buf.size += 1

	return true
}

func (buf *Ringbuffer[T]) Deque() (T, bool) {
	if buf.Empty() {
		var zeroValue T
		return zeroValue, false
	}

	out := buf.buffer[buf.read]
	buf.read = (buf.read + 1) % buf.capacity
	buf.size -= 1

	return out, true
}

func (buf *Ringbuffer[T]) Size() int {
	return buf.size
}

func (buf *Ringbuffer[T]) Full() bool {
	return buf.size == buf.capacity
}

func (buf *Ringbuffer[T]) Empty() bool {
	// initialized to -1
	if buf.write == -1 {
		return true
	}
	if buf.Size() == 0 {
		return true
	}
	return false
}
