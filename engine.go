package gandalf

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene interface {
	Setup(*World)
}

type GameHandle struct {
	game  Scene
	world *World
}

func (g *GameHandle) Update() {
	g.world.Tick()
}

func NewEngine(game Scene, size int) *Engine {
	var assetMgr = NewAssetManager()
	var engine = &Engine{
		size:     size,
		assetMgr: assetMgr,
	}
	var world = CreateWorld(size, engine, assetMgr)

	game.Setup(world)

	gameHandle := &GameHandle{
		game,
		world,
	}

	engine.game = gameHandle

	return engine
}

type Engine struct {
	size                int
	game                *GameHandle
	assetMgr            *AssetManager
	nextGame            Scene
	isChangeGamePending bool
}

func (e *Engine) Run() {
	e.Init()
	defer e.Close()

	for !rl.WindowShouldClose() {
		{
			rl.ClearBackground(rl.Black)

			rl.BeginDrawing()

			if e.isChangeGamePending {
				// fmt.Println("SWITCHING SCENES")

				world := CreateWorld(e.size, e, e.assetMgr)

				gameHandle := &GameHandle{
					game:  e.nextGame,
					world: world,
				}

				gameHandle.game.Setup(world)
				e.game = gameHandle

				e.nextGame = nil
				e.isChangeGamePending = false
			}

			e.game.Update()

			rl.EndDrawing()

		}
	}
}

func (e *Engine) Init() {
	rl.InitWindow(800, 576, "Scene Title")
	rl.SetTargetFPS(60)
}

func (e *Engine) Close() {
	rl.CloseWindow()
}

func (e *Engine) ChangeScene(newGame Scene) {
	e.nextGame = newGame
	e.isChangeGamePending = true
}
