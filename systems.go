package gandalf

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type System func([]EntityHandle)
type RSystem func(*ebiten.Image, []EntityHandle)

func NewSystemManager() *SystemManager {
	return &SystemManager{
		stores:   make(map[string]*EntityStore),
		storeSig: make(map[string]Signature),
	}
}

type SystemManager struct {
	systems          []System
	systemSignatures []Signature

	renderers        []RSystem
	renderSignatures []Signature

	stores   map[string]*EntityStore
	storeSig map[string]Signature
}

func (mgr *SystemManager) Register(system System, signature Signature) {
	// create store
	mgr.CreateStore(signature)

	mgr.systems = append(mgr.systems, system)
	mgr.systemSignatures = append(mgr.systemSignatures, signature)
}

func (mgr *SystemManager) RegisterRenderer(renderer RSystem, signature Signature) {
	mgr.CreateStore(signature)

	mgr.renderers = append(mgr.renderers, renderer)
	mgr.renderSignatures = append(mgr.renderSignatures, signature)
}

func (mgr *SystemManager) CreateStore(signature Signature) {
	_, ok := mgr.stores[signature.String()]
	// _, ok := mgr.entityStore[signature.String()]
	if !ok {
		mgr.stores[signature.String()] = NewEntityStore()
		mgr.storeSig[signature.String()] = signature
	}
}

func (mgr *SystemManager) NewEntity(entity EntityHandle) {
	// add entity to store
	for key := range mgr.stores {
		var storeSig = mgr.storeSig[key]
		var entitySig = entity.Signature()

		// check if intersection of store and entity signatures is store signature
		if (storeSig.Int() & entitySig.Int()) == storeSig.Int() {
			store := mgr.stores[key]
			store.Insert(entity)
		}
	}
}

func (mgr *SystemManager) OnRemove(world *World) {
	// remove dead entities
	for _, entity := range world.GetDeadEntities() {
		var eSig = world.GetEntitySignature(entity)

		// loop through each store
		for key := range mgr.stores {
			var storeSig = mgr.storeSig[key]
			var store = mgr.stores[key]
			// entity signature intersects store signature and store has entity
			if (storeSig.Int()&eSig.Int()) == storeSig.Int() && store.Has(entity) {
				// performs swap and pop with last element to remove entity
				// get position
				idx := store.EntityToIdxLookup[entity]
				lastIdx := store.size - 1
				lastOwnerId := store.idxToEntityLookup[lastIdx]

				// swap ----------------------------------------------------------
				store.Entities[idx], store.Entities[lastIdx] = store.Entities[lastIdx], store.Entities[idx]

				// bookkeeping ---------------------------------------------------
				store.EntityToIdxLookup[lastOwnerId] = idx // owner to idx position
				store.idxToEntityLookup[idx] = lastOwnerId // idx position to owner

				// pop ------------------------------------------------------------
				store.size -= 1
			}
		}
	}
}

func (mgr *SystemManager) Update() {
	// call systems
	for idx, system := range mgr.systems {
		var signature = mgr.systemSignatures[idx]
		var store = mgr.stores[signature.String()]
		var entities = store.Entities
		system(entities[:store.size])
	}
}

func (mgr *SystemManager) Render(screen *ebiten.Image) {
	// call systems
	for idx, renderer := range mgr.renderers {
		var signature = mgr.renderSignatures[idx]
		var store = mgr.stores[signature.String()]
		var entities = store.Entities

		renderer(screen, entities[:store.size])
	}
}

func NewEntityStore() *EntityStore {
	return &EntityStore{
		Entities:          make([]EntityHandle, 10),
		EntityToIdxLookup: make([]int, 10),
		idxToEntityLookup: make([]int, 10),
	}
}

// ! Duplicate concept: ComponentStore
type EntityStore struct {
	Entities          []EntityHandle // dense set
	EntityToIdxLookup []int          // sparse set
	idxToEntityLookup []int          // reverse set: Given index i, map owner id to Entities.
	size              int
}

func (s EntityStore) Has(id EntityId) bool {
	return s.idxToEntityLookup[s.idxToEntityLookup[id]] == id && s.idxToEntityLookup[id] < s.size
}

func (s *EntityStore) Insert(entity EntityHandle) bool {
	id := entity.Entity()
	if s.Has(id) {
		return true
	}

	// add entity to list
	s.Entities[s.size] = entity

	// bookkeeping
	s.EntityToIdxLookup[id] = s.size
	s.idxToEntityLookup[s.size] = id

	s.size += 1

	return false
}
