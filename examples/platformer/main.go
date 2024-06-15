package main

import (
	"fmt"

	g "github.com/nassorc/gandalf"
	config "github.com/nassorc/gandalf/configParser/utils"
	"github.com/nassorc/gandalf/examples/platformer/scenes"
)

// type Game struct{}

// func (*Game) Setup() {

// }

func main() {
	fmt.Println("platformer")
	engine := g.NewEngineWithConfig("./config.json", config.NewJsonConfigParser())

	engine.Run(&g.Game{
		Scene: &scenes.PlayerScene{},
	})
}
