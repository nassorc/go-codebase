package scenes

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
	"github.com/nassorc/gandalf/cmd/components"
	"github.com/nassorc/gandalf/cmd/systems"
)

type Level1 struct {
	world  *gandalf.World
	player gandalf.EntityHandle
}

func (scene *Level1) Setup(world *gandalf.World) {
	scene.world = world

	world.LoadTexture("t_witch_run", "./resources/Blue_witch/B_witch_run.png")
	ok := world.LoadAnimation(
		"witch_run",
		"t_witch_run",
		8,
		rl.NewRectangle(0, 0, 32, 48),
		rl.NewVector2(32, 48),
		rl.NewVector2(0, 0), 1, 0, 10)

	if !ok {
		fmt.Println("could not load animation witch_idle")
		panic("could not load animation witch_idle")
	}

	world.RegisterComponents2(
		components.TransformID,
		components.RigidBodyID,
		components.MotionID,
		components.RectID,
		components.AnimationID,
	)

	world.RegisterSystem2(scene.PlayerInputSystem)
	world.RegisterSystem2(scene.PlayerCollisionSystem2, components.RigidBodyID, components.TransformID)
	world.RegisterSystem2(systems.DrawEntities, components.RectID, components.TransformID)
	world.RegisterSystem2(func(entities []gandalf.EntityHandle) {
		for _, entity := range entities {
			var T *components.Transform
			var R *components.RigidBody

			entity.Unpack(&T, &R)

			rl.DrawRectangleLines(int32(T.Pos.X), int32(T.Pos.Y), int32(R.Size.X), int32(R.Size.Y), rl.Red)
			rl.DrawCircle(int32(T.Pos.X+R.Size.X/2), int32(T.Pos.Y+R.Size.Y/2), 2, rl.Yellow)
		}
	}, components.RectID, components.TransformID)
	world.RegisterSystem2(func(entities []gandalf.EntityHandle) {
		for _, entity := range entities {
			var T *components.Transform
			var A *components.Animation

			entity.Unpack(&T, &A)

			anim, _ := world.GetAnimation(A.Name)
			rl.DrawTexturePro(anim.Texture, anim.Src, rl.NewRectangle(T.Pos.X, T.Pos.Y, anim.FrmSize.X*anim.Scale, anim.FrmSize.Y*anim.Scale), rl.Vector2{10, 10}, anim.Rotation, rl.White)
		}
	}, components.AnimationID, components.TransformID)
	// world.RegisterSystem(scene.MovementSystem, &components.Motion{}, &components.Transform{})

	// world.CreateEntity(
	// 	&components.Transform{Pos: rl.Vector2{0, 0}, PrevPos: rl.Vector2{0, 0}},
	// 	&components.Animation{"witch_run"},
	// )

	// player
	var size = rl.NewVector2(10, 28)
	var pos = rl.NewVector2(0, 0)

	scene.player = world.CreateEntity(
		&components.Transform{Pos: pos},
		&components.RigidBody{Size: size, HalfSize: rl.Vector2Scale(size, 0.5), Offset: rl.Vector2Scale(size, 0.5)},
		&components.Rect{Size: size, Color: rl.Lime},
		&components.Animation{"witch_run"},
	)

	size = rl.NewVector2(100, 100)
	pos = rl.NewVector2(200, 200)

	world.CreateEntity(
		&components.Transform{Pos: pos, PrevPos: pos},
		&components.RigidBody{Size: size, HalfSize: rl.Vector2Scale(size, 0.5), Offset: rl.Vector2Scale(size, 0.5)},
		&components.Rect{Size: size, Color: rl.Red},
	)

	// platform
	size = rl.NewVector2(800, 64)
	pos = rl.NewVector2(0, 576-64)

	world.CreateEntity(
		&components.Transform{Pos: pos, PrevPos: pos},
		&components.RigidBody{Size: size, HalfSize: rl.Vector2Scale(size, 0.5), Offset: rl.Vector2Scale(size, 0.5)},
		&components.Rect{Size: size, Color: rl.White},
	)
}

func (scene *Level1) PlayerInputSystem(entities []gandalf.EntityHandle) {
	var transform *components.Transform
	scene.player.Unpack(&transform)

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
	if rl.IsKeyPressed(rl.KeyF1) {
		scene.world.PushScene(&Inventory{})
	}

	transform.PrevPos = transform.Pos
	transform.Pos = rl.Vector2Add(transform.Pos, rl.Vector2{X: mx, Y: my})
}

func (scene *Level1) PlayerCollisionSystem2(entities []gandalf.EntityHandle) {
	for _, entity := range entities {
		if scene.player.Entity() == entity.Entity() {
			continue
		}
		entityAABBCollision(scene.player, entity)
	}
}

func getOverlap(pos1, pos2, ahs, bhs rl.Vector2) rl.Vector2 {
	var distAB = rl.NewVector2(float32(math.Abs(float64(pos1.X-pos2.X))), float32(math.Abs(float64(pos1.Y-pos2.Y))))

	var dx = (ahs.X + bhs.X) - distAB.X // if dx > 0, then vertical overlap
	var dy = (ahs.Y + bhs.Y) - distAB.Y // if dy > 0, then horizontal overlap

	return rl.NewVector2(dx, dy)
}

func entityAABBCollision(target gandalf.EntityHandle, object gandalf.EntityHandle) {
	// get entity target and object data
	var (
		T1 *components.Transform
		T2 *components.Transform
		R1 *components.RigidBody
		R2 *components.RigidBody
	)

	target.Unpack(&T1, &R1)
	object.Unpack(&T2, &R2)

	if T1 == nil || T2 == nil || R1 == nil || R2 == nil {
		return
	}

	// has collision if the XY distance between A and B is less than the sum of
	// both their half sizes
	var pos1 = rl.Vector2Add(T1.Pos, R1.Offset)
	var pos2 = rl.Vector2Add(T2.Pos, R2.Offset)

	// collision if we have target vertical and horizontal overlap
	var overlap = getOverlap(pos1, pos2, R1.HalfSize, R2.HalfSize) // amount of overlap between A and B

	// no collision
	if overlap.X <= 0 || overlap.Y <= 0 {
		return
	}

	// Use previous position to determine how to resolve collision
	var prev1 = rl.Vector2Add(T1.PrevPos, R1.Offset)
	var prev2 = rl.Vector2Add(T2.PrevPos, R2.Offset)
	var prevOverlap = getOverlap(prev1, prev2, R1.HalfSize, R2.HalfSize)

	// collision resolution
	if prevOverlap.X > 0 {
		if pos1.Y < pos2.Y {
			T1.Pos.Y -= overlap.Y
		} else {
			T1.Pos.Y += overlap.Y
		}
	} else if prevOverlap.Y > 0 {
		if pos1.X < pos2.X {
			T1.Pos.X -= overlap.X
		} else {
			T1.Pos.X += overlap.X
		}
	}
}

func toCenterOrigin(pos rl.Vector2, size rl.Vector2) rl.Vector2 {
	return rl.NewVector2(
		pos.X-(size.X/2),
		pos.Y-(size.Y/2),
	)
}
