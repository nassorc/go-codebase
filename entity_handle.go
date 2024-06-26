package gandalf

func NewEntityHandle(entity EntityId, world *World, signature *Signature) EntityHandle {
	return EntityHandle{
		world,
		entity,
		signature,
	}
}

type EntityHandle struct {
	world     *World
	entity    EntityId
	signature *Signature
}

func (e *EntityHandle) Entity() EntityId {
	return e.entity
}

func (e *EntityHandle) Unpack(components ...interface{}) {
	for _, component := range components {
		// t := reflect.TypeOf(component)
		// val := reflect.ValueOf(component).Elem()

		// cIdx, ok := e.world.typeToComponent[t.Elem()]

		e.world.componentMgr.GetData(e.entity, component)

		// if !ok {
		// 	continue
		// }

		// carr := e.world.Components[cIdx]
		// idx, ok := carr.entityToData[e.entity.Id()]

		// if !ok {
		// 	continue
		// }

		// val.Set(carr.Data.Index(idx).Addr().Elem())
	}
}

func (e *EntityHandle) Signature() *Signature {
	return e.signature
}

func (e *EntityHandle) Destroy() {
	e.world.RemoveEntity(e.Entity())
}
