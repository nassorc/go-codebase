package components

import "github.com/nassorc/gandalf"

type Animation struct {
	Name string
}

var AnimationID = gandalf.CreateComponentID[Animation]()
