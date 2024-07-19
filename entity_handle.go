package gandalf

func NewEntityHandle(entity EntityId, world *World, signature Signature) EntityHandle {
	return EntityHandle{
		world,
		entity,
		signature,
	}
}

type EntityHandle struct {
	world     *World
	entity    EntityId
	signature Signature
}

func (e *EntityHandle) Entity() EntityId {
	return e.entity
}

func (e *EntityHandle) Unpack(components ...interface{}) {
	for _, component := range components {
		e.world.componentMgr.Unpack(e.entity, component)
	}
}

func (e *EntityHandle) Signature() Signature {
	return e.signature
}

func (e *EntityHandle) Destroy() {
	e.world.RemoveEntity(e.Entity())
}
