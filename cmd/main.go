package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	gandalf "github.com/nassorc/gandalf"
)

func IAmMenu(world *gandalf.World, entities []*gandalf.Entity) {
	fmt.Println("I am menu")
}

type Menu struct {
	// currentScene Scene
	world *gandalf.World
}

func (g *Menu) Setup(world *gandalf.World) {
	g.world = world

	world.RegisterSystem(IAmMenu)
}

func InputSystem(world *gandalf.World, entities []*gandalf.Entity) {
	if rl.IsKeyDown(rl.KeyD) {
		fmt.Println("I ame input")
		world.ChangeScene(&Menu{})
	}
}

type Game struct {
	// currentScene Scene
	world *gandalf.World
}

type Position struct {
	X int
	Y int
}

func (g *Game) Setup(world *gandalf.World) {
	g.world = world

	world.RegisterSystem(InputSystem, &Position{})
}

func main() {
	var engine = gandalf.NewEngine(&Game{})

	engine.Run()
	fmt.Println("working")
}
