package scenes

import (
	"context"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
	"snake-game/internal/core"
	"snake-game/internal/storage"
	"snake-game/internal/ui"
	"time"
)

const (
	MaxPlayerName = 13
)

type GameOverScene struct {
	accessor GameAccessor

	level *core.Level

	nextState core.GameState

	newGameButton   *ui.Button
	mainMenuButton  *ui.Button
	saveScoreButton *ui.Button
	nameFieldRect   image.Rectangle

	isRecordSaved bool
	playerName    []rune
}

func NewGameOverScene(accessor GameAccessor, level *core.Level) *GameOverScene {
	scene := &GameOverScene{
		accessor:  accessor,
		level:     level,
		nextState: core.GameOverState,
	}

	cfg := scene.accessor.Config()
	centerX := float64(cfg.ScreenWidth) / 2

	newGameButton := ui.NewButton(
		centerX-120,
		float64(cfg.ScreenHeight/2)+90,
		240,
		50,
		"NEW GAME",
		func() {
			accessor.StartGame(level)
			scene.nextState = core.GamePlayingState
		},
	)

	scene.newGameButton = newGameButton

	mainMenuButton := ui.NewButton(
		centerX-120,
		float64(cfg.ScreenHeight/2)+145,
		240,
		50,
		"MAIN MENU",
		func() {
			scene.nextState = core.MainMenuState
		})

	scene.mainMenuButton = mainMenuButton

	saveButton := ui.NewButton(
		centerX+130,
		float64(cfg.ScreenHeight/2)+38,
		40,
		40,
		"S",
		func() {
			if scene.isRecordSaved == true {
				return
			}
			record := storage.NewRecord(string(scene.playerName), scene.accessor.Score(), scene.accessor.GameTime(), scene.level.Name, time.Now())
			err := scene.accessor.Repository().SaveRecord(context.Background(), record)
			if err != nil {
				scene.accessor.Logger().Error("failed to save record", "error", err)
			} else {
				scene.isRecordSaved = true
			}
		})

	scene.saveScoreButton = saveButton

	inputFieldWidth := 240.0
	inputFieldHeight := 40.0
	inputX := centerX - 120
	inputY := float64(cfg.ScreenHeight/2) + 38
	scene.nameFieldRect = image.Rect(int(inputX), int(inputY), int(inputX+inputFieldWidth), int(inputY+inputFieldHeight))

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

	s.drawInputField(screen)

	s.newGameButton.Draw(screen, assets)
	s.mainMenuButton.Draw(screen, assets)
	s.saveScoreButton.Draw(screen, assets)
}

func (s *GameOverScene) drawInputField(screen *ebiten.Image) {
	assets := s.accessor.Assets()
	rect := s.nameFieldRect

	ui.DrawRectangle(screen, assets, float64(rect.Min.X-2), float64(rect.Min.Y-2), float64(rect.Dx()+4), float64(rect.Dy()+4), color.Gray{Y: 100})
	ui.DrawRectangle(screen, assets, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()), color.Black)
	text.Draw(screen, string(s.playerName), assets.UIFont, rect.Min.X+15, rect.Min.Y+28, color.White)
}

func (s *GameOverScene) Update() (core.GameState, error) {
	s.handleInput()
	s.newGameButton.Update()
	s.mainMenuButton.Update()
	s.saveScoreButton.Update()
	return s.nextState, nil
}

func (s *GameOverScene) handleInput() {
	s.playerName = ebiten.AppendInputChars(s.playerName)
	if len(s.playerName) > MaxPlayerName {
		s.playerName = s.playerName[0:MaxPlayerName]
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(s.playerName) > 0 {
			s.playerName = s.playerName[0 : len(s.playerName)-1]
		}
	}
}

func (s *GameOverScene) OnEnter() {
	s.nextState = core.GameOverState
	s.isRecordSaved = false
}
