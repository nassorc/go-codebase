package gandalf

import (
	"fmt"
	"reflect"
)

var MAX_SIGNATURE_SIZE = 16

type System func(*World, []*Entity)

func createWorld(engine *Engine) *World {
	rb := NewRingBuffer[int](100)

	for idx := 0; idx < 100; idx++ {
		Enqueue(rb, idx)
	}

	return &World{
		typeToComponent: make(map[reflect.Type]int),
		availIds:        rb,
		engine:          engine,
		systemSignature: make(map[int]*Signature),
	}
}

type World struct {
	engine *Engine

	Components      []*ComponentArray
	typeToComponent map[reflect.Type]int
	availIds        *Ringbuffer[int]
	Entities        []*Entity

	systems         []System
	systemEntities  [][]*Entity
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

func (w *World) CreateEntity(components ...interface{}) *Entity {
	id, err := Dequeue(w.availIds)

	if err != nil {
		panic("Dequeing Entity Id paniced")
	}

	var eSignature = NewSignature(MAX_SIGNATURE_SIZE)

	fmt.Println("creating new entity")
	for _, component := range components {
		t := reflect.TypeOf(component)
		val := reflect.ValueOf(component)

		if t.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		carr := w.getComponentArray(t)
		carr.AppendData(id, val)
		at := w.getComponentId(component)

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

func (w *World) update() {
	for _, system := range w.systems {
		system(w, make([]*Entity, 5))
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

func (w *World) ChangeScene(game Game) {
	fmt.Println("WHAT IS THIS", w.engine)
	w.engine.changeGame(game)
}
