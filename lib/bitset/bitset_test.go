package bitset

import "testing"

func TestBitset_NewSignature(t *testing.T) {
	s := NewBitset(4)

	if len(s.data) != 4 {
		t.Errorf("Signature expected size=4, got=%d", len(s.data))
	}
}

func TestBitset_Set(t *testing.T) {
	s := NewBitset(4)

	compareSignature(t, s, "0000")

	s.Set(0)
	compareSignature(t, s, "0001")
	s.Set(3)
	compareSignature(t, s, "1001")
	s.Set(1)
	compareSignature(t, s, "1011")
	s.Set(2)
	compareSignature(t, s, "1111")

	s = NewBitset(8)

	compareSignature(t, s, "00000000")
	s.SetAll()
	compareSignature(t, s, "11111111")
}

func TestBitset_Test(t *testing.T) {
	s := NewBitset(4)
	s.Set(0)
	s.Set(3)

	tt := []bool{true, false, false, true}

	for idx, expected := range tt {
		if actual := s.Test(idx); actual != expected {
			t.Errorf("Test Expected=%t, Got=%t", expected, actual)
		}
	}
}

func TestBitset_Reset(t *testing.T) {
	s := NewBitset(4)
	compareSignature(t, s, "0000")
	s.Set(0)
	s.Set(3)
	compareSignature(t, s, "1001")
	s.Reset(0)
	compareSignature(t, s, "1000")
	s.Reset(3)
	compareSignature(t, s, "0000")

	s.Set(1)
	s.Set(2)
	compareSignature(t, s, "0110")
	s.ResetAll()
	compareSignature(t, s, "0000")
}

func TestBitset_String(t *testing.T) {
	s := NewBitset(4)
	s.Set(0)
	s.Set(3)

	if s.String() != "1001" {
		t.Errorf("String Expected=%s, Got=%s", "1001", s.String())
	}

	s.ResetAll()
	if s.String() != "0000" {
		t.Errorf("String Expected=%s, Got=%s", "0000", s.String())
	}
}

func TestBitset_Int(t *testing.T) {
	s := NewBitset(4)

	s.Set(0)
	if s.Int() != 1 {
		t.Errorf("Int Expected=%d, Got=%d", 1, s.Int())
	}

	s.Set(1)
	if s.Int() != 3 {
		t.Errorf("Int Expected=%d, Got=%d", 3, s.Int())
	}

	s.Set(2)
	if s.Int() != 7 {
		t.Errorf("Int Expected=%d, Got=%d", 7, s.Int())
	}

	s.Set(3)
	if s.Int() != 15 {
		t.Errorf("Int Expected=%d, Got=%d", 15, s.Int())
	}
}

func TestBitset_IsSubset(t *testing.T) {
	s := NewBitset(4)
	tt := NewBitset(4)

	s.Set(1)
	s.Set(3)

	if s.IsSubset(tt) {
		t.Errorf("Expected=%v, Got=%v", false, true)
	}

	tt.Set(2)

	if s.IsSubset(tt) {
		t.Errorf("Expected=%v, Got=%v", false, true)
	}

	tt.Set(1)

	if s.IsSubset(tt) {
		t.Errorf("Expected=%v, Got=%v", false, true)
	}

	tt.Set(3)

	if !s.IsSubset(tt) {
		t.Errorf("Expected=%v, Got=%v", true, false)
	}
}

func compareSignature(t *testing.T, got Bitset, expected string) {
	if got.String() != expected {
		t.Errorf("expected=%s, got=%s", got.String(), expected)
	}
}
