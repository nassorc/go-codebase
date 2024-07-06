package main

import (
	"fmt"
	"image/color"

	gandalf "github.com/nassorc/gandalf"
	scenes "github.com/nassorc/gandalf/cmd/scenes"
)

func main() {
	var engine = gandalf.NewEngine(scenes.NewMenuScene(color.RGBA{255, 255, 255, 255}, color.RGBA{210, 20, 90, 255}), 1000)

	engine.Run()
	fmt.Println("working")
}
