package gandalf

import (
	"testing"
)

func Test_RingBuffer(t *testing.T) {
	rb := NewRingBuffer[int](3)

	Enqueue(rb, 10)
	Enqueue(rb, 20)
	Enqueue(rb, 30)

	out, err := Dequeue(rb)

	if err != nil {
		t.Error(err)
	}

	if out != 10 {
		t.Fatalf("expected=%d, got=%d", 10, out)
	}

	out, err = Dequeue(rb)
	if out != 20 {
		t.Fatalf("expected=%d, got=%d", 20, out)
	}

	Enqueue(rb, 100)

	out, err = Dequeue(rb)
	if out != 30 {
		t.Fatalf("expected=%d, got=%d", 30, out)
	}

	out, err = Dequeue(rb)
	if out != 100 {
		t.Fatalf("expected=%d, got=%d", 100, out)
	}
}
