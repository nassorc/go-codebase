package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	g "github.com/nassorc/gandalf"
	c "github.com/nassorc/gandalf/examples/platformer/components"
)

type InputSystem struct {
	entities []*g.Entity
}

func (s *InputSystem) Update(w *g.World) {
	for _, entity := range s.entities {
		var transform *c.Transform
		entity.GetData(&transform)

		const speed = 2
		var mx float32 = 0
		var my float32 = 0

		if rl.IsKeyDown(rl.KeyW) {
			my = -speed
		}
		if rl.IsKeyDown(rl.KeyS) {
			my = speed
		}
		if rl.IsKeyDown(rl.KeyA) {
			mx = -speed
		}
		if rl.IsKeyDown(rl.KeyD) {
			mx = speed
		}

		transform.PrevPos = transform.Pos
		transform.Pos.X += mx
		transform.Pos.Y += my
	}
}

func (s *InputSystem) AddEntity(e *g.Entity) {
	s.entities = append(s.entities, e)
}
