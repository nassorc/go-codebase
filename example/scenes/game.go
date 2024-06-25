package scenes

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
	"github.com/nassorc/gandalf/cmd/components"
	"github.com/nassorc/gandalf/cmd/systems"
)

type Game struct {
	player *gandalf.EntityHandle
	world  *gandalf.World
}

func (g *Game) Setup(world *gandalf.World) {
	g.world = world
	world.RegisterComponents(
		&components.Transform{},
		&components.Input{},
		&components.Motion{},
		&components.Poly{},
		&components.Lifespan{},
	)

	// render poly
	world.RegisterSystem(g.InputSystem)
	world.RegisterSystem(g.MovementSystem, &components.Motion{}, &components.Transform{})
	world.RegisterSystem(g.WallCollisionSystem, &components.Motion{}, &components.Transform{})
	world.RegisterSystem(systems.RotatePoly, &components.Poly{})
	world.RegisterSystem(systems.DrawPoly, &components.Poly{}, &components.Transform{})
	world.RegisterSystem(g.Debug)
	world.RegisterSystem(g.LifespanSystem, &components.Lifespan{})

	g.spawnEnemy(100, 100)
	g.spawnEnemy(150, 100)
	g.spawnEnemy(180, 100)
	g.spawnEnemy(200, 100)
	g.spawnEnemy(300, 100)
	g.spawnPlayer()
}

func (g *Game) spawnPlayer() {
	g.player = g.world.CreateEntity(
		components.NewPoly(3, 20, 0, rl.Green),
		components.NewTranform(250, 250),
	)
}

func (g *Game) spawnEnemy(x float32, y float32) {
	sides := rand.Int31n(5) + 2

	R := uint8(rand.Intn(100)) + 154
	G := uint8(rand.Intn(50))
	B := uint8(rand.Intn(100))

	X := float32(rand.Intn(10)) - 5
	Y := float32(rand.Intn(10)) - 5

	g.world.CreateEntity(
		components.NewPoly(sides, 20, 0, rl.NewColor(R, G, B, 255)),
		components.NewTranform(x, y),
		components.NewMotion(rl.NewVector2(X, Y), rl.NewVector2(0, 0)),
	)
}

func (g *Game) spawnBullet(src rl.Vector2, dst rl.Vector2) {
	g.world.CreateEntity(
		components.NewPoly(12, 8, 0, rl.White),
		components.NewTranform(src.X, src.Y),
		components.NewMotion(rl.NewVector2(dst.X, dst.Y), rl.NewVector2(0, 0)),
		components.NewLifespan(60),
	)
}

func (g *Game) Debug(w *gandalf.World, entities []*gandalf.EntityHandle) {
	fps := rl.GetFPS()
	rl.DrawText(fmt.Sprintf("%v", fps), 10, 10, 28, rl.Red)
}

func (g *Game) InputSystem(w *gandalf.World, entities []*gandalf.EntityHandle) {
	var transform *components.Transform
	var poly *components.Poly

	g.player.Unpack(&transform, &poly)

	// world.RemoveComponent(g.player, &Position{})

	var mx float32 = 0
	var my float32 = 0
	var speed float32 = 5

	// Movement
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

	var pos = transform.Pos

	pos.X += mx
	pos.Y += my

	// mouse click
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		// get velocity vector
		dstPos := rl.Vector2Subtract(rl.GetMousePosition(), *pos)
		// normalize
		dstPos = rl.Vector2Normalize(dstPos)
		// Give speed
		dstPos = rl.Vector2Scale(dstPos, 8)

		g.spawnBullet(*pos, dstPos)
	}
}
func (g *Game) MovementSystem(world *gandalf.World, entities []*gandalf.EntityHandle) {
	for _, entity := range entities {
		var motion *components.Motion
		var transform *components.Transform
		entity.Unpack(&motion, &transform)

		*transform.Pos = rl.Vector2Add(*transform.Pos, motion.Velocity)
	}
}

func (g *Game) WallCollisionSystem(world *gandalf.World, entities []*gandalf.EntityHandle) {
	for _, entity := range entities {
		var motion *components.Motion
		var transform *components.Transform
		entity.Unpack(&motion, &transform)

		var pos = transform.Pos

		if pos.X <= 0 || pos.X >= 500 {
			motion.Velocity.X *= -1
		}
		if pos.Y <= 0 || pos.Y >= 500 {
			motion.Velocity.Y *= -1
		}
	}
}

func (g *Game) LifespanSystem(world *gandalf.World, entities []*gandalf.EntityHandle) {
	for _, entity := range entities {
		var lifespan *components.Lifespan
		entity.Unpack(&lifespan)

		if lifespan.Remaining > 0 {
			lifespan.Remaining -= 1
		}
		if lifespan.Remaining == 0 {
			// world.
			// world.RemoveComponent(, &lifespan{})
			// entity.Destroy()
		}
	}
}
