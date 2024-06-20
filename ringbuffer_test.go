package gandalf

import (
	"testing"
)

func Test_RingBuffer(t *testing.T) {
	rb := NewRingBuffer[int](3)

	rb.Enqueue(10)
	rb.Enqueue(20)
	rb.Enqueue(30)

	out, err := rb.Dequeue()

	if err != nil {
		t.Error(err)
	}

	if out != 10 {
		t.Fatalf("expected=%d, got=%d", 10, out)
	}

	out, _ = rb.Dequeue()
	if out != 20 {
		t.Fatalf("expected=%d, got=%d", 20, out)
	}

	rb.Enqueue(100)

	out, _ = rb.Dequeue()
	if out != 30 {
		t.Fatalf("expected=%d, got=%d", 30, out)
	}

	out, _ = rb.Dequeue()
	if out != 100 {
		t.Fatalf("expected=%d, got=%d", 100, out)
	}
}
