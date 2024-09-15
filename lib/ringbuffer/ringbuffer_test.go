package ringbuffer

import (
	"testing"
)

func TestRingBuffer(t *testing.T) {
	var buf = NewRingBuffer[int](3)

	if got := buf.Empty(); got != true {
		t.Errorf("Expected=%v, Got=%v", true, got)
	}

	tt := []int{3, 5, 8}

	for _, t := range tt {
		buf.Enqueue(t)
	}

	if got := buf.Size(); got != 3 {
		t.Errorf("Expected=%v, Got=%v", 3, got)
	}

	size := buf.Size()

	for idx := 0; idx < size; idx++ {
		got, ok := buf.Deque()

		if !ok {
			t.Errorf("Expected=%v, Got=%v", true, ok)
		}
		if got != tt[idx] {
			t.Errorf("Expected=%v, Got=%v", tt[idx], got)
		}
	}

	if got := buf.Empty(); got != true {
		t.Errorf("Expected=%v, Got=%v", true, got)
	}
}

func TestRingBufferEnqueueDeque(t *testing.T) {
	var size = 2
	var buf = NewRingBuffer[int](size)

	for idx := 0; idx < size; idx++ {
		buf.Enqueue(idx)
	}

	// deque all elements, then requeue first element
	out1, _ := buf.Deque() // 0
	buf.Deque()            // 1
	buf.Enqueue(out1)      // 0

	out3, _ := buf.Deque() // 0
	if out3 != 0 {
		t.Errorf("Expected=%v Got=%v", 1, out3)
	}

	buf = NewRingBuffer[int](5)

	tt := []int{3, 5, 8, 13, 21}

	for _, t := range tt {
		buf.Enqueue(t)
	}

	if got := buf.Size(); got != 5 {
		t.Errorf("Expected=%v, Got=%v", 3, got)
	}

	if got, _ := buf.Deque(); got != 3 {
		t.Errorf("Expected=%v, Got=%v", 3, got)
	}

	if got, _ := buf.Deque(); got != 5 {
		t.Errorf("Expected=%v, Got=%v", 5, got)
	}

	buf.Enqueue(100)

	for _, expected := range []int{8, 13, 21, 100} {
		got, ok := buf.Deque()

		if !ok {
			t.Errorf("Expected=%v, Got=%v", true, ok)
		}

		if got != expected {
			t.Errorf("Expected=%v, Got=%v", expected, got)
		}
	}

	if got := buf.Empty(); got != true {
		t.Errorf("Expected=%v, Got=%v", true, got)
	}
}

func TestRingBufferIsEmpty(t *testing.T) {
	buf := NewRingBuffer[int](2)

	if !buf.Empty() {
		t.Errorf("expected ringbuffer to be empty.")
	}

	buf.Enqueue(3)
	if buf.Empty() {
		t.Errorf("expected ringbuffer to NOT be empty.")
	}

	buf.Enqueue(5)
	if buf.Empty() {
		t.Errorf("expected ringbuffer to NOT be empty.")
	}

	// test if enqueuing full buffer breaks code
	buf.Enqueue(8)
	if buf.Empty() {
		t.Errorf("expected ringbuffer to NOT be empty.")
	}

	buf.Deque()
	buf.Deque()
	// should be empty
	if !buf.Empty() {
		t.Errorf("expected ringbuffer to be empty.")
	}
	buf.Deque()
	if !buf.Empty() {
		t.Errorf("expected ringbuffer to be empty.")
	}
}

func TestRingBufferDequeReturnsFalseIfEmpty(t *testing.T) {
	buf := NewRingBuffer[int](2)

	if got := buf.Empty(); got != true {
		t.Errorf("Expected=%v, Got=%v", true, got)
	}

	if _, ok := buf.Deque(); ok != false {
		t.Errorf("Expected=%v, Got=%v", false, ok)
	}

	buf.Enqueue(2)
	if _, ok := buf.Deque(); ok != true {
		t.Errorf("Expected=%v, Got=%v", false, ok)
	}
	if _, ok := buf.Deque(); ok != false {
		t.Errorf("Expected=%v, Got=%v", false, ok)
	}
}

func BenchmarkRingBufferEnqueue(b *testing.B) {
	size := b.N
	buf := NewRingBuffer[int](size)

	for idx := 0; idx < size; idx++ {
		buf.Enqueue(idx)
	}
}

func BenchmarkRingBufferDequeue(b *testing.B) {
	size := b.N
	buf := NewRingBuffer[int](size)

	for idx := 0; idx < size; idx++ {
		buf.Enqueue(idx)
	}

	b.ResetTimer()
	for idx := 0; idx < size; idx++ {
		buf.Deque()
	}
}
