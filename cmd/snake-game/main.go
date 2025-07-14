package main

import (
	"log"
	"snake-game/internal/assets"
	"snake-game/internal/config"
	"snake-game/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// 1. Загружаем конфигурацию
	cfg := config.LoadConfig()

	// 2. Загружаем ассеты (картинки, шрифты)
	assets, err := assets.Load()
	if err != nil {
		log.Fatalf("error loading assets: %v", err)
	}

	// 3. Создаем экземпляр игры
	g, err := game.NewGame(cfg, assets)
	if err != nil {
		log.Fatalf("error creating game: %v", err)
	}

	// 4. Настраиваем и запускаем окно
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle("Змейка на Ebitengine")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("game finished with error: %v", err)
	}
}
