package main

import (
	"fmt"
	"image/color"

	"gandalf"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// type IEvent interface {}
//
// type PubSub struct {
//   a map[reflec.Type]
// }

func main() {

  // eventBus.publish(CollisionEvent{ a: EntityA, b: EntityB })

  fmt.Println()
	rl.InitWindow(960, 640, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

  world := gandalf.NewWorld()
  world.RegisterComponents(&Tag{}, &Transform{}, &Size{}, &Color{}, &Input{})

  world.RegisterSystem(&InputSystem{})

  player := world.NewEntity(&Transform{ pos: rl.NewVector2(0, 0) }, &Tag{ "player" }, &Size{32, 32}, &Color{rl.Green}, &Input{})

  // fmt.Println("system", *(world.systems[0]))

  var playerTag *Tag
  var playerPos *Transform
  player.GetData(&playerTag, &playerPos)

	for !rl.WindowShouldClose() {
    // Update Player___________________________________________________________
    // UpdatePlayer(player, world)
    // Update__________________________________________________________________

    world.Run()

		rl.BeginDrawing()
    rl.ClearBackground(color.RGBA{R: 24, G: 24, B: 24})

    for _, entity := range world.Entities {
      var size *Size
      var transform *Transform
      var color *Color
      entity.GetData(&size, &transform, &color)

      var pos = transform.pos

      rl.DrawRectangle(int32(pos.X), int32(pos.Y), int32(size.Width), int32(size.Height), color.c)
    }

    rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

    rl.DrawText("hello", 32, 32, 20, rl.Black)

		rl.EndDrawing()
	}
}

type PlayerSystem struct {
}

func (s *PlayerSystem) Update(w *gandalf.World) {
  // assumes that it will receive the player entity as the first element
  // alternative solution: 
  // 1. the AddEntity method doesn't require the system to
  // have a list of entities. Update it so that it only deals with a single
  // entity.
  // 2. the world keeps a reference to the player, since the player is a special
  // type of entity that would likely be used in many systems.
  // 3. a query system to fetch for specific entities. 
}

type InputSystem struct {
  entities []*gandalf.Entity
}

func (s *InputSystem) Update(w *gandalf.World) {
  for _, entity := range s.entities {
    var transform *Transform
    entity.GetData(&transform)

    const speed = 2
    var mx float32 = 0
    var my float32 = 0

    if rl.IsKeyDown(rl.KeyW) {
      my = -speed
    }
    if rl.IsKeyDown(rl.KeyS) {
      my = speed
    }
    if rl.IsKeyDown(rl.KeyA) {
      mx = -speed
    }
    if rl.IsKeyDown(rl.KeyD) {
      mx = speed
    }

    transform.pos.X += mx
    transform.pos.Y += my
  }

  player := s.entities[0]
  var transform *Transform
  player.GetData(&transform)

  const speed = 2
  var mx float32 = 0
  var my float32 = 0

  if rl.IsKeyDown(rl.KeyW) {
    my = -speed
  }
  if rl.IsKeyDown(rl.KeyS) {
    my = speed
  }
  if rl.IsKeyDown(rl.KeyA) {
    mx = -speed
  }
  if rl.IsKeyDown(rl.KeyD) {
    mx = speed
  }

  transform.pos.X += mx
  transform.pos.Y += my
}

func (s *InputSystem) AddEntity(e *gandalf.Entity) {
  s.entities = append(s.entities, e)
}


type Transform struct {
  pos rl.Vector2
}

type Tag struct {
  name string
}

type Size struct {
  Width int
  Height int
}

type Color struct {
  c rl.Color
}

type Input struct {
  Up      bool
  Down    bool
  Left    bool
  Right   bool
  CanJump bool
}
