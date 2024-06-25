package components

import rl "github.com/gen2brain/raylib-go/raylib"

func NewPoly(sides int32, radius float32, rotation float32, color rl.Color) *Poly {
	return &Poly{
		sides,
		radius,
		rotation,
		color,
	}
}

type Poly struct {
	Sides    int32
	Radius   float32
	Rotation float32
	Color    rl.Color
}
