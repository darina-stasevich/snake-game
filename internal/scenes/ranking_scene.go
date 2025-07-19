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
	RecordsNumber = 20
)

type RankingScene struct {
	accessor GameAccessor

	records   []storage.Record
	loadError error

	nextState core.GameState

	levelName  []rune
	playerName []rune

	isScoreAsc bool
	isTimeAsc  bool

	activeField   string
	cursorVisible bool
	cursorBlink   time.Time

	scoreButton *ui.Button
	timeButton  *ui.Button

	playerNameFieldRect image.Rectangle
	levelNameFieldRect  image.Rectangle
}

func NewRankingScene(accessor GameAccessor) *RankingScene {
	scene := &RankingScene{
		accessor: accessor,
	}

	fieldsY := 120 // Y-координата для полей ввода

	// --- 2. Инициализация полей ввода ---
	fieldWidth := 300
	fieldHeight := 40

	playerFieldX := 40
	scene.playerNameFieldRect = image.Rect(playerFieldX, fieldsY, playerFieldX+fieldWidth, fieldsY+fieldHeight)

	levelFieldX := playerFieldX + fieldWidth + 20
	scene.levelNameFieldRect = image.Rect(levelFieldX, fieldsY, levelFieldX+fieldWidth, fieldsY+fieldHeight)

	colX_Score := 300
	colX_Time := 410
	scoreTextWidth := text.BoundString(accessor.Assets().UIFont, "SCORE").Dx()
	timeTextWidth := text.BoundString(accessor.Assets().UIFont, "TIME").Dx()
	buttonSize := 25
	buttonY := float64(195)

	scene.scoreButton = ui.NewButton(
		float64(colX_Score+scoreTextWidth+5),
		buttonY,
		float64(buttonSize),
		float64(buttonSize),
		"F",
		func() {
			scene.isScoreAsc = !scene.isScoreAsc
			scene.loadRecords()
		},
	)

	scene.timeButton = ui.NewButton(
		float64(colX_Time+timeTextWidth+5),
		buttonY,
		float64(buttonSize),
		float64(buttonSize),
		"F",
		func() {
			scene.isTimeAsc = !scene.isTimeAsc
			scene.loadRecords()
		},
	)

	scene.reset()
	scene.loadRecords()

	return scene
}

func (r *RankingScene) reset() {
	r.nextState = core.BestScoresState
	r.isScoreAsc = false
	r.isTimeAsc = true
	r.loadError = nil
	r.records = make([]storage.Record, 0)
	r.levelName = make([]rune, 0)
	r.playerName = make([]rune, 0)
	r.activeField = ""
}

func (r *RankingScene) loadRecords() {
	repo := r.accessor.Repository()
	filter := storage.NewFilter(string(r.playerName), string(r.levelName), r.isScoreAsc, r.isTimeAsc, RecordsNumber)
	var err error
	r.records, err = repo.GetTopRecords(context.Background(), *filter)
	if err != nil {
		r.accessor.Logger().Error("failed to load records", "error", err)
		r.loadError = err
	} else {
		r.loadError = nil
	}
}

func (r *RankingScene) Draw(screen *ebiten.Image) {

	screen.Fill(color.NRGBA{R: 0x0A, G: 0x19, B: 0x4E, A: 0xff})

	cfg := r.accessor.Config()
	uiFont := r.accessor.Assets().UIFont
	titleFont := r.accessor.Assets().TitleFont

	title := "TOP SCORES"
	titleBounds := text.BoundString(titleFont, title)
	titleX := (cfg.ScreenWidth - titleBounds.Dx()) / 2
	text.Draw(screen, title, titleFont, titleX, 60, color.White)

	if r.loadError != nil {
		errorMsg := fmt.Sprintf("Error: Could not load records.")
		errorBounds := text.BoundString(uiFont, errorMsg)
		errorX := (cfg.ScreenWidth - errorBounds.Dx()) / 2
		text.Draw(screen, errorMsg, uiFont, errorX, cfg.ScreenHeight/2, color.RGBA{R: 255, G: 100, B: 100, A: 255})
		return
	}

	text.Draw(screen, "Player Name:", uiFont, r.playerNameFieldRect.Min.X, r.playerNameFieldRect.Min.Y-10, color.White)
	text.Draw(screen, "Level Name:", uiFont, r.levelNameFieldRect.Min.X, r.levelNameFieldRect.Min.Y-10, color.White)
	// Сами поля
	r.drawInputField(screen, string(r.playerName), r.playerNameFieldRect, "player")
	r.drawInputField(screen, string(r.levelName), r.levelNameFieldRect, "level")

	headerY := 220
	colX_Num := 40
	colX_Player := 90
	colX_Score := 300
	colX_Time := 410
	colX_Level := 510
	colX_Date := 710

	text.Draw(screen, "№", uiFont, colX_Num, headerY, color.White)
	text.Draw(screen, "PLAYER", uiFont, colX_Player, headerY, color.White)
	text.Draw(screen, "SCORE", uiFont, colX_Score, headerY, color.White)
	text.Draw(screen, "TIME", uiFont, colX_Time, headerY, color.White)
	text.Draw(screen, "LEVEL", uiFont, colX_Level, headerY, color.White)
	text.Draw(screen, "DATE", uiFont, colX_Date, headerY, color.White)

	r.scoreButton.Draw(screen, r.accessor.Assets())
	r.timeButton.Draw(screen, r.accessor.Assets())

	if len(r.records) == 0 {
		noRecordsMsg := "No records yet. Be the first!"
		noRecordsBounds := text.BoundString(uiFont, noRecordsMsg)
		noRecordsX := (cfg.ScreenWidth - noRecordsBounds.Dx()) / 2
		text.Draw(screen, noRecordsMsg, uiFont, noRecordsX, headerY+60, color.Gray{Y: 180})
	} else {
		for i, record := range r.records {
			rowY := headerY + (i+1)*30
			// #
			text.Draw(screen, fmt.Sprintf("%d.", i+1), uiFont, colX_Num, rowY, color.White)
			// PLAYER
			text.Draw(screen, record.PlayerName, uiFont, colX_Player, rowY, color.White)
			// SCORE
			text.Draw(screen, fmt.Sprintf("%d", record.Score), uiFont, colX_Score, rowY, color.White)
			// TIME (в формате ММ:СС)
			gameTime := time.Unix(0, 0).Add(record.Time)
			text.Draw(screen, gameTime.Format("04:05"), uiFont, colX_Time, rowY, color.White)
			// LEVEL
			text.Draw(screen, record.LevelName, uiFont, colX_Level, rowY, color.White)
			// DATE (в формате ГГГГ-ММ-ДД)
			text.Draw(screen, record.CreatedAt.Format("2006-01-02"), uiFont, colX_Date, rowY, color.White)
		}
	}

	exitMsg := "Press ESC to return to menu"
	exitBounds := text.BoundString(uiFont, exitMsg)
	text.Draw(screen, exitMsg, uiFont, (cfg.ScreenWidth-exitBounds.Dx())/2, cfg.ScreenHeight-40, color.White)
}

func (r *RankingScene) drawInputField(screen *ebiten.Image, content string, rect image.Rectangle, fieldName string) {
	borderColor := color.Gray{Y: 100}

	ui.DrawRectangle(screen, r.accessor.Assets(), float64(rect.Min.X-2), float64(rect.Min.Y-2), float64(rect.Dx()+4), float64(rect.Dy()+4), borderColor)
	ui.DrawRectangle(screen, r.accessor.Assets(), float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()), color.Black)
	text.Draw(screen, content, r.accessor.Assets().UIFont, rect.Min.X+5, rect.Min.Y+28, color.White)

	if r.activeField == fieldName && r.cursorVisible {
		textBounds := text.BoundString(r.accessor.Assets().UIFont, content)
		cursorX := float64(rect.Min.X+5) + float64(textBounds.Dx()) + 2
		ui.DrawRectangle(screen, r.accessor.Assets(), cursorX, float64(rect.Min.Y+8), 2, float64(rect.Dy()-16), color.White)
	}
}

func (r *RankingScene) Update() (core.GameState, error) {
	r.accessor.Logger().Debug("i'm updating", "player_name", string(r.playerName))
	r.setActiveField()
	switch r.activeField {
	case "level":
		inputChars := ebiten.AppendInputChars(r.levelName)
		if len(inputChars) > MaxLevelName {
			inputChars = inputChars[0:MaxLevelName]
		}
		if len(inputChars) != len(r.levelName) {
			r.levelName = inputChars
			r.loadRecords()
		} else {
			r.levelName = inputChars
		}
	case "player":
		inputChars := ebiten.AppendInputChars(r.playerName)
		r.accessor.Logger().Debug("i'm updting and i'm active field", "inputChars", string(inputChars))
		if len(inputChars) > MaxPlayerName {
			inputChars = inputChars[0:MaxPlayerName]
		}
		if len(inputChars) != len(r.playerName) {
			r.playerName = inputChars
			r.loadRecords()
		} else {
			r.playerName = inputChars
		}
		r.accessor.Logger().Debug("i'm updating and after check", "player_name", string(r.playerName))
		r.scoreButton.Update()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		switch r.activeField {
		case "level":
			if len(r.levelName) > 0 {
				r.levelName = r.levelName[0 : len(r.levelName)-1]
				r.loadRecords()
			}
		case "player":
			if len(r.playerName) > 0 {
				r.playerName = r.playerName[0 : len(r.playerName)-1]
				r.loadRecords()
			}
		}
	}
	if time.Since(r.cursorBlink) > time.Millisecond*500 {
		r.cursorVisible = !r.cursorVisible
		r.cursorBlink = time.Now()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return core.MainMenuState, nil
	}

	r.scoreButton.Update()
	r.timeButton.Update()
	return core.BestScoresState, nil
}

func (r *RankingScene) setActiveField() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		mousePoint := image.Pt(cursorX, cursorY)
		if mousePoint.In(r.levelNameFieldRect) {
			r.activeField = "level"
		} else if mousePoint.In(r.playerNameFieldRect) {
			r.activeField = "player"
		} else {
			r.activeField = ""
		}
	}
}

func (r *RankingScene) OnEnter() {
	r.accessor.Logger().Info("entering ranking scene")
	r.reset()
	r.loadRecords()
}
