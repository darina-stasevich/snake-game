package game

import "github.com/hajimehoshi/ebiten/v2"

type Wall struct {
	GameObject
}

func NewWall(x, y int, img *ebiten.Image) *Wall {
	return &Wall{GameObject: GameObject{Position: Position{X: x, Y: y}, Img: img}}
}
