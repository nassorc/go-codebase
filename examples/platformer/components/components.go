package components

import rl "github.com/gen2brain/raylib-go/raylib"

type RigidBody struct {
	Size   rl.Vector2
	Offset rl.Vector2
}

type Transform struct {
	Vel     rl.Vector2
	Pos     rl.Vector2
	PrevPos rl.Vector2
}

type Movable struct{}

type Tag struct {
	Name string
}

type Size struct {
	Width  int
	Height int
}

type Color struct {
	C rl.Color
}

type Input struct {
	Up      bool
	Down    bool
	Left    bool
	Right   bool
	CanJump bool
}
