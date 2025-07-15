package assets

import (
	"fmt"
	"image/color"
	_ "image/png" // Для поддержки формата PNG
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Assets struct {
	SnakeHead       *ebiten.Image
	SnakeBody       *ebiten.Image
	SnakeBodyCorner *ebiten.Image
	SnakeTail       *ebiten.Image
	Apple           *ebiten.Image
	Wall            *ebiten.Image
}

func loadImage(path string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить изображение %s: %w", path, err)
	}
	return img, nil
}

func Load(skin string) (*Assets, error) {
	var err error
	assets := &Assets{}

	skinPath := filepath.Join("internal", "assets", "images", skin)

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

	assets.Apple = ebiten.NewImage(20, 20)
	assets.Apple.Fill(color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff}) // Красное яблоко
	assets.Wall = ebiten.NewImage(20, 20)
	assets.Wall.Fill(color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff}) // Серая стена

	return assets, nil
}
