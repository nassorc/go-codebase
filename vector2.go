package gandalf

type Vec2 struct {
	X float64
	Y float64
}

func Vec2Add(v1, v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

func Vec2Scale(v Vec2, scale float64) Vec2 {
	return Vec2{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}
