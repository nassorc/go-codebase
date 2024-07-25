package gandalf

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

const SIG_SIZE = 16

func NewWorld(size int) *World {
	var entityMgr = NewEntityManager(size)
	var systemMgr = NewSystemManager()
	var componentMgr = NewComponentManager(size)

	return &World{
		entityMgr,
		systemMgr,
		componentMgr,
	}
}

type World struct {
	entityMgr    *EntityManager
	systemMgr    *SystemManager
	componentMgr *ComponentManager
}

func (world *World) RegisterSystem(system System, components ...ComponentID) {
	var sig = NewSignature(SIG_SIZE)

	// create system signature
	for _, component := range components {
		var id, _ = world.componentMgr.GetStoreId(component)
		sig.Set(id)
	}

	world.systemMgr.Register(system, sig)
}

func (world *World) RegisterRenderer(system RSystem, components ...ComponentID) {
	var sig = NewSignature(SIG_SIZE)

	// create system signature
	for _, component := range components {
		var id, _ = world.componentMgr.GetStoreId(component)
		sig.Set(id)
	}

	world.systemMgr.RegisterRenderer(system, sig)
}

func (w *World) RegisterComponents(components ...ComponentID) {
	for _, component := range components {
		if component.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		w.componentMgr.NewStore(component)
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

func (world *World) AddComponent(entity EntityId, component interface{}) {
	ok := world.componentMgr.AddDataToStore(entity, component)

	if !ok {
		panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
	}

	t := reflect.TypeOf(component)
	storeId, ok := world.componentMgr.GetStoreId(t)

	if !ok {
		panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
	}

	orgSig := world.entityMgr.GetSignature(entity)
	sig := NewSignature(SIG_SIZE)
	sig.signature = slices.Clone(orgSig.signature)

	sig.Set(storeId)

	world.systemMgr.OnSignatureEntityChange(NewEntityHandle(entity, world, sig))
}

func (world *World) RemoveComponent(entity EntityId, component ComponentID) {
	sig := world.entityMgr.GetSignature(entity)
	id, ok := world.componentMgr.GetStoreId(component)

	if !ok {
		panic(fmt.Sprintf("Component %v is not a store", reflect.TypeOf(component)))
	}

	ok = world.componentMgr.RemoveData(entity, component)

	if !ok {
		panic(fmt.Sprintf("failed to remove %v", reflect.TypeOf(component)))
	}

	world.entityMgr.UpdateSignature(entity, sig)
	sig.Reset(id)

	world.systemMgr.OnSignatureEntityChange(NewEntityHandle(entity, world, sig))
}

func (world *World) GetDeadEntities() []EntityId {
	return world.entityMgr.GetEntitiesToRemove()
}

func (world *World) GetEntitySignature(entity EntityId) Signature {
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
	world.systemMgr.Update()
}

func (world *World) Draw(screen *ebiten.Image) {
	world.systemMgr.Render(screen)
}

}

func (w *World) Query(component ComponentID) []EntityId {
	return w.componentMgr.GetOwners(component)
}

func (w *World) NewEntityHandle(entity EntityId) EntityHandle {
	return NewEntityHandle(entity, w, w.entityMgr.entitySignatures[entity])
}
