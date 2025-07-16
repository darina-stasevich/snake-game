package scenes

import "snake-game/internal/core"

type Wall struct {
	core.Position
}

func NewWall(x, y int) *Wall {
	return &Wall{Position: core.Position{X: x, Y: y}}
}
