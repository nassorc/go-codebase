package gandalf

import (
	"reflect"
)

func NewComponentManager(storeCapacity int) *ComponentManager {
	return &ComponentManager{
		capacity:        storeCapacity,
		typeToComponent: make(map[reflect.Type]int),
	}
}

type ComponentManager struct {
	capacity        int
	components      []*Store
	typeToComponent map[reflect.Type]int
}

func (mgr *ComponentManager) GetOwners(component reflect.Type) []EntityId {
	return mgr.components[mgr.typeToComponent[component]].GetOwners()
}

func (mgr *ComponentManager) NewStore(component reflect.Type) {
	var store = NewStore(component, mgr.capacity)
	var idx = len(mgr.components)

	mgr.components = append(mgr.components, store)

	// bookkeeping
	mgr.typeToComponent[component] = idx
}

func (mgr *ComponentManager) Unpack(id int, component interface{}) bool {
	var t = reflect.TypeOf(component).Elem()
	storeId, ok := mgr.GetStoreId(t)

	if !ok {
		return false
	}

	var store = mgr.components[storeId]
	var val = reflect.ValueOf(component).Elem()
	got, ok := store.Get(id)

	if !ok {
		return false
	}

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

	ok = store.Remove(entity)

	if !ok {
		panic("failed to remove data from component store.")
	}

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
	for _, entity := range world.DeadEntities() {
		for _, store := range mgr.components {
			store.Remove(entity)
		}
	}
}

func (mgr *ComponentManager) Update(world *World) {
}

func NewStore(t reflect.Type, capacity int) *Store {
	return &Store{
		capacity:       capacity,
		Data:           reflect.MakeSlice(reflect.SliceOf(t), 0, 0),
		idToDataLookup: make([]EntityId, capacity),
		dataToIdLookup: make([]EntityId, capacity),
	}
}

type Store struct {
	capacity       int
	size           int
	Data           reflect.Value
	dataToIdLookup []EntityId
	idToDataLookup []EntityId
}

func (s *Store) GetOwners() []EntityId {
	return s.dataToIdLookup[:s.Data.Len()]
}

func (s *Store) Size() EntityId {
	return s.Data.Len()
}

func (s *Store) Has(id EntityId) bool {
	if s.dataToIdLookup[s.idToDataLookup[id]] == id && s.idToDataLookup[id] < s.Size() {
		return true
	}
	return false
}

func (s *Store) Get(id EntityId) (reflect.Value, bool) {
	if !s.Has(id) {
		return reflect.Value{}, false
	}

	idx := s.idToDataLookup[id]

	return s.Data.Index(idx), true // ! or s.Data.Index(idx).Addr().Elem()
}

func (s *Store) Insert(id EntityId, value reflect.Value) {
	if s.Size() >= s.capacity {
		panic("Full component store.")
	}

	idx := s.Data.Len()
	s.Data = reflect.Append(s.Data, value)
	s.idToDataLookup[id] = idx
	s.dataToIdLookup[idx] = id
}

// This function removes the data of the given id by performing
// a move and pop with the last element.
func (s *Store) Remove(id EntityId) bool {
	if !s.Has(id) {
		return false
	}

	idx := s.idToDataLookup[id]
	lastIdx := s.Data.Len() - 1
	lastOwnerId := s.dataToIdLookup[lastIdx]

	// replace target data with the last element and create new slice excluding the value
	s.Data.Index(idx).Set(s.Data.Index(lastIdx))
	s.Data = s.Data.Slice(0, lastIdx)

	// bookkeeping

	// repositioned data
	s.idToDataLookup[lastOwnerId] = idx
	s.dataToIdLookup[idx] = lastOwnerId

	return true
}
