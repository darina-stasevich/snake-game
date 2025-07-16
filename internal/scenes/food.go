package scenes

import "snake-game/internal/core"

type Food struct {
	core.Position
}

func NewFood(x, y int) *Food {
	return &Food{Position: core.Position{X: x, Y: y}}
}
