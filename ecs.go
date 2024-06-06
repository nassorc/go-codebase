package gandalf

import (
	"fmt"
	"reflect"

	// rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_SIGNATURE_SIZE = 16

// Resources
// https://austinmorlan.com/posts/entity_component_system/#demo
// https://github.com/yohamta/donburi
// - uses unsafe pointers as the underlying type to store the component array
// - not recommended, "Packages that import unsafe may be non-portable and are not protected by the Go 1 compatibility guidelines."

// https://github.com/sedyh/mizu/
// https://github.com/ecsyjs/ecsy/tree/dev
// - both libraries are similar in how they define systems
// - mizu: uses the reflect library instead of unsafe pointers
// - mizu: requires systems to implement an interface and the game engine will
//   call the system's interface functions.
// - mizu: defines its signature or the entities it's interested in through its
//   struct fields. E.g. type PhysicsSystem struct { pos *Position }. The game
//   will query all entities with a position and set the system's field with
//   the current entity.


// How systems will be implemented:
// - Mizu's system implemenation works by defining a system as a struct with 
//   with actual component type as its fields. The struct's fields describes 
//   what entities the system is interested in. For example,
/*
type PhysicsSystem struct {
  pos *Position
  bbox *BoundingBox
}
*/
//   this mean that the system wants entities with a Position and a BoundingBox
//   component. The game engine will then loop through each entity, and for each
//   entity that matches the system, the game engine will set the system's field
//   to the entity's data.
//   After setting the value, the engine will then call system methods such as
//   Update() and Render(), which has access to the current entities data.

// - our implementation will instead have each system manage a list of entities
//   it's interested in. Through system methods such as Update() and Render(),
//   the game engine will loop through each system calling these method, and in
//   turn these methods will loop through each entity.
// - adding the entity to the list happens when the world.NewEntity() function is
//   called.

type ISystem interface {
  Update(w *World)
  AddEntity(e *Entity)
}


type Entity struct {
  id        int
  world     *World
  signature *Signature
}

func (e *Entity) GetData(components ...interface{}) {
  for _, component := range components {
    t := reflect.TypeOf(component)
    val := reflect.ValueOf(component).Elem()

    if t.Kind() != reflect.Pointer {
      panic("Add component failed. Component is not a pointer type.")
    }

    carr := e.world.GetComponentArray(t.Elem())
    idx := carr.entityToData[e.id]
    // ptrVal.Set(arr.Index(0).Addr().Elem())

    // newPtr := reflect.New(carr.data.Index(idx).Elem().Type())
    // newPtr.Elem().Set(carr.data.Index(idx).Elem())
    // val.Set(newPtr)

    val.Set(carr.data.Index(idx).Addr().Elem())


    UNUSED(idx)
  }
}

func (e *Entity) setData(components ...interface{}) {
  for _, component := range components {
    t := reflect.TypeOf(component)
    val := reflect.ValueOf(component)

    if t.Kind() != reflect.Pointer {
      panic("Add component failed. Component is not a pointer type.")
    }

    carr := e.world.GetComponentArray(t)
    carr.SetData(e.id, val)
  }
}

func NewComponentArray(t reflect.Type) *ComponentArray {
  return &ComponentArray {
    data: reflect.MakeSlice(reflect.SliceOf(t), 0, 0),
    entityToData: make(map[int]int),
  }
}

type ComponentArray struct {
  data reflect.Value
  entityToData map[int]int
}

func (c *ComponentArray) AppendData(entityId int, value reflect.Value) {
  idx := c.data.Len()
  c.data = reflect.Append(c.data, value)
  c.entityToData[entityId] = idx
}

func (c *ComponentArray) SetData(entityId int, value reflect.Value) {
  idx := c.entityToData[entityId]
  c.data.Index(idx).Set(value)
}

func (c *ComponentArray) GetData(entityId int) reflect.Value {
    idx := c.entityToData[entityId]

    return c.data.Index(idx)
}

func (c *ComponentArray) RemoveEntity(entityId int) {
}

func NewWorld() *World {
  return &World{
    typeToComponent: make(map[reflect.Type]int),
    availIds: NewRingBuffer[int](10),

    systemIdxToSignature: make(map[int]*Signature),
  }
}

type World struct {
  components           []*ComponentArray
  typeToComponent      map[reflect.Type]int

  availIds             *ringbuffer[int]
  Entities             []*Entity

  systems              []*ISystem
  systemIdxToSignature map[int]*Signature
}

func (w *World) Run() {
  for _, system := range w.systems {
    (*system).Update(w)
  }
}

func (w *World) RegisterSystem(system ISystem, components ...interface{}) {
  // t := reflect.TypeOf(system)
  idx := len(w.systems)
  w.systems = append(w.systems, &system)
  sSignature := NewSignature(MAX_SIGNATURE_SIZE)
  for _, component := range components {
    id := w.GetComponentId(component)
    sSignature.Set(id)
  }

  w.systemIdxToSignature[idx] = sSignature
}

func (w *World) RegisterComponents(components ...interface{}) {
  for _, component := range components {
    t := reflect.TypeOf(component)

    if t.Kind() != reflect.Pointer {
      panic("Add component failed. Component is not a pointer type.")
    }

    carr := NewComponentArray(t)
    idx := len(w.components)
    w.components = append(w.components, carr)
    w.typeToComponent[t] = idx
  }
}

func (w *World) GetComponentId(component interface{}) int {
    t := reflect.TypeOf(component)

    if t.Kind() != reflect.Pointer {
      panic("Add component failed. Component is not a pointer type.")
    }

    id, ok := w.typeToComponent[t]

    if !ok {
      panic(fmt.Sprintf("GetComponentId panicked. type=%v is not a component array", t))
    }

    return id
}

func (w *World) NewEntity(components ...interface{}) *Entity {
  id, err := Dequeue(w.availIds)

  if err != nil {
    panic("Dequeing Entity Id paniced")
  }

  var eSignature = NewSignature(MAX_SIGNATURE_SIZE)

  for _, component := range components {
    t := reflect.TypeOf(component)
    val := reflect.ValueOf(component)

    if t.Kind() != reflect.Pointer {
      panic("Add component failed. Component is not a pointer type.")
    }

    carr := w.GetComponentArray(t)
    carr.AppendData(id, val)
    id := w.GetComponentId(component)

    eSignature.Set(id)
  }

  entity := &Entity{
    id: id,
    world: w,
    signature: eSignature,
  }

  w.Entities = append(w.Entities, entity)

  for idx, signature := range w.systems {
    sSignature := w.systemIdxToSignature[idx]

    if (eSignature.Int() & sSignature.Int()) == sSignature.Int() {
      fmt.Println("entity with signature", eSignature.String(), "matched", sSignature.String())
      (*signature).AddEntity(entity)
    }
  }


  return entity
}

func (w *World) GetComponentArray(t reflect.Type) *ComponentArray {
    cidx := w.typeToComponent[t]
    return w.components[cidx]
}

func UNUSED(x ...interface{}) {}
