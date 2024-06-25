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
	components      []*ComponentStore
	typeToComponent map[reflect.Type]int
}

func (mgr *ComponentManager) RemoveData(id int, component reflect.Type) bool {
	idx, ok := mgr.typeToComponent[component]

	if !ok {
		return false
	}

	store := mgr.components[idx]
	store.Remove(id)
	return true
}

func (mgr *ComponentManager) NewStore(component reflect.Type) {
	store := NewComponentStore(component)
	idx := len(mgr.components)

	mgr.components = append(mgr.components, store)

	// bookkeeping
	mgr.typeToComponent[component] = idx
}

func (mgr *ComponentManager) GetStoreId(component reflect.Type) (int, bool) {
	id, ok := mgr.typeToComponent[component]

	if !ok {
		return 0, false
	}

	return id, true
}

func (mgr *ComponentManager) GetData(id int, component interface{}) bool {
	t := reflect.TypeOf(component).Elem()
	storeId, ok := mgr.GetStoreId(t)

	if !ok {
		return ok
	}

	store := mgr.components[storeId]
	// data := store.Get(id)

	val := reflect.ValueOf(component).Elem()

	got, _ := store.Get(id)
	val.Set(got)
	return true
}

func (mgr *ComponentManager) GetDataWithHandle(id int, component IComponentHandle) bool {
	t := component.TypeArg().Elem()
	storeId, ok := mgr.GetStoreId(t)

	if !ok {
		return ok
	}

	store := mgr.components[storeId]
	val := reflect.ValueOf(component).Elem()

	data, _ := store.Get(id)
	val.Set(data)
	return true
}

func (mgr *ComponentManager) AddDataToStore(id int, data interface{}) bool {
	t := reflect.TypeOf(data)
	val := reflect.ValueOf(data)

	if t.Kind() != reflect.Pointer {
		return false
	}

	storeId, ok := mgr.GetStoreId(t)

	if !ok {
		return false
	}

	store := mgr.components[storeId]
	store.Push(id, val)

	return true
}

func NewComponentStore(t reflect.Type) *ComponentStore {
	return &ComponentStore{
		Data:           reflect.MakeSlice(reflect.SliceOf(t), 0, 0),
		idToDataLookup: make(map[int]int),
		dataToIdLookup: make(map[int]int),
	}
}

type ComponentStore struct {
	Data           reflect.Value
	idToDataLookup map[int]int
	dataToIdLookup map[int]int
}

// This function removes the data of the given id by performing
// a move and pop with the last element.
func (c *ComponentStore) Remove(id int) bool {
	idx, ok := c.idToDataLookup[id]
	lastIdx := c.Data.Len() - 1
	lastOwnerId := c.dataToIdLookup[lastIdx]

	if !ok {
		return false
	}

	// replace target data with the last element and create new slice excluding the value
	c.Data.Index(idx).Set(c.Data.Index(lastIdx))
	c.Data = c.Data.Slice(0, lastIdx)

	// bookkeeping

	// repositioned data
	c.idToDataLookup[lastOwnerId] = idx
	c.dataToIdLookup[idx] = lastOwnerId
	delete(c.dataToIdLookup, lastIdx)

	// removed data
	delete(c.idToDataLookup, id)

	return true
}

func (c *ComponentStore) Push(id int, value reflect.Value) {
	idx := c.Data.Len()
	c.Data = reflect.Append(c.Data, value)
	c.idToDataLookup[id] = idx
	c.dataToIdLookup[idx] = id
}

func (c *ComponentStore) Size() int {
	return c.Data.Len()
}

func (c *ComponentStore) Get(id int) (reflect.Value, bool) {
	idx, ok := c.idToDataLookup[id]

	if !ok {
		return reflect.Value{}, false
	}

	return c.Data.Index(idx).Addr().Elem(), true // ! or c.Data.Index(idx).Addr().Elem()
}

type IComponentHandle interface {
	TypeArg() reflect.Type
}

type ComponentHandle[T any] struct {
	World *World
	Owner *Entity
	Data  T
}

func (h *ComponentHandle[T]) ReflectType() reflect.Type {
	handleType := reflect.TypeOf((*T)(nil))
	return handleType
}

func (h *ComponentHandle[T]) Destroy() bool {
	// h.World.componentManager.
	return false
}
