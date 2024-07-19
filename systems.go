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

// 0010
func (mgr *SystemManager) OnSignatureEntityChange(entity EntityHandle) {
	for key := range mgr.stores {
		var storeSig = mgr.storeSig[key]
		var store = mgr.stores[storeSig.String()]
		var eSig = entity.Signature()

		// if new entity signature matches a system signature, add entity to system entity list
		if !store.Has(entity.Entity()) && (storeSig.Int()&eSig.Int()) == storeSig.Int() {
			store.Insert(entity)
		}
		if store.Has(entity.Entity()) && (storeSig.Int()&eSig.Int()) != storeSig.Int() {
			// fmt.Println("removing")
			store.Remove(entity.Entity())
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
				store.Remove(entity)
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

const HARDCODED_STORE_SIZE = 10

func NewEntityStore() *EntityStore {
	return &EntityStore{
		Entities:          make([]EntityHandle, HARDCODED_STORE_SIZE),
		EntityToIdxLookup: make([]int, HARDCODED_STORE_SIZE),
		idxToEntityLookup: make([]int, HARDCODED_STORE_SIZE),
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
	return s.idxToEntityLookup[s.EntityToIdxLookup[id]] == id && s.idxToEntityLookup[id] < s.size
}
func (s *EntityStore) Remove(id EntityId) bool {
	if !s.Has(id) {
		return false
	}
	idx := s.EntityToIdxLookup[id]
	lastIdx := s.size - 1
	lastOwnerId := s.idxToEntityLookup[lastIdx]

	// swap ----------------------------------------------------------
	s.Entities[idx], s.Entities[lastIdx] = s.Entities[lastIdx], s.Entities[idx]

	// bookkeeping ---------------------------------------------------
	s.EntityToIdxLookup[lastOwnerId] = idx // owner to idx position
	s.idxToEntityLookup[idx] = lastOwnerId // idx position to owner

	// pop ------------------------------------------------------------
	s.size -= 1

	return true
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
