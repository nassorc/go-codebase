package gandalf

import "reflect"

type IComponentHandle interface {
	TypeArg() reflect.Type
}

type ComponentHandle[T any] struct {
	owner *Entity
	data  T
}

func (h *ComponentHandle[T]) TypeArg() reflect.Type {
	handleType := reflect.TypeOf((*T)(nil))
	return handleType
}

func NewComponentArray(t reflect.Type) *ComponentArray {
	return &ComponentArray{
		Data:         reflect.MakeSlice(reflect.SliceOf(t), 0, 0),
		entityToData: make(map[int]int),
	}
}

type ComponentArray struct {
	Data         reflect.Value
	entityToData map[int]int
}

func (c *ComponentArray) AppendData(entityId int, value reflect.Value) {
	idx := c.Data.Len()
	c.Data = reflect.Append(c.Data, value)
	c.entityToData[entityId] = idx
}

func (c *ComponentArray) SetData(entityId int, value reflect.Value) {
	idx := c.entityToData[entityId]
	c.Data.Index(idx).Set(value)
}

func (c *ComponentArray) GetData(entityId int) reflect.Value {
	idx := c.entityToData[entityId]

	return c.Data.Index(idx)
}

func (c *ComponentArray) RemoveEntity(entityId int) {
}
