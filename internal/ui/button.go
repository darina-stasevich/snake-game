package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"snake-game/internal/assets"
)

type Button struct {
	X, Y, Width, Height float64
	Text                string
	Color               color.Color
	HoverColor          color.Color

	OnClick func()

	IsHovered bool
}

func NewButton(x, y, w, h float64, text string, onClick func()) *Button {
	return &Button{
		X:     x,
		Y:     y,
		Width: w, Height: h,
		Text:       text,
		Color:      color.RGBA{R: 0x8a, G: 0x2b, B: 0xe2, A: 0xff},
		HoverColor: color.RGBA{R: 0x99, G: 0x32, B: 0xcc, A: 0xff},
		OnClick:    onClick,
		IsHovered:  false,
	}
}

func (button *Button) Update() {
	cursorX, cursorY := ebiten.CursorPosition()

	if float64(cursorX) >= button.X && float64(cursorX) < button.X+button.Width &&
		float64(cursorY) >= button.Y && float64(cursorY) < button.Y+button.Height {
		button.IsHovered = true
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if button.OnClick != nil {
				button.OnClick()
			}
		}
	} else {
		button.IsHovered = false
	}
}

func (button *Button) Draw(screen *ebiten.Image, assets *assets.Assets) {
	currentColor := button.Color
	if button.IsHovered {
		currentColor = button.HoverColor
	}

	DrawRectangle(screen, assets, button.X, button.Y, button.Width, button.Height, currentColor)
	font := assets.UIFont
	bounds := text.BoundString(font, button.Text)
	textX := button.X + (button.Width-float64(bounds.Dx()))/2
	textY := button.Y + (button.Height+float64(bounds.Dy()))/2
	text.Draw(screen, button.Text, font, int(textX), int(textY), color.White)
}
