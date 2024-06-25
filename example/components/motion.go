package components

import rl "github.com/gen2brain/raylib-go/raylib"

func NewMotion(velocity rl.Vector2, acceleration rl.Vector2) *Motion {
	return &Motion{
		velocity,
		acceleration,
	}
}

type Motion struct {
	Velocity     rl.Vector2
	Acceleration rl.Vector2
}
