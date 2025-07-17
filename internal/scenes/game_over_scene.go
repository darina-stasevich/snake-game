package scenes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"snake-game/internal/core"
	"snake-game/internal/ui"
)

type GameOverScene struct {
	accessor GameAccessor

	level *core.Level

	nextState     core.GameState
	newGameButton *ui.Button
}

func NewGameOverScene(accessor GameAccessor, level *core.Level) *GameOverScene {
	scene := &GameOverScene{
		accessor:  accessor,
		level:     level,
		nextState: core.GameOverState,
	}

	cfg := scene.accessor.Config()
	centerX := float64(cfg.ScreenWidth) / 2

	button := ui.NewButton(
		centerX-120,
		float64(cfg.ScreenHeight/2)+80,
		240,
		50,
		"NEW GAME",
		func() {
			accessor.StartGame(level)
			scene.nextState = core.GamePlayingState
		},
	)

	scene.newGameButton = button

	return scene
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	cfg := s.accessor.Config()
	assets := s.accessor.Assets()

	opOverlay := &ebiten.DrawImageOptions{}
	opOverlay.GeoM.Scale(float64(cfg.ScreenWidth), float64(cfg.ScreenHeight))

	opOverlay.ColorScale.Scale(0, 0, 0, 0.66)
	screen.DrawImage(assets.WhitePixel, opOverlay)

	uiFont := assets.UIFont
	centerX := cfg.ScreenWidth / 2

	gameOverText := "GAME OVER"
	gameOverBounds := text.BoundString(uiFont, gameOverText)
	gameOverX := centerX - gameOverBounds.Dx()/2
	gameOverY := cfg.WindowHeight()/2 - 80
	text.Draw(screen, gameOverText, uiFont, gameOverX, gameOverY, color.White)

	scoreStr := fmt.Sprintf("FINAL SCORE: %d", s.accessor.Score())
	seconds := int(s.accessor.GameTime().Seconds())
	timeStr := fmt.Sprintf("TIME: %02d:%02d", seconds/60, seconds%60)

	scoreBounds := text.BoundString(uiFont, scoreStr)
	scoreX := centerX - scoreBounds.Dx()/2
	scoreY := gameOverY + 40
	text.Draw(screen, scoreStr, uiFont, scoreX, scoreY, color.White)

	timeBounds := text.BoundString(uiFont, timeStr)
	timeX := centerX - timeBounds.Dx()/2
	timeY := scoreY + 25
	text.Draw(screen, timeStr, uiFont, timeX, timeY, color.White)

	inputFieldWidth := 240.0
	inputFieldHeight := 40.0
	inputX := float64(centerX) - inputFieldWidth/2
	inputY := float64(timeY) + 40

	opBorder := &ebiten.DrawImageOptions{}
	opBorder.GeoM.Scale(inputFieldWidth+4, inputFieldHeight+4)
	opBorder.GeoM.Translate(inputX-2, inputY-2)
	opBorder.ColorScale.Scale(0.5, 0.5, 0.5, 1) // Серый цвет
	screen.DrawImage(assets.WhitePixel, opBorder)

	opField := &ebiten.DrawImageOptions{}
	opField.GeoM.Scale(inputFieldWidth, inputFieldHeight)
	opField.GeoM.Translate(inputX, inputY)
	opField.ColorScale.Scale(0.1, 0.1, 0.1, 1) // Темно-серый
	screen.DrawImage(assets.WhitePixel, opField)

	s.newGameButton.Draw(screen, assets)
}

func (s *GameOverScene) Update() (core.GameState, error) {
	s.newGameButton.Update()
	return s.nextState, nil
}

func (s *GameOverScene) OnEnter() {
	s.nextState = core.GameOverState
}
