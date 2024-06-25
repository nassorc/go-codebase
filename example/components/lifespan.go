package components

func NewLifespan(total int) *Lifespan {
	return &Lifespan{
		Remaining: total,
		Total:     total,
	}
}

type Lifespan struct {
	Remaining int
	Total     int
}
