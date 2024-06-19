package gandalf

import (
	"reflect"
)

type Entity struct {
	id        int
	World     *World
	Signature *Signature
}

func (e *Entity) Id() int {
	return e.id
}

func (e *Entity) GetData(components ...interface{}) {
	for _, component := range components {
		t := reflect.TypeOf(component)
		val := reflect.ValueOf(component).Elem()

		if t.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		cIdx, ok := e.World.typeToComponent[t.Elem()]

		if !ok {
			continue
		}

		carr := e.World.Components[cIdx]
		idx, ok := carr.entityToData[e.Id()]

		if !ok {
			continue
		}

		val.Set(carr.Data.Index(idx).Addr().Elem())
	}
}

func (e *Entity) SetData(components ...interface{}) {
	for _, component := range components {
		t := reflect.TypeOf(component)
		val := reflect.ValueOf(component)

		if t.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		carr := e.World.getComponentArray(t)
		carr.SetData(e.id, val)
	}
}
