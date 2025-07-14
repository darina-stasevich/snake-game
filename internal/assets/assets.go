package assets

import (
	"fmt"
	"image/color"
	_ "image/png" // Для поддержки формата PNG

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Assets хранит все загруженные игровые ресурсы.
// Пока добавим только базовые, потом можно будет расширить.
type Assets struct {
	SnakeHead *ebiten.Image
	SnakeBody *ebiten.Image
	Apple     *ebiten.Image
	Wall      *ebiten.Image
}

func loadImage(path string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить изображение %s: %w", path, err)
	}
	return img, nil
}

// Load загружает все игровые ассеты.
// Пути к файлам должны быть правильными.
func Load() (*Assets, error) {
	// var err error
	assets := &Assets{}

	// Предполагается, что у вас есть эти файлы в папке internal/assets/images/
	// assets.SnakeHead, err = loadImage("internal/assets/images/head.png")
	// if err != nil { return nil, err }
	// ... и так далее для остальных

	// ЗАГЛУШКА: Если картинок пока нет, создадим их программно
	// Когда у вас появятся файлы, замените этот блок на загрузку.
	assets.SnakeHead = ebiten.NewImage(20, 20)
	assets.SnakeHead.Fill(color.RGBA{R: 0, G: 0xff, B: 0, A: 0xff}) // Зеленая голова
	assets.SnakeBody = ebiten.NewImage(20, 20)
	assets.SnakeBody.Fill(color.RGBA{R: 0, G: 0xcc, B: 0, A: 0xff}) // Тело чуть темнее
	assets.Apple = ebiten.NewImage(20, 20)
	assets.Apple.Fill(color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff}) // Красное яблоко
	assets.Wall = ebiten.NewImage(20, 20)
	assets.Wall.Fill(color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff}) // Серая стена

	return assets, nil
}
