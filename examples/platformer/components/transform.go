package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
)

func NewTranform(x float32, y float32) *Transform {
	pos := rl.NewVector2(x, y)
	prev := rl.NewVector2(x, y)
	return &Transform{
		Pos:     pos,
		PrevPos: prev,
	}
}

type Transform struct {
	Pos     rl.Vector2
	PrevPos rl.Vector2
}

var TransformID = gandalf.CreateComponentID[Transform]()
