package main

import (
	"fmt"

	gandalf "github.com/nassorc/gandalf"
)

// type Scene interface {
// 	Setup()
// }

// type MenuScene struct {

// }

//	func (ms *MenuScene) Setup() {
//		world := gandalf.CreateWorld()
//	}
// func InputSystem(scene *gandalf.Scene, entities []*gandalf.Entity) {
// 	// scene.changescene(Game{})
// }

type Game struct {
	// currentScene Scene
	world *gandalf.World
}

func (g *Game) Setup(world *gandalf.World) {
	g.world = world

	world.SetText("BLAH BLAH")

	fmt.Print(world)
	// g.currentScene.Setup()
}

func main() {
	var engine = gandalf.NewEngine(&Game{
		// sceneManager: &SceneManager{}
	})

	// scenenManager.onChangeRequest()

	engine.Run()
	fmt.Println("working")
}
