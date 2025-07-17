package scenes

import (
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"snake-game/internal/core"
	"snake-game/internal/ui"
	"strings"
)

type MainMenuScene struct {
	accessor GameAccessor

	levelNames   []string
	currentLevel int

	nextState         core.GameState
	newGameButton     *ui.Button
	createLevelButton *ui.Button
	rankingButton     *ui.Button
	quitButton        *ui.Button
}

func NewMainMenuScene(accessor GameAccessor) *MainMenuScene {
	scene := &MainMenuScene{
		accessor:  accessor,
		nextState: core.MainMenuState,
	}

	cfg := scene.accessor.Config()
	centerX := float64(cfg.ScreenWidth) / 2
	startY := float64(cfg.ScreenHeight/2) + 80
	buttonWidth := float64(240)
	buttonHeight := float64(50)
	buttonSpacing := buttonHeight + 10

	newGameButton := ui.NewButton(centerX-120, startY, buttonWidth, buttonHeight, "NEW GAME", scene.newGame)
	createLevelButton := ui.NewButton(centerX-120, startY+buttonSpacing, buttonWidth, buttonHeight, "CREATE LEVEL", scene.createLevel)
	rankingButton := ui.NewButton(centerX-120, startY+2*buttonSpacing, buttonWidth, buttonHeight, "RANKING", scene.ranking)
	quitButton := ui.NewButton(centerX-120, startY+3*buttonSpacing, buttonWidth, buttonHeight, "QUIT",
		func() {
			os.Exit(0)
		},
	)

	scene.newGameButton = newGameButton
	scene.createLevelButton = createLevelButton
	scene.rankingButton = rankingButton
	scene.quitButton = quitButton

	return scene
}

func (s *MainMenuScene) Draw(screen *ebiten.Image) {
	cfg := s.accessor.Config()
	assets := s.accessor.Assets()
	centerX := cfg.ScreenWidth / 2

	screen.Fill(color.RGBA{R: 20, G: 20, B: 40, A: 255})

	titleFont := assets.TitleFont
	titleText := "SNAKE GAME"
	titleBounds := text.BoundString(titleFont, titleText)
	titleX := centerX - titleBounds.Dx()/2
	text.Draw(screen, titleText, titleFont, titleX, 50, color.White)

	s.drawLevelSelector(screen)

	s.newGameButton.Draw(screen, assets)
	s.createLevelButton.Draw(screen, assets)
	s.rankingButton.Draw(screen, assets)
	s.quitButton.Draw(screen, assets)
}

func (s *MainMenuScene) drawLevelSelector(screen *ebiten.Image) {
	assets := s.accessor.Assets()
	cfg := s.accessor.Config()

	centerX := cfg.ScreenWidth / 2

	labelText := "Select level:"
	labelBounds := text.BoundString(assets.UIFont, labelText)
	labelX := centerX - labelBounds.Dx()/2 - 60
	labelY := float64(cfg.ScreenHeight/2) - 80
	text.Draw(screen, labelText, assets.UIFont, labelX, int(labelY), color.White)

	fieldWidth := 240.0
	fieldHeight := 40.0
	fieldX := float64(labelX) + float64(labelBounds.Dx()) + 10
	fieldY := labelY - fieldHeight/2 - 5

	ui.DrawRectangle(screen, assets, fieldX-2, fieldY-2, fieldWidth+4, fieldHeight+4, color.Gray{Y: 128})
	ui.DrawRectangle(screen, assets, fieldX, fieldY, fieldWidth, fieldHeight, color.Black)

	var levelName string
	if len(s.levelNames) > 0 {
		levelName = strings.TrimSuffix(s.levelNames[s.currentLevel], ".json")
	} else {
		levelName = "Levels not found"
	}
	levelBounds := text.BoundString(assets.UIFont, levelName)
	levelTextX := fieldX + (fieldWidth-float64(levelBounds.Dx()))/2
	levelTextY := fieldY + (fieldHeight+float64(levelBounds.Dy()))/2
	text.Draw(screen, levelName, assets.UIFont, int(levelTextX), int(levelTextY), color.White)
}

func (s *MainMenuScene) Update() (core.GameState, error) {

	s.newGameButton.Update()
	s.createLevelButton.Update()
	s.rankingButton.Update()
	s.quitButton.Update()

	s.handleInput()

	return s.nextState, nil

}

func (s *MainMenuScene) handleInput() {
	if len(s.levelNames) == 0 {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.currentLevel--
		if s.currentLevel < 0 {
			s.currentLevel = len(s.levelNames) - 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		s.currentLevel = (s.currentLevel + 1) % len(s.levelNames)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.newGame()
	}
}

func (s *MainMenuScene) OnEnter() {
	s.accessor.Logger().Info("Entering main menu, scanning for levels...")

	s.levelNames = []string{}
	s.currentLevel = 0

	if err := os.MkdirAll("levels", 0755); err != nil {
		s.accessor.Logger().Error("failed to create levels directory", "error", err)
		return
	}

	err := filepath.WalkDir("levels", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			s.levelNames = append(s.levelNames, d.Name())
		}
		return nil
	})

	if err != nil {
		s.accessor.Logger().Error("failed to scan for levels", "error", err)
	}
}

func (s *MainMenuScene) newGame() {
	if len(s.levelNames) == 0 {
		s.accessor.Logger().Warn("no levels were found")
		return
	}

	levelFileName := s.levelNames[s.currentLevel]
	levelPath := path.Join("levels", levelFileName)

	s.accessor.Logger().Info("loading level", "path", levelPath)
	data, err := os.ReadFile(levelPath)
	if err != nil {
		s.accessor.Logger().Error("failed to read level file", "path", levelPath, "error", err)
		return
	}

	var level core.Level
	if err := json.Unmarshal(data, &level); err != nil {
		s.accessor.Logger().Error("failed to parse level json", "path", levelPath, "error", err)
		return
	}

	if level.GridHeight > s.accessor.Config().ScreenHeight/s.accessor.Config().TileSize ||
		level.GridWidth > s.accessor.Config().ScreenWidth/s.accessor.Config().TileSize {
		s.accessor.Logger().Warn("unappropriated level size",
			"grid_width", level.GridWidth,
			"grid_height", level.GridHeight,
			"cfg_max_width", s.accessor.Config().ScreenHeight/s.accessor.Config().TileSize,
			"cfg_max_height", s.accessor.Config().ScreenWidth/s.accessor.Config().TileSize)
		return
	}

	s.nextState = core.GamePlayingState
	s.accessor.StartGame(&level)
}

func (s *MainMenuScene) createLevel() {
	s.accessor.Logger().Info("go to createLevelScene (not implement)")
	// s.nextState = core.LevelCreateState
}

func (s *MainMenuScene) ranking() {
	s.accessor.Logger().Info("go to rankingScene (not implement)")
	// s.nextState = core.BestScoresState
}
