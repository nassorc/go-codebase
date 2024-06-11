package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerSystem struct {
}

func (s *PlayerSystem) Update(w *World) {
	// assumes that it will receive the player entity as the first element
	// alternative solution:
	// 1. the AddEntity method doesn't require the system to
	// have a list of entities. Update it so that it only deals with a single
	// entity.
	// 2. the world keeps a reference to the player, since the player is a special
	// type of entity that would likely be used in many systems.
	// 3. a query system to fetch for specific entities.
}

type PhysicsSystem struct {
	entities []*Entity
}

type AxisCollision struct {
	X bool
	Y bool
}

func HasAxisCollision(posA rl.Vector2, posB rl.Vector2, sizeA rl.Vector2, sizeB rl.Vector2) (AxisCollision, rl.Vector2) {
	var halfSizeA = rl.Vector2Scale(sizeA, 0.5)
	var halfSizeB = rl.Vector2Scale(sizeB, 0.5)

	var distX float64 = math.Abs(float64(posA.X - posB.X))
	var distY float64 = math.Abs(float64(posA.Y - posB.Y))
	var deltaX = float32(distX) - (halfSizeA.X + halfSizeB.X)
	var deltaY = float32(distY) - (halfSizeA.Y + halfSizeB.Y)

	var hasXAxisCollision = deltaX < 0
	var hasYAxisCollision = deltaY < 0

	return AxisCollision{X: hasXAxisCollision, Y: hasYAxisCollision}, rl.NewVector2(float32(math.Abs(float64(deltaX))), float32(math.Abs(float64(deltaY))))
	// return AxisCollision{hasXAxisCollision},  AxisCollision{hasYAxisCollision}
}

func (s *PhysicsSystem) Update(w *World) {
	// Collision logic won't work with multiple entities since it directly updates
	// the movable entity's position.
	for _, entityA := range s.entities {
		// process movable entity colliding against entity with a rigid body and a transform component
		var movable *Movable
		entityA.GetData(&movable)

		if movable == nil {
			continue
		}

		for _, entityB := range s.entities {
			if entityA.id == entityB.id {
				continue
			}

			var transformA *Transform
			var rigidBodyA *RigidBody
			entityA.GetData(&transformA, &rigidBodyA)

			var transformB *Transform
			var rigidBodyB *RigidBody
			entityB.GetData(&transformB, &rigidBodyB)

			collision, overlap := HasAxisCollision(transformA.pos, transformB.pos, rigidBodyA.size, rigidBodyB.size)
			prevCollision, _ := HasAxisCollision(transformA.prevPos, transformB.prevPos, rigidBodyA.size, rigidBodyB.size)

			hasCollision := collision.X && collision.Y
			// Resolve movable (entityA) collision

			// vertical collision
			if hasCollision && prevCollision.X {
				if transformA.pos.Y < transformB.pos.Y {
					transformA.pos.Y -= overlap.Y
				} else {
					transformA.pos.Y += overlap.Y
				}
				transformA.prevPos.Y = transformA.pos.Y
			} else if hasCollision && prevCollision.Y {
				// horizontal collision
				transformA.prevPos.X = transformA.pos.X
				if transformA.pos.X < transformB.pos.X {
					transformA.pos.X -= overlap.X
				} else {
					transformA.pos.X += overlap.X
				}
			}

			// DEBUG
			var color *Color
			entityB.GetData(&color)
			if hasCollision {
				color.c = rl.Red
			} else {
				color.c = rl.Black
			}

		}
	}
}

func (s *PhysicsSystem) AddEntity(e *Entity) {
	s.entities = append(s.entities, e)
}

type InputSystem struct {
	entities []*Entity
}

func (s *InputSystem) Update(w *World) {
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

		transform.prevPos = transform.pos
		transform.pos.X += mx
		transform.pos.Y += my
	}
}

func (s *InputSystem) AddEntity(e *Entity) {
	s.entities = append(s.entities, e)
}
