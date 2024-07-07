package scenes

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
	// "github.com/nassorc/gandalf/cmd/components"
	// "github.com/nassorc/gandalf/cmd/systems"
)

type Inventory struct {
	world  *gandalf.World
	player gandalf.EntityHandle
}

func (scene *Inventory) Setup(world *gandalf.World) {
	world.RegisterSystem2(func(entities []gandalf.EntityHandle) {
		rl.DrawRectangle(100, 100, 600, 200, rl.Black)
	})

}
