package components

import "github.com/nassorc/gandalf"

func NewInput() *Input {
	return &Input{}
}

type Input struct {
	Up       bool
	Down     bool
	Left     bool
	Right    bool
	Shoot    bool
	CanShoot bool
}

var InputID = gandalf.CreateComponentID[Input]()
