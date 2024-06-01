package main

import (
	"fmt"
	"reflect"
	"slices"
	// rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_SIGNATURE = 16

func makeWorldComponent[T any](w *world) *component[T] {
	component := &component[T]{
		world: w,
	}

	t := reflect.TypeOf(component)
	elem := reflect.TypeOf((*T)(nil)).Elem()

	component.t = t
	component.elem = elem
	component.id = w.registerComponent(component)

	return component
}

func makeWorldSystem[T *ISystem](w *world) *system[T] {
	system := &system[T]{
		world: w,
	}

  // var obj T
  fmt.Println(">>", reflect.Zero(reflect.TypeOf((*T)(nil)).Elem().Elem()))
	// w.registerSystem(obj)
  // t := reflect.TypeOf(&obj)
  // zero := reflect.Zero(t)
  // fmt.Println("t", t, zero, obj)

	// z := reflect.Zero(reflect.TypeOf((*T)(nil))).Interface().(*ISystem)
 //  fmt.Println("checking", reflect.ValueOf((*T)(nil)))
	//
 //  // fmt.Println("system", reflect.Zero() )
	// zero := reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(ISystem)
	// zero2 := reflect.Zero(reflect.TypeOf((*T)(nil))).Interface().(ISystem)
 //  fmt.Println("zero", zero)
 //  fmt.Println("zero2", reflect.TypeOf(zero2))
	// // zero := reflect.Zero(reflect.TypeOf((*T)(nil)).Elem().Elem()).Interface().(ISystem)
	// w.registerSystem(zero2)

	return system
}

func newWorld(entityCapacity int) *world {
	var rb *ringbuffer[Entity] = NewRingBuffer[Entity](entityCapacity)

	for idx := Entity(0); idx < Entity(entityCapacity); idx++ {
		Enqueue(rb, idx)
	}

	return &world{
		// Component_______________________________________________________________
		componentMap:        make(map[int]*componentArray),
		componentToId:       make(map[reflect.Type]int),
		elemToComponent:     make(map[reflect.Type]reflect.Type),
		// Entity__________________________________________________________________
		availEntities:       rb,
		entityToSignature:   make(map[Entity]*Signature),
		// System__________________________________________________________________
		systems:             make([]ISystem, 0, 0),
		systemElemToIdx:     make(map[reflect.Type]int),
    systemIdToSignature: make(map[int]*Signature),
	}
}

func newComponentArray(t reflect.Type) *componentArray {
	return &componentArray{
		data: reflect.MakeSlice(reflect.SliceOf(t), 0, 0),
	}
}

type Entity = uint32

type ISystem interface {
	Update(w *world)
  Id() int
  AddEntity(entity Entity)
}

type system[T any] struct {
	world *world
  id    int
}

func (s *system[T]) Id() int {
  return s.id
}

type IComponent interface {
	Id() int
	Type() reflect.Type
	ElemType() reflect.Type
}

type component[T any] struct {
	world *world
	id    int
	t     reflect.Type
	elem  reflect.Type
}

func (c component[T]) Id() int {
	return c.id
}

func (c component[T]) Type() reflect.Type {
	return c.t
}

func (c component[T]) ElemType() reflect.Type {
	// elem := reflect.TypeOf((*T)(nil)).Elem()
	return c.elem
}

func (c *component[T]) set(entity Entity, val *T) {}
func (c *component[T]) get(entity Entity)         {}

type componentArray struct {
	data         reflect.Value
	entityToData map[Entity]int
	dataToEntity map[int]Entity
}

type world struct {
	// Components________________________________________________________________
	componentMap    map[int]*componentArray // CHANGE TO ARRAY
	componentToId   map[reflect.Type]int
	elemToComponent map[reflect.Type]reflect.Type

	// Entity____________________________________________________________________
	availEntities     *ringbuffer[Entity]
	entityToSignature map[Entity]*Signature

	// System____________________________________________________________________
	systems             []ISystem
	systemElemToIdx     map[reflect.Type]int
	systemIdToSignature map[int]*Signature
}

func UNUSED(a ...interface{}) {}

func (w *world) registerSystem(system ISystem) int {
	sType := reflect.TypeOf(system)
	idx := len(w.systems)

	w.systems = append(w.systems, system)
	w.systemElemToIdx[sType] = idx

  fmt.Println("adding system", sType, "system:", (w.systems[0]))

	return idx
}

func (w *world) setSystemSignature(id int, signature *Signature) {
  w.systemIdToSignature[id] = signature
}

func (w *world) registerComponent(component IComponent) int {
	cId := len(w.componentMap)
	cType := component.Type()
	cElmType := component.ElemType()

	// unnecessary?
	w.componentMap[cId] = newComponentArray(cElmType)
	w.componentToId[cType] = cId
	w.elemToComponent[cElmType] = cType

	return cId
}

func (w *world) getComponentId(component interface{}) int {
	cElemType := reflect.TypeOf(component)

	return w.componentToId[w.elemToComponent[cElemType]]
}

func (w *world) newEntity(components ...interface{}) Entity {
	entity, _ := Dequeue(w.availEntities)

	eSignature := NewSignature(MAX_SIGNATURE)

	for _, component := range components {
		cType := reflect.TypeOf(component)
		cId := w.getComponentId(component)

		eSignature.Set(cId)

		fmt.Println("adding component", cType, "at", cId)
	}

	w.entityToSignature[entity] = eSignature

	fmt.Println("entity signature", eSignature.String())

  // add entity to system 
  for _, system := range w.systems {
    sType := reflect.TypeOf(system)
    id := w.systemElemToIdx[sType]
    fmt.Println("checking system", sType)
    // sId := (*system).Id()
    //
    sSignature := w.systemIdToSignature[id]

    if (eSignature.Int() & sSignature.Int()) == sSignature.Int() {
      fmt.Println("adding entity", entity)
      system.AddEntity(entity)
    }
  }

	return entity
}

func (w *world) run() {
  for _, system := range w.systems {
    system.Update(w)
  }
}

type Pos struct {
	X int
	Y int
}

type PlayerSystem struct {
  id       int
	entities []Entity
}

func (s PlayerSystem) Update(w *world) {
	fmt.Println("player system update")
  for _, entity := range s.entities {
    fmt.Println("doing something with entity", entity)
  }
}

func (s PlayerSystem) Id() int {
  return s.id
}

func (s *PlayerSystem) AddEntity(entity Entity) {
  if (!slices.Contains(s.entities, entity)) {
    s.entities = append(s.entities, entity)
  }
}

func main() {
	world := newWorld(10)
	Position := makeWorldComponent[Pos](world)
	fmt.Println(Position)

  playerSystem := PlayerSystem{}
  id := world.registerSystem(&playerSystem)

	// playerSystem := makeWorldSystem[PlayerSystem](world)
  {
    signature := NewSignature(MAX_SIGNATURE)
    signature.Set(Position.Id())
    world.setSystemSignature(id, signature)
  }
	//
	//
	world.newEntity(Pos{})
  world.run()
  fmt.Println(playerSystem)

	UNUSED(fmt.Println, Position, playerSystem)
}
