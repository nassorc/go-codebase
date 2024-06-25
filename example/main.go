package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	gandalf "github.com/nassorc/gandalf"
	scenes "github.com/nassorc/gandalf/cmd/scenes"
)

func main() {
	var engine = gandalf.NewEngine(scenes.NewMenuScene(scenes.MenuItems, scenes.MenuActions, rl.White, rl.Red), 100000)

	engine.Run()
	fmt.Println("working")
}
