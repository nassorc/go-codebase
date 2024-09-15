package sparse_set

import "fmt"

const (
  PAGE_SIZE = 10
  MAX_PAGES = 1000
)

type Page [PAGE_SIZE]int
type Pages [MAX_PAGES]*Page

// Stores data that maps to a set of integers of the range 0..cap-1.
type SparseSet[T any] struct {
  size          int
  cap           int
  Dense         []T
  // Sparse's index maps to the set of integers, 0..cap-1, and its value gives the index
  // to its data in the Dense array.
  Pages         Pages // sparse set

  // ReverseLookup mirrors the Dense array, but contains the integer that owns the data.
  ReverseLookup []int // data index to id
}

func NewSparseSet[T any](cap int) *SparseSet[T] {
  var pages Pages

  for idx := 0; idx < MAX_PAGES; idx++ {
    pages[idx] = nil
  }

  return &SparseSet[T]{
    size: 0,
    cap: cap,
    Pages: pages,
    ReverseLookup: make([]int, cap, cap),
  }
}


func (s SparseSet[T]) panicInvalidIdx(idx int) {
  if idx < 0 || idx >= s.cap {
    panic(fmt.Sprintf("sparse set index %d out of bounds [0:%d].", idx, s.cap))
  }
}


func (s SparseSet[T]) Index(id int) any {
  return &s.Dense[id]
}

func (s SparseSet[T]) Get(id int) (T, bool) {
  if !s.Has(id) {
    var zero T
    return zero, false
  }

  row, col := intToPage(id)
  return s.Dense[s.Pages[row][col]], true
}

func (s SparseSet[T]) Has(id int) bool {
  row, col := intToPage(id)
  page := s.Pages[row]

  if page == nil {
    return false
  }

  idx := page[col]

  return s.ReverseLookup[idx] == id && idx < s.size
}

func intToPage(idx int) (row, col int) {
  row = idx / PAGE_SIZE
  col = idx % PAGE_SIZE
  return
}

func (s *SparseSet[T]) Insert(id int, value T) {
  s.panicInvalidIdx(id)
  if s.size >= s.cap {
		panic("Full component store.")
	}

  row, col := intToPage(id)

  // if page does NOT exist, create page
  if s.Pages[row] == nil {
    s.Pages[row] = new(Page)
  }

  if s.Has(id) {
    idx := s.Pages[row][col]
    s.Dense[idx] = value
  } else {
    // set page value to idx
    idx := s.size

    // actual data is equal or larger that the current storage capacity
    if s.size >= len(s.Dense) {
      s.Dense = append(s.Dense, value)
    } else {
      s.Dense[idx] = value
    }
    s.Pages[row][col] = idx
		s.ReverseLookup[idx] = id

		s.size += 1
  }
}

func (s *SparseSet[T]) Remove(id int) bool {
  s.panicInvalidIdx(id)
  if !s.Has(id) || s.size == 0 {
    return false
  }

	// idx := s.Sparse[id]
  row, col := intToPage(id)
	idx := s.Pages[row][col]

	lastIdx := s.size - 1
	lastOwnerId := s.ReverseLookup[lastIdx]

	// replace target data with the last element and create new slice excluding the value
  s.Dense[idx] = s.Dense[lastIdx] // swap

	// update data
  lrow, lcol := intToPage(lastOwnerId)
  s.Pages[lrow][lcol] = idx
	s.ReverseLookup[idx] = lastOwnerId
  // s.Pages
	s.size -= 1

	return true
}

func (s SparseSet[T]) Size() int {
  return s.size
}

func (s SparseSet[T]) Cap() int {
  return s.cap
}
