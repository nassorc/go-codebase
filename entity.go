package gandalf

type EntityId = int

func NewEntityManager(size int) *EntityManager {
	q := NewRingBuffer[EntityId](size)

	for idx := 0; idx < size; idx++ {
		q.Enqueue(idx)
	}

	return &EntityManager{
		availEntities:    q,
		entitiesToRemove: make([]EntityId, 0),

		entitySignatures: make([]*Signature, size),
	}
}

type EntityManager struct {
	availEntities    *Ringbuffer[EntityId]
	entitiesToRemove []EntityId

	entitySignatures []*Signature
}

func (mgr *EntityManager) CreateEntity(signature *Signature) EntityId {
	newEntity, _ := mgr.availEntities.Deque()
	mgr.entitySignatures[newEntity] = signature

	return newEntity
}

func (mgr *EntityManager) GetEntitiesToRemove() []EntityId {
	return mgr.entitiesToRemove[:]
}

func (mgr *EntityManager) GetSignature(entity EntityId) *Signature {
	return mgr.entitySignatures[entity]
}

func (mgr *EntityManager) ScheduleEntityRemoval(entity EntityId) {
	mgr.entitiesToRemove = append(mgr.entitiesToRemove, entity)
}

func (mgr *EntityManager) RemoveDeadEntities() {
	for _, removingId := range mgr.entitiesToRemove {
		// mgr.entitySignatures[removingId].ResetAll()
		mgr.entitySignatures[removingId] = nil

		// requeue dead entity's id in the available entity queue
		mgr.availEntities.Enqueue(removingId)
	}

	// clear entitiesToRemove queue
	// var cpy = mgr.entitiesToRemove[:]
	mgr.entitiesToRemove = nil
}

func (mgr *EntityManager) OnRemove(world *World) {
	mgr.RemoveDeadEntities()
}

// Update is called every game loop
func (mgr *EntityManager) Update(world *World) {
	// mgr.RemoveDeadEntities()
}
