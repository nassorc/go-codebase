package math

import "math"

func NewVec2(x float32, y float32) Vec2 {
	return Vec2{x, y}
}

type Vec2 struct {
	X float32
	Y float32
}

func Vec2Add(v1, v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

func Vec2Subtract(v1, v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
	}
}

func Vec2Scale(v Vec2, scale float32) Vec2 {
	return Vec2{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

func Vec2Divide(v Vec2, divisor float32) Vec2 {
	return Vec2{
		X: v.X / divisor,
		Y: v.Y / divisor,
	}
}

func Vec2Normalize(v Vec2) Vec2 {
  m := Vec2Magnitude(v)

  return Vec2Divide(v, m)
}

func Vec2Magnitude(v Vec2) float32 {
  return float32(math.Sqrt(float64(v.X * v.X + v.Y * v.Y)))
}

type AABB struct {
  Min Vec2
  Max Vec2
}

func (aabb *AABB) Width() float32 {
  return aabb.Max.X - aabb.Min.X
}

func (aabb *AABB) Height() float32 {
  return aabb.Max.Y - aabb.Min.Y
}

func (aabb *AABB) Center() Vec2 {
  hw := aabb.Width()/2
  hh := aabb.Height()/2

  return Vec2{ aabb.Min.X + hw, aabb.Min.Y + hh }
}
