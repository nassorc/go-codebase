package sparse_set

import "fmt"

func NewSparseSet[T any](cap int) *SparseSet[T] {
  return &SparseSet[T]{
    size: 0,
    cap: cap,
    Dense: make([]T, cap, cap),
    Sparse: make([]int, cap, cap),
    ReverseLookup: make([]int, cap, cap),
  }
}

// Stores data that maps to a set of integers of the range 0..cap-1.
type SparseSet[T any] struct {
  size          int
  cap           int
  Dense         []T
  // Sparse's index maps to the set of integers, 0..cap-1, and its value gives the index
  // to its data in the Dense array.
  Sparse        []int   
  // ReverseLookup mirrors the Dense array, but contains the integer that owns the data.
  ReverseLookup []int // data index to id
}

func (s SparseSet[T]) panicInvalidIdx(idx int) {
  if idx < 0 || idx >= s.cap {
    panic(fmt.Sprintf("sparse set index %d out of bounds [0:%d].", idx, s.cap))
  }
}

func (s SparseSet[T]) Get(id int) (T, bool) {
  s.panicInvalidIdx(id)
  if !s.Has(id) {
    var zero T
    return zero, false
  }

  return s.Dense[s.Sparse[id]], true
}

func (s SparseSet[T]) Has(id int) bool {
  return s.ReverseLookup[s.Sparse[id]] == id && s.Sparse[id] < s.size
}

func (s *SparseSet[T]) Insert(id int, value T) {
  s.panicInvalidIdx(id)
  if s.size >= s.cap {
		panic("Full component store.")
	}

	if s.Has(id) {  // update
		idx := s.Sparse[id]
    s.Dense[idx] = value
	} else {  // insert
		idx := s.size

		s.Dense[idx] = value
		s.Sparse[id] = idx
		s.ReverseLookup[idx] = id

		s.size += 1
	}
}

func (s *SparseSet[T]) Remove(id int) bool {
  s.panicInvalidIdx(id)
  if !s.Has(id) || s.size == 0 {
    return false
  }

	idx := s.Sparse[id]

	lastIdx := s.size - 1
	lastOwnerId := s.ReverseLookup[lastIdx]

	// replace target data with the last element and create new slice excluding the value
  s.Dense[idx] = s.Dense[lastIdx] // swap

	// update data
	s.Sparse[lastOwnerId] = idx
	s.ReverseLookup[idx] = lastOwnerId
	s.size -= 1

	return true
}

func (s SparseSet[T]) Size() int {
  return s.size
}

func (s SparseSet[T]) Cap() int {
  return s.cap
}
