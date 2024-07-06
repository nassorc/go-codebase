package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
	"github.com/nassorc/gandalf/cmd/components"
)

func DrawEntities(entities []gandalf.EntityHandle) {
	for _, entity := range entities {
		var transform *components.Transform
		var rect *components.Rect

		entity.Unpack(&transform, &rect)

		rl.DrawRectangle(int32(transform.Pos.X), int32(transform.Pos.Y), int32(rect.Size.X), int32(rect.Size.Y), rect.Color)
	}
}
