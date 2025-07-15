package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log/slog"
	"math/rand/v2"
	"snake-game/internal/assets"
	"snake-game/internal/config"
)

type Game struct {
	cfg    *config.Config
	assets *assets.Assets
	snake  *Snake
	food   *Food
	walls  []Wall

	score int
	ticks int

	logger *slog.Logger
}

func NewGame(cfg *config.Config, assets *assets.Assets) (*Game, error) {
	g := &Game{
		cfg:    cfg,
		assets: assets,
		logger: cfg.Logger,
	}
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
	startPos := Position{
		X: (g.cfg.ScreenWidth / g.cfg.TileSize) / 2,
		Y: (g.cfg.ScreenHeight / g.cfg.TileSize) / 2,
	}

	snake, err := NewSnake(startPos.X, startPos.Y, g.cfg.InitialSnakeLen, g.cfg.InitialSpeed)
	if err != nil {
		g.logger.Error("failed to create snake", "error", err)
		return fmt.Errorf("не удалось создать змею: %w", err)
	}

	g.snake = snake

	g.score = 0
	g.ticks = 0
	g.spawnFood()
	// TODO: g.loadLevel(1)

	g.logger.Info("game restarted successfully")
	return nil
}

func (g *Game) Update() error {

	if !g.snake.IsAlive {
		return nil
	}

	g.ticks++
	g.handleInput()
	moved := g.snake.Update()
	if moved {
		err := g.checkCollisions()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.snake.SetNextDirection(Up)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.snake.SetNextDirection(Down)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.snake.SetNextDirection(Left)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.snake.SetNextDirection(Right)
	}
}

func (g *Game) spawnFood() {
	occupiedCells := make(map[Position]bool)
	for _, wall := range g.walls {
		occupiedCells[wall.Position] = true
	}
	for _, snakePart := range g.snake.Body {
		occupiedCells[snakePart.Position] = true
	}

	freeCells := make([]Position, 0)
	maxWidth := g.cfg.ScreenWidth / g.cfg.TileSize
	maxHeight := g.cfg.ScreenHeight / g.cfg.TileSize
	for i := 0; i < maxWidth; i++ {
		for j := 0; j < maxHeight; j++ {
			if (occupiedCells[Position{i, j}] == false) {
				freeCells = append(freeCells, Position{i, j})
			}
		}
	}

	randomIndex := rand.IntN(len(freeCells))
	g.food = NewFood(freeCells[randomIndex].X, freeCells[randomIndex].Y)
	g.logger.Info("new food created", "position X", freeCells[randomIndex].X, "position Y", freeCells[randomIndex].Y)
}

func (g *Game) checkCollisions() error {

	for _, wall := range g.walls {
		if wall.Position == g.snake.Body[0].Position {
			g.snake.IsAlive = false
			g.logger.Info("snake crashed in wall")
			break
		}
	}

	if g.snake.Body[0].X < 0 || g.snake.Body[0].X >= g.cfg.ScreenWidth/g.cfg.TileSize {
		g.snake.IsAlive = false
		g.logger.Info("snake crashed in border")
	}
	if g.snake.Body[0].Y < 0 || g.snake.Body[0].Y >= g.cfg.ScreenHeight/g.cfg.TileSize {
		g.snake.IsAlive = false
		g.logger.Info("snake crashed in border")
	}

	if g.food.Position == g.snake.Body[0].Position {
		g.logger.Info("snake ate food")
		g.score++
		if g.score%g.cfg.SpeedIncreaseInterval == 0 {
			g.snake.DecreaseMoveInterval(g.cfg.SpeedIncreaseAmount)
			g.logger.Info("change snake speed")
		}
		g.spawnFood()
	} else {
		err := g.snake.cutTail()
		if err != nil {
			return err
		}
	}

	for i := 1; i < len(g.snake.Body); i++ {
		if g.snake.Body[i].Position == g.snake.Body[0].Position {
			g.snake.IsAlive = false
			g.logger.Info("snake crashed in own body")
			break
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xff})

	if g.snake.IsAlive {
		for i, segment := range g.snake.Body {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(segment.Position.X*g.cfg.TileSize), float64(segment.Position.Y*g.cfg.TileSize))

			img := g.assets.SnakeBody
			if i == 0 {
				img = g.assets.SnakeHead
				// TODO: Поворачивать голову в зависимости от направления
			}
			if i == len(g.snake.Body)-1 {
				// TODO: img = g.assets.SnakeTail
				// TODO: Поворачивать хвост в зависимости от состояния
			}
			// TODO: поворачивать тело в зависимости от положения
			screen.DrawImage(img, op)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.food.Position.X*g.cfg.TileSize), float64(g.food.Position.Y*g.cfg.TileSize))
	img := g.assets.Apple
	screen.DrawImage(img, op)

	// TODO: Рисовать стены, счет
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}
