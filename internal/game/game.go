package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
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

	// TODO: добавить score int
	// TODO: добавить ticks int

}

func NewGame(cfg *config.Config, assets *assets.Assets) (*Game, error) {
	g := &Game{
		cfg:    cfg,
		assets: assets,
	}
	if err := g.Reset(); err != nil {
		return nil, fmt.Errorf("не удалось инициализировать игру: %w", err)
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
		return fmt.Errorf("не удалось создать змею: %w", err)
	}

	g.snake = snake

	// TODO: add score
	// TODO: add ticks
	g.spawnFood()
	// TODO: g.loadLevel(1)

	return nil
}

func (g *Game) Update() error {
	// TODO ticks ++
	g.handleInput()
	g.snake.Update()
	// TODO: Добавить логику столкновений
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
	g.food = NewFood(freeCells[randomIndex].X, freeCells[randomIndex].Y, g.assets.Apple)
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
			}
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
