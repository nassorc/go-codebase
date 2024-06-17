package ecs

import (
	"fmt"
	"reflect"

	g "github.com/nassorc/gandalf"
)

type System func(*World, []*Entity)

func NewWorld() *World {
	rb := g.NewRingBuffer[int](100)

	for idx := 0; idx < 100; idx++ {
		g.Enqueue(rb, idx)
	}

	return &World{
		// ActionMap:            make(map[int]string),
		typeToComponent: make(map[reflect.Type]int),
		availIds:        rb,
		systemSignature: make(map[int]*g.Signature),
	}
}

type World struct {
	// ActionMap            map[int]string
	Components      []*ComponentArray
	typeToComponent map[reflect.Type]int
	availIds        *g.Ringbuffer[int]
	Entities        []*Entity

	systems         []System
	systemEntities  [][]*Entity
	systemSignature map[int]*g.Signature
}

func (w *World) Run() {
	// for _, system := range w.systems {
	// 	system(w, []*Entity{})
	// }
	w.Tick()
}

func (w *World) RegisterAction(key int, name string) {
	// w.ActionMap[key] = name
}

func (w *World) RegisterSystem(system System, components ...interface{}) {
	// t := reflect.TypeOf(system)
	idx := len(w.systems)
	w.systems = append(w.systems, system)

	sSignature := g.NewSignature(g.MAX_SIGNATURE_SIZE)
	for _, component := range components {
		Id := w.GetComponentId(component)
		sSignature.Set(Id)
	}

	w.systemSignature[idx] = sSignature
	w.systemEntities = append(w.systemEntities, make([]*Entity, 0))
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

func (w *World) GetComponentId(component interface{}) int {
	t := reflect.TypeOf(component)

	if t.Kind() != reflect.Pointer {
		panic("Add component failed. Component is not a pointer type.")
	}

	Id, ok := w.typeToComponent[t]

	if !ok {
		panic(fmt.Sprintf("GetComponentId panicked. type=%v is not a component array", t))
	}

	return Id
}

func (w *World) CreateEntityFromPreFab(prefab interface{}) *Entity {
	value := reflect.ValueOf(prefab).Elem()

	var components []interface{}

	for idx := 0; idx < value.NumField(); idx++ {
		// entity := w.CreateEntity()
		components = append(components, value.Field(idx).Interface())
	}

	entity := w.CreateEntity(components...)

	return entity
}

func (w *World) CreateEntity(components ...interface{}) *Entity {
	id, err := g.Dequeue(w.availIds)

	if err != nil {
		panic("Dequeing Entity Id paniced")
	}

	var eSignature = g.NewSignature(g.MAX_SIGNATURE_SIZE)

	fmt.Println("creating new entity")
	for _, component := range components {
		t := reflect.TypeOf(component)
		val := reflect.ValueOf(component)

		if t.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		carr := w.GetComponentArray(t)
		carr.AppendData(id, val)
		at := w.GetComponentId(component)

		eSignature.Set(at)
	}

	entity := &Entity{
		id:        id,
		World:     w,
		Signature: eSignature,
	}

	w.Entities = append(w.Entities, entity)

	for idx := 0; idx < len(w.systems); idx++ {
		sSignature := w.systemSignature[idx]

		if (eSignature.Int() & sSignature.Int()) == sSignature.Int() {
			w.systemEntities[idx] = append(w.systemEntities[idx], entity)
		}
	}

	return entity
}

func (w *World) RemoveEntity() {
	// remove data from component list
	// remove from entities
}

func (w *World) GetData(entity *Entity, components ...interface{}) {
	for _, component := range components {
		t := reflect.TypeOf(component)
		val := reflect.ValueOf(component).Elem()

		if t.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		cIdx, ok := entity.World.typeToComponent[t.Elem()]

		if !ok {
			continue
		}

		carr := entity.World.Components[cIdx]
		idx, ok := carr.entityToData[entity.Id()]

		if !ok {
			continue
		}

		val.Set(carr.Data.Index(idx).Addr().Elem())
	}
}

func (w *World) GetComponentArray(t reflect.Type) *ComponentArray {
	cidx := w.typeToComponent[t]
	return w.Components[cidx]
}

func (w *World) Tick() {
	for idx, system := range w.systems {
		system(w, w.systemEntities[idx])
	}
}

// func (w *World) ChangeScene(newScene g.Scene) {
// }

func UNUSED(x ...interface{}) {}
