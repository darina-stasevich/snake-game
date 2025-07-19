package scenes

import (
	"context"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
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

	levelName  string
	playerName string

	isScoreAsc bool
	isTimeAsc  bool

	scoreButton *ui.Button
	timeButton  *ui.Button
}

func NewRankingScene(accessor GameAccessor) *RankingScene {
	scene := &RankingScene{
		accessor: accessor,
	}

	scene.reset()
	scene.loadRecords()

	return scene
}

func (r *RankingScene) reset() {
	r.nextState = core.BestScoresState
	r.isScoreAsc = false
	r.isTimeAsc = true
	r.loadError = nil
}

func (r *RankingScene) loadRecords() {
	repo := r.accessor.Repository()
	filter := storage.NewFilter(r.playerName, r.levelName, r.isScoreAsc, r.isTimeAsc, RecordsNumber)
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
	// 1. Рисуем фон
	screen.Fill(color.NRGBA{R: 0x00, G: 0x40, B: 0x00, A: 0xff}) // Привычный темно-зеленый

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

	headerY := 120
	colX_Num := 40
	colX_Player := 90
	colX_Score := 250
	colX_Time := 350
	colX_Level := 450
	colX_Date := 650

	text.Draw(screen, "№", uiFont, colX_Num, headerY, color.White)
	text.Draw(screen, "PLAYER", uiFont, colX_Player, headerY, color.White)
	text.Draw(screen, "SCORE", uiFont, colX_Score, headerY, color.White)
	text.Draw(screen, "TIME", uiFont, colX_Time, headerY, color.White)
	text.Draw(screen, "LEVEL", uiFont, colX_Level, headerY, color.White)
	text.Draw(screen, "DATE", uiFont, colX_Date, headerY, color.White)

	// 5. Рисуем строки с рекордами
	if len(r.records) == 0 {
		// Сообщение, если рекордов еще нет
		noRecordsMsg := "No records yet. Be the first!"
		noRecordsBounds := text.BoundString(uiFont, noRecordsMsg)
		noRecordsX := (cfg.ScreenWidth - noRecordsBounds.Dx()) / 2
		text.Draw(screen, noRecordsMsg, uiFont, noRecordsX, headerY+60, color.Gray{Y: 180})
	} else {
		for i, record := range r.records {
			rowY := headerY + (i+1)*30 // Смещаем каждую строку на 30 пикселей вниз

			// №
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

	// 6. Подсказка для выхода (пока не реализовано, но пусть будет)
	exitMsg := "Press ESC to return to menu"
	exitBounds := text.BoundString(uiFont, exitMsg)
	text.Draw(screen, exitMsg, uiFont, (cfg.ScreenWidth-exitBounds.Dx())/2, cfg.ScreenHeight-40, color.White)
}

func (r *RankingScene) Update() (core.GameState, error) {
	// TODO setActiveField
	// TODO handleInput
	// r.scoreButton.Update()
	// r.timeButton.Update()
	return core.BestScoresState, nil
}

func (r *RankingScene) OnEnter() {
	r.accessor.Logger().Info("entering ranking scene")
	r.reset()
	r.loadRecords()
}
