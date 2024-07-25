package gandalf

type EntityId = int

func NewEntityManager(cap int) *EntityManager {
	q := NewRingBuffer[EntityId](cap)

	for idx := 0; idx < cap; idx++ {
		q.Enqueue(idx)
	}

	return &EntityManager{
		idPool:           q,
		entitiesToRemove: make([]EntityId, 0),
		entitySignatures: make([]Signature, cap),
		alive:            make([]bool, cap),
		cap:              cap,
	}
}

type EntityManager struct {
	idPool           *Ringbuffer[EntityId]
	entitiesToRemove []EntityId
	entitySignatures []Signature
	alive            []bool
	cap              int
}

func (mgr *EntityManager) Alive(entity EntityId) bool {
	return mgr.alive[entity]
}

func (mgr *EntityManager) Create(sig Signature) EntityId {
	if mgr.idPool == nil {
		panic("uninitialized entity id pool")
	}

	entity, ok := mgr.idPool.Deque()

	if !ok {
		panic("failed to create entity")
	}

	mgr.alive[entity] = true
	mgr.entitySignatures[entity] = sig

	return entity
}

func (mgr *EntityManager) GetEntitiesToRemove() []EntityId {
	return mgr.entitiesToRemove[:]
}

func (mgr *EntityManager) OnRemove(world *World) {
	mgr.RemoveDeadEntities()
}

func (mgr *EntityManager) RemoveDeadEntities() {
	for _, entity := range mgr.entitiesToRemove {
		mgr.entitySignatures[entity].ResetAll()
		mgr.alive[entity] = false
		mgr.idPool.Enqueue(entity) // requeue id back to id pool
	}

	// clear queue
	mgr.entitiesToRemove = nil
}

func (mgr *EntityManager) SetSignature(entity EntityId, sig Signature) bool {
	if !mgr.alive[entity] {
		return false
	}

	mgr.entitySignatures[entity] = sig

	return true
}

func (mgr *EntityManager) Signature(entity EntityId) Signature {
	return mgr.entitySignatures[entity]
}

func (mgr *EntityManager) ScheduleEntityRemoval(entity EntityId) {
	mgr.entitiesToRemove = append(mgr.entitiesToRemove, entity)
}

// Update is called every game loop
func (mgr *EntityManager) Update(world *World) {}
