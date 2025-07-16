package assets

import (
	"embed"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png" // Для поддержки формата PNG
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed images/* fonts/*
var assetsFS embed.FS

type Assets struct {
	SnakeHead       *ebiten.Image
	SnakeBody       *ebiten.Image
	SnakeBodyCorner *ebiten.Image
	SnakeTail       *ebiten.Image
	Apple           *ebiten.Image
	Wall            *ebiten.Image

	UIFont font.Face
}

func loadImage(path string) (*ebiten.Image, error) {
	file, err := assetsFS.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func Load(skin string) (*Assets, error) {
	var err error
	assets := &Assets{}

	skinPath := filepath.Join("images", skin)

	assets.SnakeHead, err = loadImage(filepath.Join(skinPath, "head.png"))
	if err != nil {
		return nil, err
	}
	assets.SnakeBody, err = loadImage(filepath.Join(skinPath, "body.png"))
	if err != nil {
		return nil, err
	}
	assets.SnakeBodyCorner, err = loadImage(filepath.Join(skinPath, "body_corner.png"))
	if err != nil {
		return nil, err
	}
	assets.SnakeTail, err = loadImage(filepath.Join(skinPath, "tail.png"))
	if err != nil {
		return nil, err
	}

	assets.Apple, err = loadImage(filepath.Join(skinPath, "food.png"))

	assets.Wall = ebiten.NewImage(20, 20)
	assets.Wall.Fill(color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff}) // Серая стена

	fontData, err := assetsFS.ReadFile("fonts/PressStart2P-Regular.ttf")
	if err != nil {
		return nil, err
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	const dpi = 72
	assets.UIFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	return assets, nil
}
