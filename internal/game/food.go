package game

import "github.com/hajimehoshi/ebiten/v2"

type Food struct {
	GameObject
}

func NewFood(x, y int, img *ebiten.Image) *Food {
	return &Food{GameObject: GameObject{Position: Position{X: x, Y: y}, Img: img}}
}
