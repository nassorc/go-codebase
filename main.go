package main

import (
	"fmt"
	"image/color"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
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


func UpdatePlayer(player *Entity, w *World) {
  var transform *Transform
  player.GetData(&transform)

  const speed = 2
  var mx float32 = 0
  var my float32 = 0

  if rl.IsKeyDown(rl.KeyW) {
    my = -speed
  }
  if rl.IsKeyDown(rl.KeyS) {
    my = speed
  }
  if rl.IsKeyDown(rl.KeyA) {
    mx = -speed
  }
  if rl.IsKeyDown(rl.KeyD) {
    mx = speed
  }

  transform.pos.X += mx
  transform.pos.Y += my
}

func main() {
  fmt.Println()
	rl.InitWindow(960, 640, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

  world := NewWorld()
  world.RegisterComponents(&Tag{}, &Transform{}, &Size{}, &Color{}, &Input{})

  world.RegisterSystem(&InputSystem{})

  player := world.NewEntity(&Transform{ pos: rl.NewVector2(0, 0) }, &Tag{ "player" }, &Size{32, 32}, &Color{rl.Green}, &Input{})

  // fmt.Println("system", *(world.systems[0]))

  var playerTag *Tag
  var playerPos *Transform
  player.GetData(&playerTag, &playerPos)

	for !rl.WindowShouldClose() {
    // Update Player___________________________________________________________
    // UpdatePlayer(player, world)
    // Update__________________________________________________________________

    world.Run()

		rl.BeginDrawing()
    rl.ClearBackground(color.RGBA{R: 24, G: 24, B: 24})

    for _, entity := range world.entities {
      var size *Size
      var transform *Transform
      var color *Color
      entity.GetData(&size, &transform, &color)

      var pos = transform.pos

      rl.DrawRectangle(int32(pos.X), int32(pos.Y), int32(size.Width), int32(size.Height), color.c)
    }

    rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

    rl.DrawText("hello", 32, 32, 20, rl.Black)

		rl.EndDrawing()
	}
}

type ISystem interface {
  Update(w *World)
  AddEntity(e *Entity)
}

type PlayerSystem struct {
}

func (s *PlayerSystem) Update(w *World) {
  // assumes that it will receive the player entity as the first element
  // alternative solution: 
  // 1. the AddEntity method doesn't require the system to
  // have a list of entities. Update it so that it only deals with a single
  // entity.
  // 2. the world keeps a reference to the player, since the player is a special
  // type of entity that would likely be used in many systems.
  // 3. a query system to fetch for specific entities. 
}

type InputSystem struct {
  entities []*Entity
}

func (s *InputSystem) Update(w *World) {
  for _, entity := range s.entities {
    var transform *Transform
    entity.GetData(&transform)

    const speed = 2
    var mx float32 = 0
    var my float32 = 0

    if rl.IsKeyDown(rl.KeyW) {
      my = -speed
    }
    if rl.IsKeyDown(rl.KeyS) {
      my = speed
    }
    if rl.IsKeyDown(rl.KeyA) {
      mx = -speed
    }
    if rl.IsKeyDown(rl.KeyD) {
      mx = speed
    }

    transform.pos.X += mx
    transform.pos.Y += my
  }

  player := s.entities[0]
  var transform *Transform
  player.GetData(&transform)

  const speed = 2
  var mx float32 = 0
  var my float32 = 0

  if rl.IsKeyDown(rl.KeyW) {
    my = -speed
  }
  if rl.IsKeyDown(rl.KeyS) {
    my = speed
  }
  if rl.IsKeyDown(rl.KeyA) {
    mx = -speed
  }
  if rl.IsKeyDown(rl.KeyD) {
    mx = speed
  }

  transform.pos.X += mx
  transform.pos.Y += my
}

func (s *InputSystem) AddEntity(e *Entity) {
  s.entities = append(s.entities, e)
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
  entities             []*Entity

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

  w.entities = append(w.entities, entity)

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

type Transform struct {
  pos rl.Vector2
}

type Tag struct {
  name string
}

type Size struct {
  Width int
  Height int
}

type Color struct {
  c rl.Color
}

type Input struct {
  Up      bool
  Down    bool
  Left    bool
  Right   bool
  CanJump bool
}

func UNUSED(x ...interface{}) {}
