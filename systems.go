package gandalf

import "github.com/hajimehoshi/ebiten/v2"

type System func([]EntityHandle)
type RSystem func(*ebiten.Image, []EntityHandle)

func NewSystemManager() *SystemManager {
	return &SystemManager{
		entityStore:      make(map[string][]EntityHandle),
		storeToSignature: make(map[string]Signature),
	}
}

type SystemManager struct {
	systems          []System
	systemSignatures []Signature

	renderers        []RSystem
	renderSignatures []Signature

	entityStore      map[string][]EntityHandle
	storeToSignature map[string]Signature
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
	_, ok := mgr.entityStore[signature.String()]
	if !ok {
		mgr.entityStore[signature.String()] = make([]EntityHandle, 0)
		mgr.storeToSignature[signature.String()] = signature
	}
}

func (mgr *SystemManager) NewEntity(entity EntityHandle) {
	// add entity to store
	for key := range mgr.entityStore {
		var storeSig = mgr.storeToSignature[key]
		var entitySig = entity.Signature().Int()

		if (storeSig.Int() & entitySig) == storeSig.Int() {
			mgr.entityStore[key] = append(mgr.entityStore[key], entity)
		}
	}
}

func (mgr *SystemManager) RemoveEntity(entity EntityId, entitySig *Signature) {
	for key := range mgr.entityStore {
		var storeSig = mgr.storeToSignature[key]
		if (storeSig.Int() & entitySig.Int()) == storeSig.Int() {
			for idx := 0; idx < len(mgr.entityStore[key]); idx++ {
				// find entities idx in store
				if entity == mgr.entityStore[key][idx].Entity() {
					// swap and pop with lastIdx
					lastIdx := len(mgr.entityStore[key]) - 1
					mgr.entityStore[key][idx] = mgr.entityStore[key][lastIdx]
					mgr.entityStore[key] = mgr.entityStore[key][0:lastIdx]
					break
				}
			}
		}
	}
}

func (mgr *SystemManager) OnRemove(world *World) {
	// remove dead entities
	for _, entity := range world.GetDeadEntities() {
		var sig = world.GetEntitySignature(entity)
		mgr.RemoveEntity(entity, sig)
	}
}

func (mgr *SystemManager) Update() {
	// call systems
	for idx, system := range mgr.systems {
		var signature = mgr.systemSignatures[idx]
		var entities = mgr.entityStore[signature.String()]
		system(entities)
	}
}

func (mgr *SystemManager) Render(screen *ebiten.Image) {
	// call systems
	for idx, renderer := range mgr.renderers {
		var signature = mgr.renderSignatures[idx]
		var entities = mgr.entityStore[signature.String()]
		renderer(screen, entities)
	}
}
