package gandalf

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type System func([]EntityHandle)
type Renderer func(*ebiten.Image, []EntityHandle)

func NewSystemManager(cap int) *SystemManager {
	return &SystemManager{
		stores:         make(map[string]*EntityStore),
		storeSignature: make(map[string]Signature),
		cap:            cap,
	}
}

type SystemManager struct {
	systems          []System
	systemSignatures []Signature
	renderers        []Renderer
	renderSignatures []Signature
	stores           map[string]*EntityStore
	storeSignature   map[string]Signature
	cap              int
}

func (mgr *SystemManager) CreateStore(sig Signature) {
	_, ok := mgr.stores[sig.String()]

	if !ok {
		mgr.stores[sig.String()] = NewEntityStore(mgr.cap)
		mgr.storeSignature[sig.String()] = sig
	}
}

func (mgr *SystemManager) EntitySignatureChange(entity EntityId, sig Signature) {
	for key, store := range mgr.stores {
		var storeSig = mgr.storeSignature[key]

		// if signature matches system signature and entity is not a member of the system, insert
		if !store.Has(entity) && storeSig.IsSubset(sig) {
			store.Insert(entity)
		}
		// if store has entity but signatures no longer match, remove
		if store.Has(entity) && !storeSig.IsSubset(sig) {
			store.Remove(entity)
		}
	}
}

func (mgr *SystemManager) NewEntity(entity EntityId, sig Signature) {
	// add entity to store
	for key := range mgr.stores {
		var storeSig = mgr.storeSignature[key]

		// insert entity to store if the entity's signature is a subset of the store's signature
		if storeSig.IsSubset(sig) {
			store := mgr.stores[key]
			store.Insert(entity)
		}
	}
}

func (mgr *SystemManager) OnRemove(world *World) {
	// remove dead entities
	for _, entity := range world.DeadEntities() {
		var entitySig = world.EntitySignature(entity)

		// loop through each store
		for key := range mgr.stores {
			var storeSig = mgr.storeSignature[key]
			var store = mgr.stores[key]

			// store has entity and entity signature intersects store signature
			if store.Has(entity) && storeSig.IsSubset(entitySig) {
				store.Remove(entity)
			}
		}
	}
}

func (mgr *SystemManager) Register(system System, sig Signature) {
	mgr.CreateStore(sig)

	mgr.systems = append(mgr.systems, system)
	mgr.systemSignatures = append(mgr.systemSignatures, sig)
}

func (mgr *SystemManager) RegisterRenderer(renderer Renderer, sig Signature) {
	mgr.CreateStore(sig)

	mgr.renderers = append(mgr.renderers, renderer)
	mgr.renderSignatures = append(mgr.renderSignatures, sig)
}

func (mgr *SystemManager) Render(world *World, screen *ebiten.Image) {
	// call systems
	for idx, renderer := range mgr.renderers {
		var signature = mgr.renderSignatures[idx]
		var store = mgr.stores[signature.String()]
		var entities = store.Entities
		var out = make([]EntityHandle, store.size)

		//! Every tick recreates this
		for idx := 0; idx < store.size; idx++ {
			out[idx] = NewEntityHandle(world, entities[idx])
		}

		renderer(screen, out)
	}
}

func (mgr *SystemManager) Update(world *World) {
	// call systems
	for idx, system := range mgr.systems {
		var signature = mgr.systemSignatures[idx]
		var store = mgr.stores[signature.String()]

		var entities = store.Entities
		var out = make([]EntityHandle, store.size)

		// fmt.Println(mgr.)

		//! Every tick creates this
		for idx := 0; idx < store.size; idx++ {
			out[idx] = NewEntityHandle(world, entities[idx])
		}
		system(out)
	}
}

func NewEntityStore(cap int) *EntityStore {
	return &EntityStore{
		Entities:    make([]EntityId, cap),
		EntityToIdx: make([]int, cap),
		idxToEntity: make([]int, cap),
		cap:         cap,
	}
}

// ! Duplicate concept: ComponentStore
type EntityStore struct {
	Entities    []EntityId // dense set
	EntityToIdx []int      // sparse set
	idxToEntity []int      // reverse set: Given index i, map owner id to Entities.
	cap         int
	size        int
}

func (s EntityStore) Has(id EntityId) bool {
	return s.idxToEntity[s.EntityToIdx[id]] == id && s.EntityToIdx[id] < s.size
}

func (s *EntityStore) Insert(entity EntityId) bool {
	if s.size >= s.cap {
		return false
	}

	if s.size >= s.cap {
		return false
	}

	if s.Has(entity) {
		return true
	}

	// add entity to list
	s.Entities[s.size] = entity

	// bookkeeping
	s.EntityToIdx[entity] = s.size
	s.idxToEntity[s.size] = entity

	s.size += 1

	return true
}

func (s *EntityStore) Remove(entity EntityId) bool {
	if !s.Has(entity) {
		return false
	}

	idx := s.EntityToIdx[entity]
	lastIdx := s.size - 1
	lastOwnerId := s.idxToEntity[lastIdx]

	// swap ----------------------------------------------------------
	s.Entities[idx], s.Entities[lastIdx] = s.Entities[lastIdx], s.Entities[idx]

	// bookkeeping ---------------------------------------------------
	s.EntityToIdx[lastOwnerId] = idx // owner to idx position
	s.idxToEntity[idx] = lastOwnerId // idx position to owner

	// pop ------------------------------------------------------------
	s.size -= 1

	return true
}
