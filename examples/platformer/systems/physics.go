package systems

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	ecs "github.com/nassorc/gandalf/ecs"
	c "github.com/nassorc/gandalf/examples/platformer/components"
)

func PhysicsSystem(world *ecs.World, entities []*ecs.Entity) {
	// Collision logic won't work with multiple entities since it directly updates
	// the movable entity's position.
	for _, entityA := range entities {
		// process movable entity colliding against entity with a rigid body and a transform component
		var movable *c.Movable
		entityA.GetData(&movable)

		if movable == nil {
			continue
		}

		for _, entityB := range entities {
			if entityA.Id() == entityB.Id() {
				continue
			}

			var transformA *c.Transform
			var rigidBodyA *c.RigidBody
			entityA.GetData(&transformA, &rigidBodyA)

			var transformB *c.Transform
			var rigidBodyB *c.RigidBody
			entityB.GetData(&transformB, &rigidBodyB)

			collision, overlap := HasAxisCollision(transformA.Pos, transformB.Pos, rigidBodyA.Size, rigidBodyB.Size)
			prevCollision, _ := HasAxisCollision(transformA.PrevPos, transformB.PrevPos, rigidBodyA.Size, rigidBodyB.Size)

			hasCollision := collision.X && collision.Y
			// Resolve movable (entityA) collision

			// vertical collision
			if hasCollision && prevCollision.X {
				if transformA.Pos.Y < transformB.Pos.Y {
					transformA.Pos.Y -= overlap.Y
				} else {
					transformA.Pos.Y += overlap.Y
				}
				transformA.PrevPos.Y = transformA.Pos.Y
			} else if hasCollision && prevCollision.Y {
				// horizontal collision
				transformA.PrevPos.X = transformA.Pos.X
				if transformA.Pos.X < transformB.Pos.X {
					transformA.Pos.X -= overlap.X
				} else {
					transformA.Pos.X += overlap.X
				}
			}

			// DEBUG
			var color *c.Color
			entityB.GetData(&color)
			if hasCollision {
				color.C = rl.Red
			} else {
				color.C = rl.Black
			}

		}
	}

}

// type PhysicsSystem struct {
// 	entities []*ecs.Entity
// }

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

// func (s *PhysicsSystem) Update(w *ecs.World) {

// }

// func (s *PhysicsSystem) AddEntity(e *ecs.Entity) {
// 	s.entities = append(s.entities, e)
// }
