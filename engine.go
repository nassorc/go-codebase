package gandalf

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game interface {
	Setup(*World)
}

type GameHandle struct {
	game  Game
	world *World
}

func (g *GameHandle) Update() {
	g.world.update()
}

func NewEngine(game Game) *Engine {
	engine := &Engine{}

	world := createWorld(engine)
	game.Setup(world)

	gameHandle := &GameHandle{
		game,
		world,
	}

	engine.game = gameHandle

	return engine
}

type Engine struct {
	game                *GameHandle
	nextGame            Game
	isChangeGamePending bool
}

func (e *Engine) Run() {
	e.init()
	defer e.close()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.Black)

			if e.isChangeGamePending {
				// fmt.Println("SWITCHING SCENES")

				world := createWorld(e)

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

func (e *Engine) changeGame(newGame Game) {
	e.nextGame = newGame
	e.isChangeGamePending = true
}
