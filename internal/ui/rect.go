package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"snake-game/internal/assets"
)

func DrawRectangle(screen *ebiten.Image, assets *assets.Assets, x, y, width, height float64, clr color.Color) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(width, height)
	opts.GeoM.Translate(x, y)
	opts.ColorScale.ScaleWithColor(clr)
	screen.DrawImage(assets.WhitePixel, opts)
}
