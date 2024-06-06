package gandalf

import (
	"fmt"
	"math"
	"strings"
)

func NewSignature(size int) *Signature {
	signature := make([]bool, size)

	return &Signature{signature}
}

type Signature struct {
	signature []bool
}

func (s Signature) ttt() {}

func (s *Signature) Set(idx int) error {
	if idx < 0 || idx >= len(s.signature) {
		return fmt.Errorf("idx is out of bounds")
	}

	s.signature[idx] = true

	return nil
}

func (s *Signature) SetAll() error {
	for idx := 0; idx < len(s.signature); idx++ {
		err := s.Set(idx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Signature) Test(idx int) (bool, error) {
	if idx < 0 || idx >= len(s.signature) {
		return false, fmt.Errorf("idx is out of bounds")
	}
	return s.signature[idx], nil
}

func (s *Signature) Reset(idx int) error {
	if idx < 0 || idx >= len(s.signature) {
		return fmt.Errorf("idx is out of bounds")
	}

	s.signature[idx] = false

	return nil
}

func (s *Signature) ResetAll() error {
	for idx := 0; idx < len(s.signature); idx++ {
		err := s.Reset(idx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Signature) String() string {
	var b strings.Builder

	for idx := len(s.signature) - 1; idx >= 0; idx-- {
		if s.signature[idx] {
			fmt.Fprint(&b, "1")
		} else {
			fmt.Fprint(&b, "0")
		}
	}

	return b.String()
}

func (s *Signature) isValidIdx(idx int) error {
	if idx < 0 || idx >= len(s.signature) {
		return fmt.Errorf("idx is out of bounds")
	}
	return nil
}

func (s *Signature) IsEmpty() bool {
	if len(s.signature) > 0 {
		return true
	} else {
		return false
	}
}

func (s *Signature) Int() int {
	out := 0

	for idx := 0; idx < len(s.signature); idx++ {
    if s.signature[idx] {
      out |= int(math.Pow(2, float64(idx)))
    }
	}

	return out
}
