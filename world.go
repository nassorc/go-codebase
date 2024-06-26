package gandalf

import (
	"fmt"
	"reflect"
)

const SIG_SIZE = 16

func CreateWorld(engine *Engine, size int) *World {
	var entityMgr = NewEntityManager(size)
	var systemMgr = NewSystemManager()
	var componentMgr = NewComponentManager(size)

	return &World{
		engine,
		entityMgr,
		systemMgr,
		componentMgr,
	}
}

type World struct {
	engine       *Engine
	entityMgr    *EntityManager
	systemMgr    *SystemManager
	componentMgr *ComponentManager
}

func (world *World) RegisterSystem(system System, components ...interface{}) {
	var sig = NewSignature(SIG_SIZE)

	// create system signature
	for _, component := range components {
		var t = reflect.TypeOf(component)
		var id, _ = world.componentMgr.GetStoreId(t)
		sig.Set(id)
	}

	world.systemMgr.Register(system, sig)
}

func (w *World) RegisterComponents(components ...interface{}) {
	for _, component := range components {
		t := reflect.TypeOf(component)

		if t.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		w.componentMgr.NewStore(t)
	}
}

func (world *World) CreateEntity(components ...interface{}) EntityHandle {

	var eSignature = NewSignature(SIG_SIZE)
	var entity = world.entityMgr.CreateEntity(eSignature)

	for _, component := range components {
		ok := world.componentMgr.AddDataToStore(entity, component)
		if !ok {
			panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
		}
		t := reflect.TypeOf(component)
		storeId, ok := world.componentMgr.GetStoreId(t)

		if !ok {
			panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
		}

		eSignature.Set(storeId)
	}

	var entityHdl = NewEntityHandle(entity, world, eSignature)
	world.systemMgr.NewEntity(entityHdl)

	return entityHdl
}

func (world *World) RemoveEntity(entity EntityId) {
	world.entityMgr.ScheduleEntityRemoval(entity)
}

func (world *World) GetDeadEntities() []EntityId {
	return world.entityMgr.GetEntitiesToRemove()
}

func (world *World) GetEntitySignature(entity EntityId) *Signature {
	return world.entityMgr.GetSignature(entity)
}

func (world *World) Tick() {
	// Update systems first. Updating entityMgr clears the entitiesToRemove array,
	// which the system manager uses to remove the entities from its store.
	// Calling system update first guarantees entities are removed.

	world.systemMgr.OnRemove(world)
	world.componentMgr.OnRemove(world)
	world.entityMgr.OnRemove(world)

	world.entityMgr.Update(world)
	world.componentMgr.Update(world)
	world.systemMgr.Update(world)
}

func (w *World) ChangeScene(game Scene) {
	w.engine.ChangeScene(game)
}
