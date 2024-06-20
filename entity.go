package gandalf

import (
	"reflect"
)

type EntityId = int

type Entity struct {
	id        EntityId
	Signature *Signature
}

func (e *Entity) Id() EntityId {
	return e.id
}

func NewEntityHandle(world *World, entity *Entity) *EntityHandle {
	return &EntityHandle{
		entity,
		world,
	}
}

type EntityHandle struct {
	entity *Entity
	world  *World
}

// func (e *EntityHandle) UnpackToHandle(components ...IComponentHandle) {
// for _, _component := range components {
// 	// argType := component.TypeArg()
// }
// }

func (e *EntityHandle) Unpack(components ...interface{}) {
	for _, component := range components {
		t := reflect.TypeOf(component)
		val := reflect.ValueOf(component).Elem()

		if t.Kind() != reflect.Pointer {
			panic("Add component failed. Component is not a pointer type.")
		}

		cIdx, ok := e.world.typeToComponent[t.Elem()]

		if !ok {
			continue
		}

		carr := e.world.Components[cIdx]
		idx, ok := carr.entityToData[e.entity.Id()]

		if !ok {
			continue
		}

		val.Set(carr.Data.Index(idx).Addr().Elem())
	}
}

func (e *EntityHandle) Destroy() {}

func newEntityManager(size int) *EntityManager {
	var availIds = NewRingBuffer[int](size)
	var entityToSignature = make(map[EntityId]int)
	var signatureToEntity = make(map[int]EntityId)

	for idx := 0; idx < size; idx++ {
		availIds.Enqueue(idx)
	}

	return &EntityManager{
		availIds:          availIds,
		entityToSignature: entityToSignature,
		signatureToEntity: signatureToEntity,
	}
}

type EntityManager struct {
	availIds          *Ringbuffer[int]
	signatures        []*Signature
	entityToSignature map[EntityId]int
	signatureToEntity map[int]EntityId
}

func (manager *EntityManager) newEntity() (*Entity, bool) {
	id, ok := manager.availIds.Deque()

	if !ok {
		return nil, false
	}

	// entity's signature bookkeeping
	signatureIdx := len(manager.signatures)
	manager.signatures = append(manager.signatures, nil)
	manager.entityToSignature[id] = signatureIdx
	manager.signatureToEntity[signatureIdx] = id

	return &Entity{
		id: id,
	}, true
}

func (manager *EntityManager) setSignature(entity *Entity, signature *Signature) bool {
	idx, ok := manager.entityToSignature[entity.Id()]

	if !ok {
		return false
	}

	manager.signatures[idx] = signature

	return true
}

func (manager *EntityManager) getSignature(entity *Entity) (*Signature, bool) {
	idx, ok := manager.entityToSignature[entity.Id()]

	if !ok {
		return nil, false
	}

	return manager.signatures[idx], true
}

func (manager *EntityManager) removeEntity(entity *Entity) {
	id := entity.Id()
	// reset signature
	sigantureIdx := manager.entityToSignature[id]
	manager.signatures[sigantureIdx].ResetAll()
	manager.signatures[sigantureIdx] = nil

	// bookkeeping
	// swap and pop
	// swap last signtaure with deleted signature's values
	lastIdx := len(manager.signatures) - 1
	lastEntityId := manager.signatureToEntity[lastIdx]

	// swap
	manager.signatures[sigantureIdx] = manager.signatures[lastIdx]
	manager.signatures[lastIdx] = nil

	// update entities signature position
	manager.entityToSignature[lastEntityId] = sigantureIdx
	manager.signatureToEntity[sigantureIdx] = lastEntityId

	// remove deleted entity's records
	delete(manager.entityToSignature, id)
	delete(manager.signatureToEntity, lastIdx)

	// make removed entity's id available again
	manager.availIds.Enqueue(id)
}

func (manager *EntityManager) empty() bool {
	return manager.availIds.Empty()
}
