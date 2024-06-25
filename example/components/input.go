package components

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
