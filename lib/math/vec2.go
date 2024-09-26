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

func (v *Vec2) Add(other Vec2) *Vec2 {
  v.X += other.X
  v.Y += other.Y
  return v
}

func Vec2Sub(v1, v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
	}
}

func (v *Vec2) Sub(other Vec2) *Vec2 {
  v.X -= other.X
  v.Y -= other.Y
  return v
}

func Vec2Scale(v Vec2, scale float32) Vec2 {
	return Vec2{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

func (v *Vec2) Scale(scalar float32) *Vec2 {
  v.X *= scalar
  v.Y *= scalar
  return v
}

func Vec2Normalize(v Vec2) Vec2 {
  return Vec2Scale(v, (1 / v.Length()))
}

func (v *Vec2) Normalize() *Vec2 {
  return v.Scale(1/v.Length())
}

func Vec2Length(v Vec2) float32 {
  return float32(math.Sqrt(float64(v.X * v.X + v.Y * v.Y)))
}

func (v *Vec2) Length() float32 {
  return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

type Rotation struct {
  C float32
  S float32
}

func SinCos(angle float32) Rotation {
  a := float64(angle)
  return Rotation{float32(math.Cos(a)), float32(math.Sin(a))}
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
