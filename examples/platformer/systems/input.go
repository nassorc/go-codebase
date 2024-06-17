package systems

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/nassorc/gandalf/ecs"
	c "github.com/nassorc/gandalf/examples/platformer/components"
)

func InputSystem(world *ecs.World, entities []*ecs.Entity) {
	for _, entity := range entities {
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
		if rl.IsKeyDown(rl.KeyRight) {
			// world
		}

		transform.PrevPos = transform.Pos
		transform.Pos.X += mx
		transform.Pos.Y += my
		fmt.Println(transform.Pos.X)
	}

}
