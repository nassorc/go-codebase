package gandalf

import (
	"fmt"
	"image"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

const SIG_SIZE = 16

func NewWorld(size int) *World {
	var entityMgr = NewEntityManager(size)
	var systemMgr = NewSystemManager()
	var componentMgr = NewComponentManager(size)
	var assetMgr = NewAssetManager()

	return &World{
		entityMgr,
		systemMgr,
		componentMgr,
		assetMgr,
	}
}

type World struct {
	entityMgr    *EntityManager
	systemMgr    *SystemManager
	componentMgr *ComponentManager
	assetMgr     *AssetManager
}

func (world *World) RegisterSystem(system System, components ...ComponentID) {
	var sig = NewSignature(SIG_SIZE)

	// create system signature
	for _, component := range components {
		var id, _ = world.componentMgr.GetStoreId(component)
		sig.Set(id)
	}

	world.systemMgr.Register(system, *sig)
}

func (world *World) RegisterRenderer(system RSystem, components ...ComponentID) {
	var sig = NewSignature(SIG_SIZE)

	// create system signature
	for _, component := range components {
		var id, _ = world.componentMgr.GetStoreId(component)
		sig.Set(id)
	}

	world.systemMgr.RegisterRenderer(system, *sig)
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

	world.assetMgr.Update()
	world.entityMgr.Update(world)
	world.componentMgr.Update(world)
	world.systemMgr.Update()
}

func (world *World) Draw(screen *ebiten.Image) {
	world.systemMgr.Render(screen)
}

func (w *World) LoadTexture(name string, img image.Image) error {
	return w.assetMgr.loadTexture(name, img)
}

func (w *World) LoadAnimation(
	animName string,
	textName string,
	totalFrames int,
	src image.Rectangle,
	frmSize Vec2,
	frmOffset Vec2,
	scale float32,
	rotation float32,
	speed float32,
) bool {
	return w.assetMgr.loadAnimation(animName, textName, totalFrames, src, frmSize, frmOffset, scale, rotation, speed)
}

func (w *World) GetTexture(name string) (*ebiten.Image, bool) {
	return w.assetMgr.getTexture(name)
}

func (w *World) GetAnimation(name string) (*Animation, bool) {
	return w.assetMgr.getAnimation(name)
}

func (w *World) Query(component ComponentID) []EntityId {
	return w.componentMgr.GetOwners(component)
}

func (w *World) NewEntityHandle(entity EntityId) EntityHandle {
	return NewEntityHandle(entity, w, w.entityMgr.entitySignatures[entity])
}
