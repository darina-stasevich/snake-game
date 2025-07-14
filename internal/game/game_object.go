package game

import "github.com/hajimehoshi/ebiten/v2"

type Position struct {
	X, Y int
}

type GameObject struct {
	Position
	Img *ebiten.Image
}
