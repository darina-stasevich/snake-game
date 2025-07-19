package main

import (
	"log/slog"
	"os"
	"snake-game/internal/assets"
	"snake-game/internal/config"
	"snake-game/internal/game"
	"snake-game/internal/storage"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Форматируем время для удобства
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	// 1. Загружаем конфигурацию
	cfg := config.LoadConfig()
	cfg.SetLogger(logger)

	// 2. Загружаем ассеты (картинки, шрифты)
	assets, err := assets.Load("snake")
	if err != nil {
		logger.Error("Failed to initialize assets", "err", err)
		os.Exit(1)
	} else {
		logger.Info("Assets successfully loaded")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		logger.Warn("DATABASE_URL environment variable is not set. Running without database.")
	}

	var repo storage.Repository

	if connStr != "" {
		repo, err = storage.NewPostgresRepository(connStr, logger)
		if err != nil {
			logger.Error("failed to connect to database", "error", err)
			os.Exit(1) // Выходим, если не можем подключиться к БД
		}
		defer repo.Close() // Гарантируем закрытие соединения при выходе
	}

	g, err := game.NewGame(cfg, assets, repo)
	if err != nil {
		logger.Error("failed to create game", "err", err)
		os.Exit(1)
	} else {
		logger.Info("Game successfully initialized")
	}

	// 4. Настраиваем и запускаем окно
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight+cfg.TopBarHeight)
	ebiten.SetWindowTitle("Змейка на Ebitengine")

	if err := ebiten.RunGame(g); err != nil {
		logger.Error("game finished with error", "error", err)
	}
}
