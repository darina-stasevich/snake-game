package scenes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"math/rand/v2"
	"snake-game/internal/core"
	"snake-game/internal/ui"
)

type PlayingScene struct {
	snake *core.Snake
	food  *core.Food
	walls []core.Wall

	level *core.Level

	whitePixelImage *ebiten.Image

	accessor GameAccessor
}

func NewPlayingScene(accessor GameAccessor, level *core.Level) (*PlayingScene, error) {
	scene := &PlayingScene{
		accessor: accessor,
		level:    level,
		walls:    level.Walls,
	}

	err := scene.Reset()
	if err != nil {
		return nil, err
	}
	return scene, nil
}

func (p *PlayingScene) Reset() error {
	p.accessor.Logger().Info("playing scene  resetting...")
	cfg := p.accessor.Config()
	level := p.level
	var err error

	startPos := core.Position{
		X: level.GridWidth / 2,
		Y: level.GridHeight / 2,
	}

	snake, err := core.NewSnake(startPos.X, startPos.Y, cfg.InitialSnakeLen, cfg.InitialSpeed, cfg.MaxSpeed)
	if err != nil {
		p.accessor.Logger().Error("FATAL: failed to create snake during reset", "error", err)
		return fmt.Errorf("не удалось создать змею: %w", err)
	}

	p.accessor.Logger().Info("snake created successfully")
	p.snake = snake
	p.spawnFood()
	return nil
}

func (p *PlayingScene) spawnFood() {
	occupiedCells := make(map[core.Position]bool)
	for _, wall := range p.walls {
		occupiedCells[wall.Position] = true
	}
	for _, snakePart := range p.snake.Body {
		occupiedCells[snakePart.Position] = true
	}

	freeCells := make([]core.Position, 0)
	maxWidth := p.level.GridWidth
	maxHeight := p.level.GridHeight
	for i := 0; i < maxWidth; i++ {
		for j := 0; j < maxHeight; j++ {
			if (occupiedCells[core.Position{X: i, Y: j}] == false) {
				freeCells = append(freeCells, core.Position{X: i, Y: j})
			}
		}
	}

	if len(freeCells) == 0 {
		p.accessor.Logger().Warn("failed to created food: no free space left")
	} else {
		randomIndex := rand.IntN(len(freeCells))
		p.food = core.NewFood(freeCells[randomIndex].X, freeCells[randomIndex].Y)
		p.accessor.Logger().Info("new food created", "position X", freeCells[randomIndex].X, "position Y", freeCells[randomIndex].Y)
	}
}

func (p *PlayingScene) Update() (core.GameState, error) {
	p.accessor.Logger().Info("updating playing scene")
	if !p.snake.IsAlive {
		return core.GameOverState, nil
	}

	p.handleInput()

	moved := p.snake.Update()
	if moved {
		state, err := p.checkCollisions()
		if err != nil {
			return 0, err
		}
		return state, nil
	}
	return core.GamePlayingState, nil
}

func (p *PlayingScene) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		p.snake.SetNextDirection(core.Up)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		p.snake.SetNextDirection(core.Down)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		p.snake.SetNextDirection(core.Left)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		p.snake.SetNextDirection(core.Right)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		err := p.Reset()
		if err != nil {
			p.accessor.Logger().Error("failed to reset game", err)
			return
		}
		err = p.accessor.Reset()
		if err != nil {
			p.accessor.Logger().Error("failed to reset game", err)
		}
	}
}

func (p *PlayingScene) checkCollisions() (core.GameState, error) {
	logger := p.accessor.Logger()
	logger.Info("checking collisions")
	cfg := p.accessor.Config()
	level := p.level

	for _, wall := range p.walls {
		if wall.Position == p.snake.Body[0].Position {
			p.snake.IsAlive = false
			logger.Info("snake crashed in wall")
			break
		}
	}

	if p.snake.Body[0].X < 0 || p.snake.Body[0].X >= level.GridWidth {
		p.snake.IsAlive = false
		logger.Info("snake crashed in border")
	}
	if p.snake.Body[0].Y < 0 || p.snake.Body[0].Y >= level.GridHeight {
		p.snake.IsAlive = false
		logger.Info("snake crashed in border")
	}

	if p.food.Position == p.snake.Body[0].Position {
		logger.Info("snake ate food")
		needSpeedIncrease := p.accessor.NotifyFoodEaten()

		if needSpeedIncrease {
			p.snake.DecreaseMoveInterval(cfg.SpeedIncreaseAmount)
			logger.Info("change snake speed")
		}

		p.spawnFood()
	} else {
		err := p.snake.CutTail()
		if err != nil {
			return 0, err
		}
	}

	p.snake.CheckCollisionsWithSelf()
	if p.snake.IsAlive == false {
		logger.Info("snake crashed in its body")
	}

	if p.snake.IsAlive == false {
		return core.GameOverState, nil
	}

	return core.GamePlayingState, nil
}

func (p *PlayingScene) Draw(screen *ebiten.Image) {
	cfg := p.accessor.Config()
	assets := p.accessor.Assets()

	screen.Fill(color.RGBA{R: 5, G: 5, B: 15, A: 255})

	fieldWidthInPixels := float64(p.level.GridWidth * cfg.TileSize)
	fieldHeightInPixels := float64(p.level.GridHeight * cfg.TileSize)

	topBarHeight := float64(cfg.TopBarHeight)
	ui.DrawRectangle(
		screen,
		assets,
		0,
		topBarHeight,
		fieldWidthInPixels,
		fieldHeightInPixels,
		color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xff},
	)

	opBar := &ebiten.DrawImageOptions{}
	opBar.GeoM.Scale(float64(cfg.ScreenWidth), topBarHeight)
	screen.DrawImage(assets.WhitePixel, opBar)

	scoreStr := fmt.Sprintf("SCORE: %d", p.accessor.Score())
	seconds := int(p.accessor.GameTime().Seconds())
	timeStr := fmt.Sprintf("TIME: %02d:%02d", seconds/60, seconds%60)

	text.Draw(screen, scoreStr, assets.UIFont, 10, 25, color.Black)
	text.Draw(screen, timeStr, assets.UIFont, cfg.ScreenWidth-200, 25, color.Black)

	if p.snake.IsAlive {
		p.drawSnake(screen)
	}

	p.drawFood(screen)
	p.drawWalls(screen)
}

func (p *PlayingScene) drawSnake(screen *ebiten.Image) {
	assets := p.accessor.Assets()
	cfg := p.accessor.Config()

	for i, segment := range p.snake.Body {

		var (
			img      *ebiten.Image
			rotation float64
		)

		if i == 0 {
			img = assets.SnakeHead
			rotation = core.DirectionToRotationAngle(p.snake.Direction)

		} else if i == len(p.snake.Body)-1 {
			img = assets.SnakeTail
			direction := core.GetDirection(p.snake.Body[i-1].Position, segment.Position)
			rotation = core.DirectionToRotationAngle(direction)

		} else {
			newDirection := core.GetDirection(p.snake.Body[i-1].Position, segment.Position)
			oldDirection := core.GetDirection(segment.Position, p.snake.Body[i+1].Position)
			if newDirection == oldDirection {
				img = assets.SnakeBody
				rotation = core.DirectionToRotationAngle(oldDirection)
			} else {
				img = assets.SnakeBodyCorner
				rotation = core.CornerToRotationAngle(oldDirection, newDirection)
			}
		}

		if img != nil {
			op := &ebiten.DrawImageOptions{}

			originalWidth := img.Bounds().Dx()
			originalHeight := img.Bounds().Dy()

			var scaleX, scaleY float64
			if originalWidth > 0 {
				scaleX = float64(cfg.TileSize) / float64(originalWidth)
			}
			if originalHeight > 0 {
				scaleY = float64(cfg.TileSize) / float64(originalHeight)
			}

			op.GeoM.Scale(scaleX, scaleY)
			op.GeoM.Translate(-float64(cfg.TileSize)/2, -float64(cfg.TileSize)/2)

			op.GeoM.Rotate(rotation)

			op.GeoM.Translate(float64(segment.X*cfg.TileSize)+float64(cfg.TileSize)/2, float64(segment.Y*cfg.TileSize)+float64(cfg.TileSize)/2+float64(cfg.TopBarHeight))

			screen.DrawImage(img, op)
		}
	}
}

func (p *PlayingScene) drawFood(screen *ebiten.Image) {
	cfg := p.accessor.Config()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.food.X*cfg.TileSize), float64(p.food.Y*cfg.TileSize)+float64(cfg.TopBarHeight))
	img := p.accessor.Assets().Apple
	screen.DrawImage(img, op)
}

func (p *PlayingScene) drawWalls(screen *ebiten.Image) {
	cfg := p.accessor.Config()
	for _, wall := range p.walls {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(wall.X*cfg.TileSize), float64(wall.Y*cfg.TileSize)+float64(cfg.TopBarHeight))
		img := p.accessor.Assets().Wall
		screen.DrawImage(img, op)
	}
}

func (p *PlayingScene) OnEnter() {
	p.accessor.Logger().Info("Entering playing scene", "level", p.level.Name)
	// p.isPaused = false
}
