package core

type Food struct {
	Position
}

func NewFood(x, y int) *Food {
	return &Food{Position: Position{X: x, Y: y}}
}
