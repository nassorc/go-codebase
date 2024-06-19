package main

import (
	engine "github.com/nassorc/gandalf"
	config "github.com/nassorc/gandalf/configParser/utils"
	"github.com/nassorc/gandalf/examples/platformer/scenes"
)

func main() {
	engine := engine.NewEngine(&scenes.GameScene{}, "./config.json", config.NewJsonConfigParser())

	engine.Run()
}
