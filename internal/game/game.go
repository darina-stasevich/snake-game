package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log/slog"
	"snake-game/internal/assets"
	"snake-game/internal/config"
	"snake-game/internal/core"
	"snake-game/internal/scenes"
	"time"
)

type Game struct {
	cfg    *config.Config
	assets *assets.Assets

	score    int
	gameTime time.Duration

	logger *slog.Logger

	scenes       map[core.GameState]scenes.Scene
	currentScene scenes.Scene
	playingScene *scenes.PlayingScene
}

func (g *Game) Config() *config.Config {
	return g.cfg
}

func (g *Game) Assets() *assets.Assets {
	return g.assets
}

func (g *Game) Logger() *slog.Logger {
	return g.logger
}

func (g *Game) Score() int {
	return g.score
}

func (g *Game) GameTime() time.Duration {
	return g.gameTime
}

func NewGame(cfg *config.Config, assets *assets.Assets) (*Game, error) {
	g := &Game{
		cfg:    cfg,
		assets: assets,
		logger: cfg.Logger,
	}

	playingScene, err := scenes.NewPlayingScene(g)

	if err != nil {
		return nil, err
	}

	gameOverScene := scenes.NewGameOverScene(g)

	g.scenes = map[core.GameState]scenes.Scene{
		core.GamePlayingState: playingScene,
		core.GameOverState:    gameOverScene,
	}

	g.currentScene = playingScene
	g.playingScene = playingScene
	g.currentScene.OnEnter()

	if err := g.Reset(); err != nil {
		g.logger.Error("failed to initialize game on creation", "error", err)
		return nil, fmt.Errorf("не удалось инициализировать игру: %w", err)
	} else {
		g.logger.Info("game created successfully",
			slog.Group("config",
				"screen_width", cfg.ScreenWidth,
				"screen_height", cfg.ScreenHeight,
				"tile_size", cfg.TileSize,
				"initial_snake_len", cfg.InitialSnakeLen,
				"initial_speed", cfg.InitialSpeed,
			),
		)

	}
	return g, nil
}

func (g *Game) Reset() error {
	g.score = 0
	g.gameTime = 0
	err := g.playingScene.Reset()
	if err != nil {
		g.logger.Warn("failed to restart game")
		return err
	}

	g.logger.Info("game restarted successfully")
	return nil
}

func (g *Game) Update() error {
	newState, err := g.currentScene.Update()
	if err != nil {
		return err
	}
	if g.scenes[newState] != g.currentScene {
		newScene, ok := g.scenes[newState]
		if !ok {
			return fmt.Errorf("unknown game state: %v", newState)
		}

		g.logger.Info("changing scene")
		g.currentScene = newScene

		g.currentScene.OnEnter()
	}
	if g.currentScene == g.playingScene {
		g.gameTime += time.Second / time.Duration(ebiten.TPS())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentScene.Draw(screen)
}

func (g *Game) NotifyFoodEaten() bool {
	g.score++
	g.logger.Info("snake ate food", "new_score", g.score)

	if g.cfg.SpeedIncreaseInterval > 0 && g.score%g.cfg.SpeedIncreaseInterval == 0 {
		g.logger.Info("speed increase condition met")
		return true
	}

	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.ScreenWidth, g.cfg.WindowHeight()
}
