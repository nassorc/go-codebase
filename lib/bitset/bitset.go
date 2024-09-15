package bitset

import (
	"fmt"
	"strings"
)

func NewBitset(size int) Bitset {
	signature := make([]bool, size)

	return Bitset{ size, signature }
}

type Bitset struct {
  size int
	data []bool
}

func (s *Bitset) Count() int {
	count := 0

	for idx := 0; idx < s.size; idx++ {
		if s.data[idx] {
      count += 1
    }
  }

  return count
}

func (s *Bitset) IsSubset(other Bitset) bool {
	return (s.Int() & other.Int()) == s.Int()
}

func (s *Bitset) PanicIfNotValidIdx(idx int) {
	if idx < 0 || idx >= s.size {
		panic(fmt.Errorf("idx is out of bounds"))
	}
}

func (s *Bitset) Int() int {
	sum := 0

  // convert the bit array to an int by iterating over the array of bits and 
  // adding the bit value at 2^idx to sum if it is set to true.
	for idx := 0; idx < s.size; idx++ {
		if s.data[idx] {
      sum |= 1 << idx
			// out |= int(math.Pow(2, float64(idx)))
		}
	}

	return sum
}

func (s *Bitset) Reset(idx int) {
  s.PanicIfNotValidIdx(idx)
	s.data[idx] = false
}

func (s *Bitset) ResetAll() {
	for idx := 0; idx < s.size; idx++ {
		s.Reset(idx)
	}
}

func (s *Bitset) Set(idx int) {
  s.PanicIfNotValidIdx(idx)
	s.data[idx] = true
}

func (s *Bitset) SetAll() {
	for idx := 0; idx < s.size; idx++ {
		s.Set(idx)
	}
}

func (s *Bitset) String() string {
	var b strings.Builder

	for idx := s.size - 1; idx >= 0; idx-- {
		if s.data[idx] {
			fmt.Fprint(&b, "1")
		} else {
			fmt.Fprint(&b, "0")
		}
	}

	return b.String()
}

func (s *Bitset) Test(idx int) bool {
  s.PanicIfNotValidIdx(idx)
	return s.data[idx]
}
