package main

import (
	"fmt"

	ecs "github.com/nassorc/gandalf/ecs"
)

func NewPlayer(x int, y int) *Player {
	return &Player{
		Position: &Position{x, y},
	}
}

type Player struct {
	*Position
}

func main() {
	fmt.Println("working")
	world := ecs.NewWorld()
	world.RegisterComponents(&Position{})
	world.RegisterSystem(&MoveRight{})

	world.CreateEntity(&Position{67, 68})

	world.CreateEntityFromPreFab(NewPlayer(0, 5))

	fmt.Println(world.Components[0].Data)
	fmt.Println(world.Components[0].Data.Index(0))
	fmt.Println(world.Components[0].Data.Index(1))

	world.Run()
}

type Position struct {
	X int
	Y int
}

type MoveRight struct {
	entities []*ecs.Entity
}

func (s *MoveRight) Update(w *ecs.World)     {}
func (s *MoveRight) Render(w *ecs.World)     {}
func (s *MoveRight) AddEntity(e *ecs.Entity) {}
