package main

import (
	"fmt"

	g "github.com/nassorc/gandalf"
	config "github.com/nassorc/gandalf/configParser/utils"
	"github.com/nassorc/gandalf/examples/platformer/scenes"
)

func main() {
	fmt.Println("platformer")
	engine := g.NewEngine(&scenes.PlayerScene{}, "./config.json", config.NewJsonConfigParser())

	// engine.Run(&g.Game{
	// 	Scene: &scenes.PlayerScene{},
	// })
	engine.Run()
}
