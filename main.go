package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_SIGNATURE_SIZE = 16

type WindowConfig struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Title  string `json:"title"`
}

type TextureInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type AnimationInfo struct {
	Name        string `json:"name"`
	TextureName string `json:"texture_name"`
	Frames      int    `json:"frames"`
	Rows        int    `json:"rows"`
	Cols        int    `json:"cols"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
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

func (e *Engine) Run(game IGame) {
	e.init()
	defer e.close()

	game.Setup()

	for !rl.WindowShouldClose() {

		game.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		game.Render()

		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
		rl.DrawText("hello", 32, 32, 20, rl.Black)

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

func (e *Engine) close() {
	rl.CloseWindow()
}

type IGame interface {
	Setup()
	Update()
	Render()
}

type IScene interface {
	Setup()
	Update()
	Render()
}

type Game struct {
	scene IScene
}

func (g *Game) Update() {
	g.scene.Update()
}

func (g *Game) Render() {
	g.scene.Render()
}

func (g *Game) Setup() {
	g.scene.Setup()
}

func NewLevelLoader(w *World) *LevelLoader {
	return &LevelLoader{
		w,
	}
}

type LevelLoader struct {
	w *World
}

func (l *LevelLoader) load(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		panic("could not load path")
	}
	for _, line := range bytes.Split(content, []byte("\n")) {
		fmt.Println("line", string(line))
		args := bytes.Split(line, []byte(" "))

		if string(args[0]) == "Box" {
			x, err := strconv.Atoi(string(args[1]))
			if err != nil {
				panic("could not convert point to int")
			}
			y, err := strconv.Atoi(string(args[2]))
			if err != nil {
				panic("could not convert point to int")
			}

			pos := rl.NewVector2(float32(100+(x*64)), float32(100+(y*64)))

			l.w.NewEntity(
				&Transform{pos: pos, prevPos: pos},
				&Tag{"tile"},
				&Size{64, 64},
				&Color{rl.Green},
				&RigidBody{
					size: rl.NewVector2(64, 64),
				},
			)
		}
	}
}

func main() {
	engine := NewEngine("./config.json", NewJsonConfigParser())

	engine.Run(&Game{
		scene: &PlayerScene{},
	})
}

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

type ActionType = int

const (
	Start ActionType = iota
	End
)

// func NewAction(k int, t ActionType) Action {
// 	return Action {
// 		Type: t,
// 		Key: k,
// 	}
// }

// type Action struct {
// 	Type ActionType
// 	Key int
// }

func NewWorld() *World {
	rb := NewRingBuffer[int](100)

	for idx := 0; idx < 100; idx++ {
		Enqueue(rb, idx)
	}

	return &World{
		// ActionMap:            make(map[int]string),
		typeToComponent:      make(map[reflect.Type]int),
		availIds:             rb,
		systemIdxToSignature: make(map[int]*Signature),
	}
}

type World struct {
	// ActionMap            map[int]string
	Components           []*ComponentArray
	typeToComponent      map[reflect.Type]int
	availIds             *ringbuffer[int]
	Entities             []*Entity
	Systems              []*ISystem
	systemIdxToSignature map[int]*Signature
}

func (w *World) Run() {
	for _, system := range w.Systems {
		(*system).Update(w)
	}
}

func (w *World) RegisterAction(key int, name string) {
	// w.ActionMap[key] = name
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
