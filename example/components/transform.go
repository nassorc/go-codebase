package components

import rl "github.com/gen2brain/raylib-go/raylib"

func NewTranform(x float32, y float32) *Transform {
	pos := rl.NewVector2(x, y)
	return &Transform{
		Pos: &pos,
	}
}

type Transform struct {
	Pos *rl.Vector2
}
