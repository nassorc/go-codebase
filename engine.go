package gandalf

import (
	"fmt"

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
	game *GameHandle
}

func (e *Engine) Run() {
	e.init()
	defer e.close()

	for !rl.WindowShouldClose() {
		// e.game.Update()

		e.game.Update()
		rl.BeginDrawing()
		{
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

func (e *Engine) changeGame(newGame Game) {
	world := createWorld(e)
	newGame.Setup(world)
	fmt.Println("NEW GAME", newGame)
	fmt.Println("NEW world", world)

	fmt.Println("old", e)

	e.game.game = newGame
	e.game.world = world
}
