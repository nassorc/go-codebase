package gandalf

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

type AABB struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
}
