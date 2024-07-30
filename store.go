package gandalf

import (
	"fmt"
	"reflect"
)

type DataStore[T any] struct {
	data          []T   // dense set
	idToIdxLookup []int // sparse set
	idxToIdLookup []int // reverse lookup: data idx to id
	size          int
}

func (s DataStore[T]) Has(id int) bool {
	return s.idxToIdLookup[s.idToIdxLookup[id]] == id && s.idxToIdLookup[id] < s.size
}

func (s *DataStore[T]) Insert(entity EntityHandle) error {
	return nil
	// id := entity.Entity()

	// if !s.Has(id) {
	// 	// add entity to list
	// 	s.data[s.size] = entity
	// }

	// // bookkeeping
	// s.EntityToIdxLookup[id] = s.size
	// s.idxToEntityLookup[s.size] = id

	// s.size += 1

	// return nil
}

func (s *DataStore[T]) Remove(id int) error {
	if !s.Has(id) {
		t := reflect.TypeOf((*T)(nil))
		return fmt.Errorf("record with id %v does not exist in store %v", id, t)
	}

	idx := s.idToIdxLookup[id]
	lastIdx := s.size - 1
	lastOwnerId := s.idxToIdLookup[lastIdx]

	// swap ----------------------------------------------------------
	s.data[idx], s.data[lastIdx] = s.data[lastIdx], s.data[idx]

	// bookkeeping ---------------------------------------------------
	s.idToIdxLookup[lastOwnerId] = idx // owner to idx position
	s.idxToIdLookup[idx] = lastOwnerId // idx position to owner

	// pop ------------------------------------------------------------
	s.size -= 1

	return nil
}
