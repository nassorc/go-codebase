package gandalf

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene interface {
	Setup(*World)
}

type SceneNode struct {
	name  string
	scene Scene
	world *World
	next  *SceneNode
	prev  *SceneNode
}

// func (g *SceneNode) Update() {
// 	g.world.Tick()
// }

func NewEngine(scene Scene, size int) *Engine {
	var assetMgr = NewAssetManager()
	var engine = &Engine{
		size:     size,
		assetMgr: assetMgr,
	}
	var world = CreateWorld(size, engine, assetMgr)

	scene.Setup(world)

	sceneNode := &SceneNode{
		scene: scene,
		world: world,
	}

	engine.curScene = sceneNode

	return engine
}

type Engine struct {
	size                 int
	curScene             *SceneNode
	assetMgr             *AssetManager
	newScene             Scene
	isChangeScenePending bool
	popScene             bool
	pushScene            bool
}

func (e *Engine) Run() {
	e.Init()
	defer e.Close()

	for !rl.WindowShouldClose() {
		{
			rl.ClearBackground(rl.White)

			rl.BeginDrawing()

			// push scene
			if e.pushScene {
				world := CreateWorld(e.size, e, e.assetMgr)

				c := e.curScene

				for c.next != nil {
					c = c.next
				}

				// e.curScene = e.curScene.next
				// ! REVERSE
				c.next = &SceneNode{
					scene: e.newScene,
					prev:  c,
				}

				c.next.scene.Setup(world)
				c.next.world = world

				e.newScene = nil
				e.pushScene = false

			} else if e.isChangeScenePending { //  new scene
				world := CreateWorld(e.size, e, e.assetMgr)

				e.newScene.Setup(world)

				e.curScene = &SceneNode{
					scene: e.newScene,
					world: world,
				}

				// bookkeeping
				e.isChangeScenePending = false
				e.newScene = nil

			} else if e.popScene {
				out := e.curScene
				e.curScene = e.curScene.prev
				e.curScene.next = nil
				out.prev = nil
				out.next = nil
				out.world = nil
				out.scene = nil

				e.popScene = false
			}

			// e.scene.Update()
			e.updateScenes()

			rl.EndDrawing()

		}
	}
}

func (e *Engine) updateScenes() {
	// e.curScene.world.Tick()
	cur := e.curScene
	for cur != nil {
		cur.world.Tick()
		cur = cur.next
	}
}

func (e *Engine) Init() {
	rl.InitWindow(800, 576, "Scene Title")
	rl.SetTargetFPS(60)
}

func (e *Engine) Close() {
	rl.CloseWindow()
}

func (e *Engine) PushScene(newScene Scene) {
	e.newScene = newScene
	e.pushScene = true
}

func (e *Engine) PopScene() {
	e.popScene = true
}

func (e *Engine) ChangeScene(newScene Scene) {
	e.newScene = newScene
	e.isChangeScenePending = true
}

func (e *Engine) IsCurrentScene(name string) bool {
	return e.curScene.name == name
}
