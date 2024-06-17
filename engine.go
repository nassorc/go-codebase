package gandalf

import (
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	config "github.com/nassorc/gandalf/configParser"
)

const MAX_SIGNATURE_SIZE = 16

type Scene interface {
	Setup(e *Engine)
	Update()
	Render()
}

func NewEngine(scene Scene, configPath string, configParser config.IConfigParser) *Engine {
	engine := &Engine{
		configPath:   configPath,
		configParser: configParser,
	}

	scene.Setup(engine)

	engine.scene = scene

	return engine
}

type Engine struct {
	scene        Scene
	configPath   string
	configParser config.IConfigParser
}

func (e *Engine) Run() {
	e.init()
	defer e.close()

	// game.Setup()

	for !rl.WindowShouldClose() {

		e.scene.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		e.scene.Render()

		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}

}

func (e *Engine) init() {
	// Load files
	file, err := os.Open(e.configPath)

	if err != nil {
		panic("failed to open configuration path")
	}

	buf := make([]byte, 1024)

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

func (e *Engine) ChangeScene(newScene Scene) {
	newScene.Setup(e)
	e.scene = newScene
}
