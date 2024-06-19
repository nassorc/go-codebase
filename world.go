package gandalf

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func CreateWorld() *World {
	return &World{}
}

type World struct {
	text string
}

func (w *World) SetText(text string) {
	w.text = text
}

func (w *World) Update() {
	rl.DrawText(w.text, 100, 100, 32, rl.Purple)
}
