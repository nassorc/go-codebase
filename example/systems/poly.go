package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
	"github.com/nassorc/gandalf/cmd/components"
)

func RotatePoly(w *gandalf.World, entities []*gandalf.EntityHandle) {
	for _, entity := range entities {
		var poly *components.Poly
		entity.Unpack(&poly)

		poly.Rotation += 6
	}
}

func DrawPoly(w *gandalf.World, entities []*gandalf.EntityHandle) {
	for _, entity := range entities {
		var poly *components.Poly
		var transform *components.Transform
		entity.Unpack(&poly, &transform)

		rl.DrawPoly(*transform.Pos, poly.Sides, poly.Radius, poly.Rotation, poly.Color)
	}
}
