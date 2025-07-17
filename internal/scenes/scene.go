package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log/slog"
	"snake-game/internal/assets"
	"snake-game/internal/config"
	"snake-game/internal/core"
	"time"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update() (core.GameState, error)
	OnEnter()
}

type GameAccessor interface {
	// Методы для доступа к общим ресурсам
	Config() *config.Config
	Assets() *assets.Assets
	Logger() *slog.Logger
	Score() int
	GameTime() time.Duration

	// Методы для управления состоянием
	NotifyFoodEaten() (shouldIncreaseSpeed bool)
	Reset() error
	StartGame(level *core.Level)
}
