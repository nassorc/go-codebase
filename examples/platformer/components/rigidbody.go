package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
)

func NewRigidBody(size rl.Vector2, offset rl.Vector2) *RigidBody {
	return &RigidBody{
		Size:     size,
		HalfSize: rl.Vector2Scale(size, 0.5),
		Offset:   offset,
	}
}

type RigidBody struct {
	Size     rl.Vector2
	HalfSize rl.Vector2
	Offset   rl.Vector2
}

var RigidBodyID = gandalf.CreateComponentID[RigidBody]()
