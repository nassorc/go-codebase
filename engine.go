package gandalf

import rl "github.com/gen2brain/raylib-go/raylib"

type Game interface {
	Setup(*World)
}

type GameHandle struct {
	game  Game
	world *World
}

func (g *GameHandle) Update() {
	g.world.Update()
}

func NewEngine(game Game) *Engine {
	world := CreateWorld()
	game.Setup(world)

	gameHandle := GameHandle{
		game,
		world,
	}

	return &Engine{
		gameHandle,
	}
}

type Engine struct {
	game GameHandle
}

func (e *Engine) Run() {
	e.init()
	defer e.close()

	for !rl.WindowShouldClose() {
		// e.game.Update()

		rl.BeginDrawing()
		{
			e.game.Update()
			rl.DrawText("hello world", 0, 0, 18, rl.White)
		}
		rl.EndDrawing()
	}
}

func (e *Engine) init() {
	rl.InitWindow(500, 500, "Game Title")
	rl.SetTargetFPS(60)
}

func (e *Engine) close() {
	rl.CloseWindow()
}
