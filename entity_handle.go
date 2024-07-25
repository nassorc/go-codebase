package gandalf

func NewEntityHandle(world *World, entity EntityId) EntityHandle {
	return EntityHandle{entity, world}
}

type EntityHandle struct {
	id    EntityId
	world *World
}

func (e *EntityHandle) Add(component interface{}) {
	e.world.AddComponent(e.Id(), component)
}

func (e *EntityHandle) Destroy() {
	e.world.RemoveEntity(e.Id())
}

func (e *EntityHandle) Id() EntityId {
	return e.id
}

func (e *EntityHandle) Remove(component ComponentID) {
	e.world.RemoveComponent(e.Id(), component)
}

func (e *EntityHandle) Unpack(components ...interface{}) {
	for _, component := range components {
		e.world.componentMgr.Unpack(e.Id(), component)
	}
}

func (e *EntityHandle) Signature() Signature {
	return e.world.EntitySignature(e.Id())
}
