package scenes

import (
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"snake-game/internal/core"
	"snake-game/internal/ui"
	"strconv"
	"strings"
	"time"
)

const (
	MinimalWidth int = 3
	MinimalHeight
	MaxLevelName = 30
)

type CreateLevelScene struct {
	accessor GameAccessor

	walls  map[core.Position]bool
	width  int
	height int

	maximalWidth  int
	maximalHeight int

	LevelName []rune
	widthStr  []rune
	heightStr []rune

	isNameValid   bool
	isWidthValid  bool
	isHeightValid bool

	saveButton  *ui.Button
	resetButton *ui.Button

	nameFieldRect   image.Rectangle
	widthFieldRect  image.Rectangle
	heightFieldRect image.Rectangle

	activeField   string
	cursorVisible bool
	cursorBlink   time.Time

	nextState core.GameState
}

func NewCreateLevelScene(accessor GameAccessor) *CreateLevelScene {
	scene := &CreateLevelScene{
		accessor: accessor,
	}

	scene.layoutUI()
	scene.reset()

	return scene
}

func (c *CreateLevelScene) layoutUI() {
	cfg := c.accessor.Config()
	topBarY := float64(cfg.TopBarHeight)
	fieldHeight := min(40.0, topBarY-2)

	centerY := topBarY / 2
	currentX := 10.0

	resetButtonWidth := 100.0
	c.resetButton = ui.NewButton(currentX, centerY-20, resetButtonWidth, 40, "Reset", c.reset)
	currentX += resetButtonWidth + 20

	saveButtonWidth := 100.0
	c.saveButton = ui.NewButton(currentX, centerY-20, saveButtonWidth, 40, "Save", c.save)
	currentX += saveButtonWidth + 40

	currentX += 80
	nameFieldWidth := 200.0
	c.nameFieldRect = image.Rect(int(currentX), 0, int(currentX+nameFieldWidth), int(fieldHeight))
	currentX += nameFieldWidth + 20

	currentX += 35
	widthFieldWidth := 60.0
	c.widthFieldRect = image.Rect(int(currentX), 0, int(currentX+widthFieldWidth), int(fieldHeight))
	currentX += widthFieldWidth + 20

	currentX += 35
	heightFieldWidth := 60.0
	c.heightFieldRect = image.Rect(int(currentX), 0, int(currentX+heightFieldWidth), int(fieldHeight))

}

func (c *CreateLevelScene) reset() {
	c.width = MinimalWidth
	c.height = MinimalHeight
	c.maximalWidth = c.accessor.Config().ScreenWidth / c.accessor.Config().TileSize
	c.maximalHeight = c.accessor.Config().ScreenHeight / c.accessor.Config().TileSize

	c.LevelName = []rune("new_level")
	c.widthStr = []rune(strconv.Itoa(c.width))
	c.heightStr = []rune(strconv.Itoa(c.height))
	c.isNameValid = true
	c.isWidthValid = true
	c.isHeightValid = true
	c.activeField = ""

	c.walls = make(map[core.Position]bool)

	c.nextState = core.LevelCreateState
}

func (c *CreateLevelScene) save() {
	if !c.isNameValid || !c.isWidthValid || !c.isHeightValid {
		c.accessor.Logger().Warn("Save aborted: invalid data in fields")
		return
	}

	walls := c.wallsInSlice()
	level := core.NewLevel(string(c.LevelName), c.width, c.height, walls)
	levelsDir := "levels"

	levelJson, err := json.Marshal(level)
	if err != nil {
		c.accessor.Logger().Error("failed to marshal level to json", "error", err)
		c.nextState = core.MainMenuState
		return
	}

	if err := os.MkdirAll(levelsDir, 0755); err != nil {
		c.accessor.Logger().Error("Failed to create levels directory", "path", levelsDir, "error", err)
		return
	}

	finalPath, err := c.findAvailableFilename(levelsDir, level.Name)
	if err != nil {
		c.accessor.Logger().Error("Failed to find available filename", "error", err)
		return
	}

	file, err := os.Create(finalPath)
	if err != nil {
		c.accessor.Logger().Error("Failed to create file for writing", "path", finalPath, "error", err)
		return
	}
	defer file.Close()

	err = os.WriteFile(finalPath, levelJson, 0644)
	if err != nil {
		c.accessor.Logger().Error("Failed to write level file", "path", finalPath, "error", err)
		return
	}

	c.accessor.Logger().Info("Level saved", "name", string(c.LevelName), "w", c.width, "h", c.height)
	c.nextState = core.MainMenuState
}

func (c *CreateLevelScene) findAvailableFilename(dir, levelName string) (string, error) {

	basepath := filepath.Join(dir, levelName+".json")

	if _, err := os.Stat(basepath); os.IsNotExist(err) {
		return basepath, nil
	}
	i := 1
	for {
		newPath := filepath.Join(dir, levelName+"_"+strconv.Itoa(i)+".json")
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath, nil
		}
		i++
	}
}

func (c *CreateLevelScene) wallsInSlice() []core.Wall {
	walls := make([]core.Wall, 0)
	for pos, val := range c.walls {
		if val == true {
			walls = append(walls, *core.NewWall(pos.X, pos.Y))
		}
	}
	return walls
}

func (c *CreateLevelScene) Draw(screen *ebiten.Image) {
	cfg := c.accessor.Config()
	assets := c.accessor.Assets()

	screen.Fill(color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xff})

	topBarImg := assets.WhitePixel
	opBar := &ebiten.DrawImageOptions{}
	opBar.GeoM.Scale(float64(cfg.ScreenWidth), float64(cfg.TopBarHeight))
	opBar.ColorScale.ScaleWithColor(color.Gray{Y: 100})
	screen.DrawImage(topBarImg, opBar)

	c.resetButton.Draw(screen, assets)
	c.saveButton.Draw(screen, assets)

	uiFont := assets.UIFont
	textColor := color.White
	textYOffset := c.nameFieldRect.Min.Y + 28

	text.Draw(screen, "Name:", uiFont, c.nameFieldRect.Min.X-80, textYOffset, textColor)
	c.drawInputField(screen, string(c.LevelName), c.nameFieldRect, "name", c.isNameValid)

	text.Draw(screen, "W:", uiFont, c.widthFieldRect.Min.X-35, textYOffset, textColor)
	c.drawInputField(screen, string(c.widthStr), c.widthFieldRect, "width", c.isWidthValid)

	text.Draw(screen, "H:", uiFont, c.heightFieldRect.Min.X-35, textYOffset, textColor)
	c.drawInputField(screen, string(c.heightStr), c.heightFieldRect, "height", c.isHeightValid)

	gridOriginX := 0.0
	gridOriginY := float64(cfg.TopBarHeight)
	gridWidthPixels := float64(c.width * cfg.TileSize)
	gridHeightPixels := float64(c.height * cfg.TileSize)

	// Фон для активной области сетки
	ui.DrawRectangle(screen, assets, gridOriginX, gridOriginY, gridWidthPixels, gridHeightPixels, color.NRGBA{38, 38, 48, 255})

	// Отрисовка стен
	for pos, isSet := range c.walls {
		if !isSet {
			continue
		}
		ui.DrawRectangle(screen, assets,
			gridOriginX+float64(pos.X*cfg.TileSize),
			gridOriginY+float64(pos.Y*cfg.TileSize),
			float64(cfg.TileSize), float64(cfg.TileSize),
			color.Gray{Y: 120},
		)
	}

}

func (c *CreateLevelScene) drawInputField(screen *ebiten.Image, content string, rect image.Rectangle, fieldName string, isValid bool) {
	assets := c.accessor.Assets()

	var borderColor color.Color

	if isValid {
		borderColor = color.Gray{Y: 100}
	} else {
		borderColor = color.RGBA{R: 200, G: 0, B: 0, A: 255}
	}

	ui.DrawRectangle(
		screen, assets,
		float64(rect.Min.X-2), float64(rect.Min.Y-2), // Координаты (x, y)
		float64(rect.Dx()+4), float64(rect.Dy()+4), // Размеры (width, height)
		borderColor, // Цвет
	)

	ui.DrawRectangle(
		screen, assets,
		float64(rect.Min.X), float64(rect.Min.Y),
		float64(rect.Dx()), float64(rect.Dy()),
		color.Black,
	)

	text.Draw(screen, content, assets.UIFont, rect.Min.X+5, rect.Min.Y+28, color.White)

	if c.activeField == fieldName && c.cursorVisible {
		textBounds := text.BoundString(assets.UIFont, content)
		cursorX := float64(rect.Min.X+5) + float64(textBounds.Dx()) + 2

		ui.DrawRectangle(
			screen, assets,
			cursorX, float64(rect.Min.Y+8),
			2, float64(rect.Dy()-16),
			color.White,
		)
	}
}

func (c *CreateLevelScene) Update() (core.GameState, error) {
	c.accessor.Logger().Info("updating create scene", "current_active_field", c.activeField, "name_valid", c.isNameValid, "w_valid", c.isWidthValid, "h_valid", c.isHeightValid)
	c.handleInput()

	c.saveButton.Update()
	c.resetButton.Update()

	if time.Since(c.cursorBlink) > time.Millisecond*500 {
		c.cursorVisible = !c.cursorVisible
		c.cursorBlink = time.Now()
	}

	return c.nextState, nil
}

func (c *CreateLevelScene) handleInput() {
	c.setActiveField()
	switch c.activeField {
	case "field":
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			cursorX, cursorY := ebiten.CursorPosition()
			tileX := cursorX / c.accessor.Config().TileSize
			tileY := cursorY / c.accessor.Config().TileSize
			c.walls[core.Position{X: tileX, Y: tileY}] = !c.walls[core.Position{X: tileX, Y: tileY}]
		}
	case "name":
		inputChars := ebiten.AppendInputChars(c.LevelName)
		if len(inputChars) > MaxLevelName {
			inputChars = inputChars[0:MaxLevelName]
		}
		c.LevelName = inputChars
	case "height":
		c.heightStr = ebiten.AppendInputChars(c.heightStr)
	case "width":
		c.widthStr = ebiten.AppendInputChars(c.widthStr)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		switch c.activeField {
		case "name":
			if len(c.LevelName) > 0 {
				c.LevelName = c.LevelName[0 : len(c.LevelName)-1]
			}
		case "height":
			if len(c.heightStr) > 0 {
				c.heightStr = c.heightStr[0 : len(c.heightStr)-1]
			}
		case "width":
			if len(c.widthStr) > 0 {
				c.widthStr = c.widthStr[0 : len(c.widthStr)-1]
			}
		}
	}
	c.validateInputs()
}

func (c *CreateLevelScene) validateInputs() {
	w, err := strconv.Atoi(string(c.widthStr))
	if err == nil {
		if w < MinimalWidth {
			c.isWidthValid = false
		} else if w > c.maximalWidth {
			c.isWidthValid = false
		} else {
			c.width = w
			c.isWidthValid = true
			c.widthStr = []rune(strconv.Itoa(c.width))
		}
	} else {
		c.isWidthValid = false
	}

	h, err := strconv.Atoi(string(c.heightStr))
	if err == nil {
		if h < MinimalHeight {
			c.isHeightValid = false
		} else if h > c.maximalHeight {
			c.isHeightValid = false
		} else {
			c.height = h
			c.isHeightValid = true
			c.heightStr = []rune(strconv.Itoa(c.height))
		}
	} else {
		c.isHeightValid = false
	}

	name := string(c.LevelName)
	invalidChars := "/\\:*?\"<>| "
	if len(name) == 0 || strings.ContainsAny(name, invalidChars) {
		c.isNameValid = false
	} else {
		c.isNameValid = true
	}

}

func (c *CreateLevelScene) setActiveField() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		mousePoint := image.Pt(cursorX, cursorY)
		if cursorX >= 0 && cursorX < c.accessor.Config().TileSize*c.width &&
			cursorY >= c.accessor.Config().TopBarHeight && cursorY < c.accessor.Config().TopBarHeight+c.accessor.Config().TileSize*c.height {
			c.activeField = "field"
		} else if mousePoint.In(c.nameFieldRect) {
			c.activeField = "name"
		} else if mousePoint.In(c.heightFieldRect) {
			c.activeField = "height"
		} else if mousePoint.In(c.widthFieldRect) {
			c.activeField = "width"
		} else {
			c.activeField = ""
		}

	}
}

func (c *CreateLevelScene) OnEnter() {
	c.reset()
}
