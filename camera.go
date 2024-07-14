package gandalf

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	ViewPort AABB // X, Y, Width, Height
	Position Vec2 // TOP-LEFT point
	Zoom     float64
}

func (c *Camera) Draw(world *ebiten.Image, screen *ebiten.Image) {
	if c.Zoom == 0 {
		c.Zoom = 1
	}
	ops := &ebiten.DrawImageOptions{}

	ops.GeoM.Translate(-float64(c.Position.X), -float64(c.Position.Y))
	ops.GeoM.Scale(c.Zoom, c.Zoom)
	screen.DrawImage(world, ops)
}

func (c *Camera) Move(x float32, y float32) {
	c.Position.X += x
	c.Position.Y += y
	// sw := c.viewPort.X2 / float32(c.Zoom)
	// sh := c.viewPort.Y2 / float32(c.Zoom)
	// if c.Position.X < 0 {
	// 	c.Position.X = 0
	// }
	// if (c.Position.X + sw) > (worldWidth) {
	// 	c.Position.X = worldWidth - sw
	// }
	// if c.Position.Y < 0 {
	// 	c.Position.Y = 0
	// }
	// if (c.Position.Y + sh) > (worldHeight) {
	// 	c.Position.Y = worldHeight - sh
	// }
}
