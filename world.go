package gandalf

import (
	"fmt"
	"reflect"
)

var MAX_SIGNATURE_SIZE = 16

type System func(*World, []*EntityHandle)

func createWorld(engine *Engine, size int) *World {
	return &World{
		entityManager:    newEntityManager(size),
		componentManager: NewComponentManager(),

		// typeToComponent: make(map[reflect.Type]int),
		engine:          engine,
		systemSignature: make(map[int]*Signature),
	}
}

type World struct {
	engine *Engine

	entityManager    *EntityManager
	componentManager *ComponentManager

	systems         []System
	systemEntities  [][]*EntityHandle
	systemSignature map[int]*Signature
}

func (w *World) RegisterSystem(system System, components ...interface{}) {
	// t := reflect.TypeOf(system)
	idx := len(w.systems)
	w.systems = append(w.systems, system)

	sSignature := NewSignature(MAX_SIGNATURE_SIZE)

	// create system signature
	for _, component := range components {
		t := reflect.TypeOf(component)
		id, ok := w.componentManager.GetStoreId(t)

		if !ok {
			panic("System cannot register component that does not exist.")
		}
		sSignature.Set(id)
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

		w.componentManager.NewStore(t)
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
		ok := w.componentManager.AddDataToStore(entity.Id(), component)
		if !ok {
			panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
		}
		t := reflect.TypeOf(component)
		storeId, ok := w.componentManager.GetStoreId(t)

		if !ok {
			panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
		}

		eSignature.Set(storeId)
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

func (w *World) RemoveEntity(entity *Entity) {
	// w.entityManager.
}

func (w *World) RemoveComponent(entity *Entity, component interface{}) bool {
	var entityId = entity.Id()
	var t = reflect.TypeOf(component)
	var storeId, ok = w.componentManager.GetStoreId(t)
	if !ok {
		panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
	}
	eSignature, ok := w.entityManager.getSignature(entity)
	if !ok {
		panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
	}

	ok = w.componentManager.RemoveData(entityId, t)
	if !ok {
		panic("Failed remove component")
	}

	// find all systems that match the entity's siganture and remove from the system
	for idx := 0; idx < len(w.systems); idx++ {
		sSignature := w.systemSignature[idx]
		entities := w.systemEntities[idx]

		// if entity signature does not match system
		if sSignature.Int()&eSignature.Int() != sSignature.Int() {
			continue
		}

		// find entity in array
		// !Find better solution========================
		var found = -1
		for idx = 0; idx < len(entities); idx++ {
			if entityId == entities[idx].Entity.Id() {
				found = idx
				break
			}
		}

		var lastIdx = len(entities) - 1

		if found >= 0 {
			entities[idx] = entities[lastIdx]
			w.systemEntities[idx] = entities[0:lastIdx]
		}
	}

	eSignature.Reset(storeId)

	return true
}

func (w *World) update() {
	for idx, system := range w.systems {
		system(w, w.systemEntities[idx])
	}
}

func (w *World) ChangeScene(game Scene) {
	w.engine.changeScene(game)
}
