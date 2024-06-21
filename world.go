package gandalf

import (
	"fmt"
	"reflect"
)

var MAX_SIGNATURE_SIZE = 16

type System func(*World, []*EntityHandle)

func createWorld(engine *Engine, size int) *World {
	return &World{
		entityManager:   newEntityManager(size),
		typeToComponent: make(map[reflect.Type]int),
		engine:          engine,
		systemSignature: make(map[int]*Signature),
	}
}

type World struct {
	engine *Engine

	entityManager *EntityManager

	Components      []*ComponentArray
	typeToComponent map[reflect.Type]int

	systems         []System
	systemEntities  [][]*EntityHandle
	systemSignature map[int]*Signature
}

func (w *World) RegisterSystem(system System, components ...interface{}) {
	// t := reflect.TypeOf(system)
	idx := len(w.systems)
	w.systems = append(w.systems, system)

	sSignature := NewSignature(MAX_SIGNATURE_SIZE)
	for _, component := range components {
		Id := w.getComponentId(component)
		sSignature.Set(Id)
	}

	w.systemSignature[idx] = sSignature
	w.systemEntities = append(w.systemEntities, make([]*EntityHandle, 0))
}

func (w *World) RegisterComponents(components ...interface{}) {
	for _, component := range components {
		t := reflect.TypeOf(component)

		if t.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		carr := NewComponentArray(t)
		idx := len(w.Components)
		w.Components = append(w.Components, carr)
		w.typeToComponent[t] = idx
	}
}

func (w *World) CreateEntity(components ...interface{}) *EntityHandle {
	var eSignature = NewSignature(MAX_SIGNATURE_SIZE)
	var entity, ok = w.entityManager.newEntity()

	if !ok {
		panic("No available entities")
	}

	w.entityManager.setSignature(entity, eSignature)

	for _, component := range components {
		t := reflect.TypeOf(component)
		val := reflect.ValueOf(component)

		if t.Kind() != reflect.Pointer {
			panic("Component is not a pointer type.")
		}

		carr := w.getComponentArray(t)
		carr.AppendData(entity.Id(), val)
		at := w.getComponentId(component)

		eSignature.Set(at)
	}

	for idx := 0; idx < len(w.systems); idx++ {
		sSignature := w.systemSignature[idx]

		if (eSignature.Int() & sSignature.Int()) == sSignature.Int() {
			w.systemEntities[idx] = append(w.systemEntities[idx], NewEntityHandle(w, entity))
		}
	}

	return NewEntityHandle(w, entity)
}

func (w *World) CreateEntityFromPrefab(prefab interface{}) *EntityHandle {
	value := reflect.ValueOf(prefab).Elem()

	var components []interface{}

	for idx := 0; idx < value.NumField(); idx++ {
		components = append(components, value.Field(idx).Interface())
	}

	entity := w.CreateEntity(components...)

	return entity

}

func (w *World) update() {
	for idx, system := range w.systems {
		system(w, w.systemEntities[idx])
	}
}

func (w *World) getComponentId(component interface{}) int {
	t := reflect.TypeOf(component)

	if t.Kind() != reflect.Pointer {
		panic("Add component failed. Component is not a pointer type.")
	}

	Id, ok := w.typeToComponent[t]

	if !ok {
		panic(fmt.Sprintf("getComponentId panicked. type=%v is not a component array", t))
	}

	return Id
}

func (w *World) getComponentArray(t reflect.Type) *ComponentArray {
	cidx := w.typeToComponent[t]
	return w.Components[cidx]
}

func (w *World) ChangeScene(game Scene) {
	w.engine.changeScene(game)
}
