package gandalf

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
	Entity *Entity
	world  *World
}

// func (e *EntityHandle) UnpackToHandle(components ...IComponentHandle) {
// for _, _component := range components {
// 	// argType := component.TypeArg()
// }
// }

func (e *EntityHandle) Unpack(components ...interface{}) {
	for _, component := range components {
		e.world.componentManager.GetData(e.Entity.Id(), component)
	}
}

func (e *EntityHandle) Destroy() {
	// e.world.entityManager
}

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

	entitiesToRemove []*Entity
}

func (manager *EntityManager) scheduleEntityRemoval(entity *Entity) {
	manager.entitiesToRemove = append(manager.entitiesToRemove, entity)
}

func (manager *EntityManager) removeEntities() {
	// for _, entity := range manager.entitiesToRemove {

	// }
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
