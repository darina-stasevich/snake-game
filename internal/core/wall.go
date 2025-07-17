package core

type Wall struct {
	Position
}

func NewWall(x, y int) *Wall {
	return &Wall{Position: Position{X: x, Y: y}}
}
