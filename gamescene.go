package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerScene struct {
  world *World
}

func (s *PlayerScene) Setup() {
	world := NewWorld()
  s.world = world;

	world.RegisterAction(rl.KeyW, "Up")
	world.RegisterAction(rl.KeyS, "Down")
	world.RegisterAction(rl.KeyA, "Left")
	world.RegisterAction(rl.KeyD, "Right")

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

	world.NewEntity(
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
	UNUSED(currentFrame, textureWidth, textureHeight, witch, frameCount, frames)
}

func (s *PlayerScene) Update() {
  s.world.Run()
}

func (s *PlayerScene) Render() {
	for _, entity := range s.world.Entities {
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
}
