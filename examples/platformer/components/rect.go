package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
)

type Rect struct {
	Size  rl.Vector2
	Color rl.Color
}

var RectID = gandalf.CreateComponentID[Rect]()
