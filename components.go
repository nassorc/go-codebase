package main

import rl "github.com/gen2brain/raylib-go/raylib"

type RigidBody struct {
	size   rl.Vector2
	offset rl.Vector2
}

type Transform struct {
	vel     rl.Vector2
	pos     rl.Vector2
	prevPos rl.Vector2
}

type Movable struct{}

type Tag struct {
	name string
}

type Size struct {
	Width  int
	Height int
}

type Color struct {
	c rl.Color
}

type Input struct {
	Up      bool
	Down    bool
	Left    bool
	Right   bool
	CanJump bool
}
