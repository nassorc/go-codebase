package scenes

import (
	// rl "github.com/gen2brain/raylib-go/raylib"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	engine "github.com/nassorc/gandalf"
	ecs "github.com/nassorc/gandalf/ecs"
	// c "github.com/nassorc/gandalf/examples/platformer/components"
	// "github.com/nassorc/gandalf/examples/platformer/systems"
)

type Scene2 struct {
	engine *engine.Engine
	world  *ecs.World
}

func (s *Scene2) Setup(e *engine.Engine) {
	s.engine = e
}
func (s *Scene2) Update() {
	if rl.IsKeyDown(rl.KeyD) {
		fmt.Println("KEY D PRESSED")
		s.engine.ChangeScene(&Scene1{})
	}
}
func (s *Scene2) Render() {
	rl.DrawText("Scene 2. AHHHHHHHHHHHHHHHHHHH", 50, 50, 28, rl.Lime)
}
