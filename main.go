package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_SIGNATURE_SIZE = 16

type WindowConfig struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Title  string `json:"title"`
}

type TextureInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
}

type AnimationInfo struct {
	Name         string `json:"name"`
	TextureName  string `json:"texture_name"`
  Frames       int    `json:"frames"`
  Rows         int    `json:"rows"`
  Cols         int    `json:"cols"`
  Width        int    `json:"width"`
  Height       int    `json:"height"`
} 

type Config struct {
	Window    WindowConfig    `json:"window"`
  Textures  []TextureInfo   `json:"textures"`
  Animation []AnimationInfo `json:"animations"`
}

type IConfigParser interface {
	ParseConfig(data []byte) (*Config, error)
}

func NewJsonConfigParser() IConfigParser {
	return &JsonConfigParser{}
}

type JsonConfigParser struct{}

func (p JsonConfigParser) ParseConfig(data []byte) (*Config, error) {
	var out = Config{}
	err := json.Unmarshal(data, &out)

	if err != nil {
		return &out, err
	}

	return &out, nil
}

func NewEngine(configPath string, configParser IConfigParser) *Engine {
	return &Engine{
		configPath,
		configParser,
	}
}

type Engine struct {
	configPath   string
	configParser IConfigParser
}

func (e *Engine) Run() {
	e.init()
	defer e.close()

	world := NewWorld()

	// world.registerAction(rl.KeyW, "Up")

	world.RegisterComponents(
		&Tag{},
		&Transform{},
		&Size{},
		&Color{},
		&Input{},
		&RigidBody{},
		&Movable{},
	)

	world.RegisterSystem(&InputSystem{}, &Input{})
	world.RegisterSystem(&PhysicsSystem{}, &RigidBody{}, &Transform{})

	player := world.NewEntity(
		&Movable{},
		&Transform{pos: rl.NewVector2(0, 0)},
		&Tag{"player"},
		&Size{32, 32},
		&Color{rl.Green},
		&Input{},
		&RigidBody{
			size: rl.NewVector2(32, 32),
		},
	)

	world.NewEntity(
		&Transform{pos: rl.NewVector2(60, 60), prevPos: rl.NewVector2(60, 60)},
		&Tag{"tile"},
		&Size{32, 32},
		&Color{rl.Blue},
		&RigidBody{
			size: rl.NewVector2(32, 32),
		},
	)

	world.NewEntity(
		&Transform{pos: rl.NewVector2(60+32, 60), prevPos: rl.NewVector2(60+32, 60)},
		&Tag{"tile"},
		&Size{32, 32},
		&Color{rl.Blue},
		&RigidBody{
			size: rl.NewVector2(32, 32),
		},
	)

	world.NewEntity(
		&Transform{pos: rl.NewVector2(60, 60+32), prevPos: rl.NewVector2(60, 60+32)},
		&Tag{"tile"},
		&Size{32, 32},
		&Color{rl.Blue},
		&RigidBody{
			size: rl.NewVector2(32, 32),
		},
	)

	var xOffset float32 = 128

	world.NewEntity(
		&Transform{pos: rl.NewVector2(60+xOffset, 60), prevPos: rl.NewVector2(60+xOffset, 60)},
		&Tag{"tile"},
		&Size{32, 32},
		&Color{rl.Blue},
		&RigidBody{
			size: rl.NewVector2(32, 32),
		},
	)

	world.NewEntity(
		&Transform{pos: rl.NewVector2(60+32+xOffset, 60), prevPos: rl.NewVector2(60+32+xOffset, 60)},
		&Tag{"tile"},
		&Size{32, 32},
		&Color{rl.Blue},
		&RigidBody{
			size: rl.NewVector2(32, 32),
		},
	)

	world.NewEntity(
		&Transform{pos: rl.NewVector2(60+xOffset, 60+32), prevPos: rl.NewVector2(60+xOffset, 60+32)},
		&Tag{"tile"},
		&Size{32, 32},
		&Color{rl.Blue},
		&RigidBody{
			size: rl.NewVector2(32, 32),
		},
	)

	witch := rl.LoadTexture("./resources/Blue_witch/B_witch_idle.png")

	var (
		frameCount    int     = 0
		frames        float32 = 6
		currentFrame  float32 = 0
		textureWidth  float32 = 32
		textureHeight float32 = 48
	)
	UNUSED(currentFrame, textureWidth, textureHeight, witch)

	for !rl.WindowShouldClose() {
		// Update Player___________________________________________________________
		// UpdatePlayer(player, world)
		// Update__________________________________________________________________
		world.Run()

		// fmt.Println("component", world.Components[2].Data)

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		for _, entity := range world.Entities {
			var size *Size
			var transform *Transform
			var color *Color
			var rigidBody *RigidBody

			entity.GetData(&size, &transform, &color, &rigidBody)

			var pos = transform.pos

			rl.DrawRectangle(int32(pos.X), int32(pos.Y), int32(size.Width), int32(size.Height), color.c)
			rl.DrawCircle(int32(pos.X), int32(pos.Y), 2, rl.Red)

			if rigidBody != nil {
				rl.DrawRectangleLines(int32(pos.X), int32(pos.Y), int32(size.Width), int32(size.Height), rl.Red)
			}
		}

		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
		rl.DrawText("hello", 32, 32, 20, rl.Black)
		// rl.DrawTexture(witch, 32, 100, rl.White)
		var transform *Transform
		player.GetData(&transform)
		rl.DrawTextureRec(witch, rl.NewRectangle(0, currentFrame*textureHeight, textureWidth, textureHeight), rl.NewVector2(transform.pos.X, transform.pos.Y), rl.White)

		fmt.Println(currentFrame)
		currentFrame = float32((frameCount / 5) % int(frames))
		frameCount += 1

		rl.EndDrawing()
	}

}

func (e *Engine) init() {
	// load files
	file, err := os.Open(e.configPath)

	if err != nil {
		panic("failed to open configuration path")
	}

	buf := make([]byte, 1024, 1024)

	n, err := file.Read(buf)

	if err != nil {
		panic("failed to read configuration path")
	}

	config, err := e.configParser.ParseConfig(buf[:n-1])

	if err != nil {
		panic(err)
	}

	rl.InitWindow(int32(config.Window.Width), int32(config.Window.Height), config.Window.Title)
	rl.SetTargetFPS(60)
}

func (e *Engine) loadAssets() {

}

func (e *Engine) close() {
	rl.CloseWindow()
}

func main() {
	engine := NewEngine("./config.json", NewJsonConfigParser())

	engine.Run()
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

type PhysicsSystem struct {
	entities []*Entity
}

type AxisCollision struct {
	X bool
	Y bool
}

func HasAxisCollision(posA rl.Vector2, posB rl.Vector2, sizeA rl.Vector2, sizeB rl.Vector2) (AxisCollision, rl.Vector2) {
	var halfSizeA = rl.Vector2Scale(sizeA, 0.5)
	var halfSizeB = rl.Vector2Scale(sizeB, 0.5)

	var distX float64 = math.Abs(float64(posA.X - posB.X))
	var distY float64 = math.Abs(float64(posA.Y - posB.Y))
	var deltaX = float32(distX) - (halfSizeA.X + halfSizeB.X)
	var deltaY = float32(distY) - (halfSizeA.Y + halfSizeB.Y)

	var hasXAxisCollision = deltaX < 0
	var hasYAxisCollision = deltaY < 0

	return AxisCollision{X: hasXAxisCollision, Y: hasYAxisCollision}, rl.NewVector2(float32(math.Abs(float64(deltaX))), float32(math.Abs(float64(deltaY))))
	// return AxisCollision{hasXAxisCollision},  AxisCollision{hasYAxisCollision}
}

func (s *PhysicsSystem) Update(w *World) {
	// Collision logic won't work with multiple entities since it directly updates
	// the movable entity's position.
	for _, entityA := range s.entities {
		// process movable entity colliding against entity with a rigid body and a transform component
		var movable *Movable
		entityA.GetData(&movable)

		if movable == nil {
			continue
		}

		for _, entityB := range s.entities {
			if entityA.id == entityB.id {
				continue
			}

			var transformA *Transform
			var rigidBodyA *RigidBody
			entityA.GetData(&transformA, &rigidBodyA)

			var transformB *Transform
			var rigidBodyB *RigidBody
			entityB.GetData(&transformB, &rigidBodyB)

			collision, overlap := HasAxisCollision(transformA.pos, transformB.pos, rigidBodyA.size, rigidBodyB.size)
			prevCollision, _ := HasAxisCollision(transformA.prevPos, transformB.prevPos, rigidBodyA.size, rigidBodyB.size)

			hasCollision := collision.X && collision.Y
			// Resolve movable (entityA) collision

			// vertical collision
			if hasCollision && prevCollision.X {
				if transformA.pos.Y < transformB.pos.Y {
					transformA.pos.Y -= overlap.Y
				} else {
					transformA.pos.Y += overlap.Y
				}
				transformA.prevPos.Y = transformA.pos.Y
			} else if hasCollision && prevCollision.Y {
				// horizontal collision
				transformA.prevPos.X = transformA.pos.X
				if transformA.pos.X < transformB.pos.X {
					transformA.pos.X -= overlap.X
				} else {
					transformA.pos.X += overlap.X
				}
			}

			// DEBUG
			var color *Color
			entityB.GetData(&color)
			if hasCollision {
				color.c = rl.Red
			} else {
				color.c = rl.Black
			}

		}
	}
}

func (s *PhysicsSystem) AddEntity(e *Entity) {
	s.entities = append(s.entities, e)
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

		transform.prevPos = transform.pos
		transform.pos.X += mx
		transform.pos.Y += my
	}
}

func (s *InputSystem) AddEntity(e *Entity) {
	s.entities = append(s.entities, e)
}

type RigidBody struct {
	size   rl.Vector2
	offset rl.Vector2
}

type Transform struct {
	vel     rl.Vector2
	pos     rl.Vector2
	prevPos rl.Vector2
}

type Movable struct{}

type Tag struct {
	name string
}

type Size struct {
	Width  int
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

// [1]https://www.youtube.com/playlist?list=PL_xRyXins848nDj2v-TJYahzvs-XW9sVV
// https://rivermanmedia.com/object-oriented-game-programming-the-scene-system/
// - Scene Management (similar to [1])

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

		cIdx, ok := e.world.typeToComponent[t.Elem()]

		if !ok {
			continue
		}

		carr := e.world.Components[cIdx]
		// carr := e.world.GetComponentArray(t.Elem())
		idx, ok := carr.entityToData[e.id]

		if !ok {
			continue
		}
		// ptrVal.Set(arr.Index(0).Addr().Elem())

		// newPtr := reflect.New(carr.data.Index(idx).Elem().Type())
		// newPtr.Elem().Set(carr.data.Index(idx).Elem())
		// val.Set(newPtr)

		val.Set(carr.Data.Index(idx).Addr().Elem())

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
	return &ComponentArray{
		Data:         reflect.MakeSlice(reflect.SliceOf(t), 0, 0),
		entityToData: make(map[int]int),
	}
}

type ComponentArray struct {
	Data         reflect.Value
	entityToData map[int]int
}

func (c *ComponentArray) AppendData(entityId int, value reflect.Value) {
	idx := c.Data.Len()
	c.Data = reflect.Append(c.Data, value)
	c.entityToData[entityId] = idx
}

func (c *ComponentArray) SetData(entityId int, value reflect.Value) {
	idx := c.entityToData[entityId]
	c.Data.Index(idx).Set(value)
}

func (c *ComponentArray) GetData(entityId int) reflect.Value {
	idx := c.entityToData[entityId]

	return c.Data.Index(idx)
}

func (c *ComponentArray) RemoveEntity(entityId int) {
}

func NewWorld() *World {
	rb := NewRingBuffer[int](100)

	for idx := 0; idx < 100; idx++ {
		Enqueue(rb, idx)
	}

	return &World{
		typeToComponent: make(map[reflect.Type]int),
		availIds:        rb,

		systemIdxToSignature: make(map[int]*Signature),
	}
}

type World struct {
	Components      []*ComponentArray
	typeToComponent map[reflect.Type]int

	availIds *ringbuffer[int]
	Entities []*Entity

	Systems              []*ISystem
	systemIdxToSignature map[int]*Signature
}

func (w *World) Run() {
	for _, system := range w.Systems {
		(*system).Update(w)
	}
}

func (w *World) RegisterSystem(system ISystem, components ...interface{}) {
	// t := reflect.TypeOf(system)
	idx := len(w.Systems)
	w.Systems = append(w.Systems, &system)
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
		idx := len(w.Components)
		w.Components = append(w.Components, carr)
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

	fmt.Println("NEW ENTITY", id)

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
		id:        id,
		world:     w,
		signature: eSignature,
	}

	w.Entities = append(w.Entities, entity)

	for idx, signature := range w.Systems {
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
	return w.Components[cidx]
}

func UNUSED(x ...interface{}) {}
