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
}

func (c *Camera) Follow(x float32, y float32) {
	sw := c.ViewPort.X2 / float32(c.Zoom)
	sh := c.ViewPort.Y2 / float32(c.Zoom)

	x -= sw / 2
	y -= sh / 2

	c.Position.X = x
	c.Position.Y = y
}
