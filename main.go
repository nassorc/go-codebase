package main

import (
	"fmt"
	"image/color"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var MAX_SIGNATURE_SIZE = 16

func UNUSED(name ...interface{}) {}

type CTest struct {
  val int
}

type CComp struct {
  val int
}

type PlayerSystem struct {
  entities []Entity
}

func (s *PlayerSystem) update(w *World) {}
func (s *PlayerSystem) render(w *World) {}
func (s *PlayerSystem) addEntity(entity Entity) {
  s.entities = append(s.entities, entity)
}

type IScene interface{
  setup()
}

type ComponentArr struct {
  id           int            // position in the world.components array
  data         reflect.Value
  entityToData map[Entity]int
}

func (c *ComponentArr) Set(entity Entity, value reflect.Value) {
  var idx = c.data.Len()
  c.data = reflect.Append(c.data, value)
  c.entityToData[entity] = idx
}

func (c *ComponentArr) Get(entity Entity) reflect.Value {
  var idx = c.entityToData[entity]
  return c.data.Index(idx)
}

type Entity = int

func NewWorld() *World {
  availEntities := NewRingBuffer[Entity](10)

  for idx := 0; idx < 10; idx++ {
    Enqueue(availEntities, idx)
  }

  return &World{
    componentToId: make(map[reflect.Type]int),
    availEntities: availEntities,
    entityToSignature: make(map[Entity]*Signature),
  }
}

type World struct {

  systems           []ISystem
  systemSignatures  []*Signature 

  components        []*ComponentArr
  componentToId     map[reflect.Type]int  // maps component element type to its position in the world.components array

  availEntities     *ringbuffer[Entity]
  entityToSignature map[Entity]*Signature
}

type ISystem interface {
  update(*World)
  render(*World)
  addEntity(Entity)
}

func (w *World) registerSystem(system ISystem, signature *Signature) {
  w.systems = append(w.systems, system)
  w.systemSignatures = append(w.systemSignatures, signature)
}

func (w *World) registerComponent(components ...interface{}) {
  for _, component := range components {
    var cType = reflect.TypeOf(component)
    var id = len(w.componentToId)
    var componentArr = ComponentArr {
      id: id,
      data: reflect.MakeSlice(reflect.SliceOf(cType), 0, 0),
      entityToData: make(map[Entity]int),
    }

    w.components = append(w.components, &componentArr) 
    w.componentToId[cType] = id
  }
}

func (w *World) getComponent(entity Entity, component interface{}) {
  // fmt.Println(world.components[0].data.Index(0))
  t := reflect.TypeOf(component).Elem()
  value := reflect.ValueOf(component)
  fmt.Println("getting component", t.Elem())
  var componentArr = w.components[w.componentToId[t]]
  var idx = componentArr.entityToData[entity]
  // value.Set(componentArr.data.Index(idx))
  value.Elem().Set(componentArr.data.Index(idx))
  // fmt.Println(value, value)
}

func (w *World) getComponentId(component interface{}) int {
    var cType = reflect.TypeOf(component)
    id, ok := w.componentToId[cType]

    if !ok {
      panic(fmt.Errorf("getComponentId: %v is not a component array", cType))
    }

    return id
}

func (w *World) newEntity(components ...interface{}) Entity {
  entity, err := Dequeue(w.availEntities)
  if err != nil {
    panic(err)
  }

  var eSignature = NewSignature(MAX_SIGNATURE_SIZE)

  for _, component := range components {
    var id = w.getComponentId(component)
    var cValue = reflect.ValueOf(component)

    w.components[id].Set(entity, cValue)
    eSignature.Set(id)
  }

  w.entityToSignature[entity] = eSignature

  // inform systems of new entity
  for idx, system := range w.systems {
    var sSignature = w.systemSignatures[idx]

    if (eSignature.Int() & sSignature.Int()) == sSignature.Int() {
      system.addEntity(entity)
    }
  }

  return entity
}

func (w *World) update() {
  for _, system := range w.systems {
    system.update(w)
  }
}

func (w *World) render() {
  for _, system := range w.systems {
    system.render(w)
  }
}

type Game struct {}

func (g *Game) setup() {
  world := NewWorld()
  world.registerComponent(CTest{})
}

type Engine struct{} 

func (e *Engine) run(scene IScene) {
	rl.InitWindow(960, 640, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

    // Update__________________________________________________________________
		rl.BeginDrawing()
    rl.ClearBackground(color.RGBA{R: 24, G: 24, B: 24})

    rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

    rl.DrawText("hello", 32, 32, 20, rl.Black)

		rl.EndDrawing()
	}
}

func main() {
  // var engine = Engine{}
  // engine.run(&Game{})

  var world = NewWorld()
  world.registerComponent(&CTest{})
  world.registerComponent(&CComp{})


  var signature = NewSignature(MAX_SIGNATURE_SIZE)
  signature.Set(world.getComponentId(&CComp{}))
  world.registerSystem(&PlayerSystem{}, signature)

  world.newEntity(&CTest{val: 100})

  world.newEntity(&CComp{val: 100})
  player := world.newEntity(&CComp{val: 100}, &CTest{val: 5})

  var playerTest *CTest

  world.getComponent(player, &playerTest)
  fmt.Println("did i get it?", playerTest)
  // fmt.Println(world.components[0].data.Index(0))
}

// ctest := newWorldComponent[CTest](world)
// ctest.get()
// ctest.set()
// ctest.all()
// player.get(CTest{})

// Issues
// current way of handling components is error prone. For example, we must pass
// a reference to a component to every world function that handles components.
// forgetting to add '&' will treat the component type as a different componentArray
// Possible Solutions:
// 1. Do it once. Users can create a type that will deal with correctly passing
//    in the component structs.
//    Position := newWorldType[Pos](world)
//    Position.get(entity)
//    Position.set(entity, Pos{10, 20})
// 2. Make it so that users will receive an error, or the program panicking if
//    the user doesn't pass in a reference to a component type.

// a component should have its own store
// a comonent should have an id
// a user should be able to query the components id
// a user should be able to create an entity
