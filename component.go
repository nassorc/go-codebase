package gandalf

import (
	"reflect"
)

func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		typeToComponent: make(map[reflect.Type]int),
	}
}

type ComponentManager struct {
	components      []*Store
	typeToComponent map[reflect.Type]int
}

func (mgr *ComponentManager) NewStore(component reflect.Type) {
	var store = NewStore(component)
	var idx = len(mgr.components)

	mgr.components = append(mgr.components, store)

	// bookkeeping
	mgr.typeToComponent[component] = idx
}

func (mgr *ComponentManager) GetData(id int, component interface{}) bool {
	var t = reflect.TypeOf(component).Elem()
	storeId, ok := mgr.GetStoreId(t)

	if !ok {
		return ok
	}

	var store = mgr.components[storeId]
	var val = reflect.ValueOf(component).Elem()
	var got, _ = store.Get(id)

	val.Set(got)

	return true
}

func (mgr *ComponentManager) AddDataToStore(id int, data interface{}) bool {
	var t = reflect.TypeOf(data)
	var val = reflect.ValueOf(data)

	if t.Kind() != reflect.Pointer {
		return false
	}

	var storeId, ok = mgr.GetStoreId(t)

	if !ok {
		return false
	}

	var store = mgr.components[storeId]
	store.Insert(id, val)

	return true
}

func (mgr *ComponentManager) RemoveData(entity EntityId, component reflect.Type) bool {
	var idx, ok = mgr.typeToComponent[component]

	if !ok {
		return false
	}

	var store = mgr.components[idx]
	store.Remove(entity)
	return true
}

func (mgr *ComponentManager) GetStoreId(component reflect.Type) (int, bool) {
	var id, ok = mgr.typeToComponent[component]

	if !ok {
		return 0, false
	}

	return id, true
}

func (mgr *ComponentManager) OnRemove(world *World) {
	for _, entity := range world.GetDeadEntities() {
		for _, store := range mgr.components {
			store.Remove(entity)
		}
	}
}

func (mgr *ComponentManager) Update(world *World) {
}

func NewStore(t reflect.Type) *Store {
	return &Store{
		Data:           reflect.MakeSlice(reflect.SliceOf(t), 0, 0),
		idToDataLookup: make(map[int]int),
		dataToIdLookup: make(map[int]int),
	}
}

type Store struct {
	Data           reflect.Value
	idToDataLookup map[EntityId]EntityId
	dataToIdLookup map[int]EntityId
}

func (s *Store) Size() EntityId {
	return s.Data.Len()
}

func (s *Store) Get(id EntityId) (reflect.Value, bool) {
	idx, ok := s.idToDataLookup[id]

	if !ok {
		return reflect.Value{}, false
	}

	return s.Data.Index(idx), true // ! or s.Data.Index(idx).Addr().Elem()
}

func (s *Store) Insert(id EntityId, value reflect.Value) {
	idx := s.Data.Len()
	s.Data = reflect.Append(s.Data, value)
	s.idToDataLookup[id] = idx
	s.dataToIdLookup[idx] = id
}

// This function removes the data of the given id by performing
// a move and pop with the last element.
func (s *Store) Remove(id EntityId) bool {
	idx, ok := s.idToDataLookup[id]
	lastIdx := s.Data.Len() - 1
	lastOwnerId := s.dataToIdLookup[lastIdx]

	if !ok {
		return false
	}

	// replace target data with the last element and create new slice excluding the value
	s.Data.Index(idx).Set(s.Data.Index(lastIdx))
	s.Data = s.Data.Slice(0, lastIdx)

	// bookkeeping

	// repositioned data
	s.idToDataLookup[lastOwnerId] = idx
	s.dataToIdLookup[idx] = lastOwnerId
	delete(s.dataToIdLookup, lastIdx)

	// removed data
	delete(s.idToDataLookup, id)

	return true
}
