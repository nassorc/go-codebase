package scenes

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	engine "github.com/nassorc/gandalf"
	ecs "github.com/nassorc/gandalf/ecs"
	c "github.com/nassorc/gandalf/examples/platformer/components"
	"github.com/nassorc/gandalf/examples/platformer/systems"
)

type PlayerScene struct {
	world *ecs.World
}

func (s *PlayerScene) Setup(e *engine.Engine) {
	world := ecs.NewWorld()
	s.world = world

	world.RegisterAction(rl.KeyW, "Up")
	world.RegisterAction(rl.KeyS, "Down")
	world.RegisterAction(rl.KeyA, "Left")
	world.RegisterAction(rl.KeyD, "Right")

	world.RegisterComponents(
		&c.Tag{},
		&c.Transform{},
		&c.Size{},
		&c.Color{},
		&c.Input{},
		&c.RigidBody{},
		&c.Movable{},
	)

	world.RegisterSystem(systems.InputSystem, &c.Input{})
	world.RegisterSystem(systems.PhysicsSystem, &c.RigidBody{}, &c.Transform{})

	world.CreateEntity(
		&c.Movable{},
		&c.Transform{Pos: rl.NewVector2(0, 0)},
		&c.Tag{Name: "player"},
		&c.Size{Width: 32, Height: 32},
		&c.Color{C: rl.Green},
		&c.Input{},
		&c.RigidBody{
			Size: rl.NewVector2(32, 32),
		},
	)

	world.CreateEntity(
		&c.Transform{Pos: rl.NewVector2(60, 60), PrevPos: rl.NewVector2(60, 60)},
		&c.Tag{Name: "tile"},
		&c.Size{Width: 32, Height: 32},
		&c.Color{C: rl.Blue},
		&c.RigidBody{
			Size: rl.NewVector2(32, 32),
		},
	)

	world.CreateEntity(
		&c.Transform{Pos: rl.NewVector2(60+32, 60), PrevPos: rl.NewVector2(60+32, 60)},
		&c.Tag{Name: "tile"},
		&c.Size{Width: 32, Height: 32},
		&c.Color{C: rl.Blue},
		&c.RigidBody{
			Size: rl.NewVector2(32, 32),
		},
	)

	// witch := rl.LoadTexture("./resources/Blue_witch/B_witch_idle.png")

	// var (
	// 	frameCount    int     = 0
	// 	frames        float32 = 6
	// 	currentFrame  float32 = 0
	// 	textureWidth  float32 = 32
	// 	textureHeight float32 = 48
	// )

	// LevelLoader := g.NewLevelLoader(world)
	// LevelLoader.Load("./level1.txt")

}

func (s *PlayerScene) Update() {
	s.world.Run()
}

func (s *PlayerScene) Render() {
	for _, entity := range s.world.Entities {
		var size *c.Size
		var transform *c.Transform
		var color *c.Color
		var rigidBody *c.RigidBody

		entity.GetData(&size, &transform, &color, &rigidBody)

		var pos = transform.Pos

		rl.DrawRectangle(int32(pos.X), int32(pos.Y), int32(size.Width), int32(size.Height), color.C)
		rl.DrawCircle(int32(pos.X), int32(pos.Y), 2, rl.Red)

		if rigidBody != nil {
			rl.DrawRectangleLines(int32(pos.X), int32(pos.Y), int32(size.Width), int32(size.Height), rl.Red)
		}
	}
}
