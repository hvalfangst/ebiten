package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	_ "image/png"
	"log"
	"os"
	"rattle/zones"
	"strconv"
	"strings"
)

type Game struct {
	currentZone     *zones.Zone
	currentPos      Position
	movementHandled bool
}

type Position struct {
	x              int
	y              int
	justChanged    bool
	previousX      int
	previousY      int
	previousChange bool
}

var (
	backgroundImage, playerImage  *ebiten.Image
	game                          *Game
	playerSprite                  *ebiten.Image
	spriteImages                  []*ebiten.Image
	currentSpriteFrameIndexPlayer int
)

func init() {
	InitializeGame()
}

func InitializeGame() {
	game = &Game{
		&zones.Zone{},
		Position{
			448,
			576,
			false,
			0,
			0,
			false},
		false,
	}
	backgroundImage = createBackgroundImage("assets/areas/reindal_960x640_32.png", game)
	playerSprite = createImage("assets/models/avatars_576x576_12x8.png", game)
	spriteImages = createSpriteImages(game, playerSprite)
	currentSpriteFrameIndexPlayer = 9

}

func (g *Game) Update() error {
	g.HandleKeyPress()
	if g.currentPos.justChanged {
		if g.currentZone.Name == "reindal" && g.currentPos.x >= 600 && g.currentPos.x < 620 && g.currentPos.y > 380 && g.currentPos.y < 410 {

			g.currentPos.x = 372
			g.currentPos.y = 416
			currentSpriteFrameIndexPlayer = 45
			backgroundImage = createBackgroundImage("assets/areas/koala_960x640_32.png", game)
		} else if g.currentZone.Name == "koala" && g.currentPos.x == 372 && g.currentPos.y == 448 {
			backgroundImage = createBackgroundImage("assets/areas/reindal_960x640_32.png", game)
			g.currentPos.x = 632
			g.currentPos.y = 400
			currentSpriteFrameIndexPlayer = 9
		} else {
			fmt.Printf("Zone: %s, Pos X: %d, Pos Y: %d\n", g.currentZone.Name, g.currentPos.x, g.currentPos.y)
		}

	}
	return nil
}

func (g *Game) HandleKeyPress() {
	if !g.movementHandled {

		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.currentPos.Increment(0, -1)
			g.movementHandled = true

			// Alternate between sprite frames 45, 46 and 47
			if currentSpriteFrameIndexPlayer == 45 {
				currentSpriteFrameIndexPlayer = 46
			} else if currentSpriteFrameIndexPlayer == 46 {
				currentSpriteFrameIndexPlayer = 47
			} else {
				currentSpriteFrameIndexPlayer = 45
			}

		} else if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.currentPos.Increment(0, 1)
			g.movementHandled = true

			// Alternate between sprite frames 9, 10 and 11
			if currentSpriteFrameIndexPlayer == 9 {
				currentSpriteFrameIndexPlayer = 10
			} else if currentSpriteFrameIndexPlayer == 10 {
				currentSpriteFrameIndexPlayer = 11
			} else {
				currentSpriteFrameIndexPlayer = 9
			}

		} else if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.currentPos.Increment(-1, 0)
			g.movementHandled = true

			// Alternate between sprite frames 21, 22 and 23
			if currentSpriteFrameIndexPlayer == 21 {
				currentSpriteFrameIndexPlayer = 22
			} else if currentSpriteFrameIndexPlayer == 22 {
				currentSpriteFrameIndexPlayer = 23
			} else {
				currentSpriteFrameIndexPlayer = 21
			}

		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.currentPos.Increment(1, 0)
			g.movementHandled = true

			// Alternate between sprite frames 33, 34 and 35
			if currentSpriteFrameIndexPlayer == 33 {
				currentSpriteFrameIndexPlayer = 34
			} else if currentSpriteFrameIndexPlayer == 34 {
				currentSpriteFrameIndexPlayer = 35
			} else {
				currentSpriteFrameIndexPlayer = 33
			}
		}
	}

	if !ebiten.IsKeyPressed(ebiten.KeyW) && !ebiten.IsKeyPressed(ebiten.KeyS) &&
		!ebiten.IsKeyPressed(ebiten.KeyA) && !ebiten.IsKeyPressed(ebiten.KeyD) {
		g.movementHandled = false
	}

	g.currentPos.justChanged = g.movementHandled
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw the background image
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(backgroundImage, op)

	op.GeoM.Translate(float64(g.currentPos.x), float64(g.currentPos.y))

	// Draw the current sprite image
	screen.DrawImage(spriteImages[currentSpriteFrameIndexPlayer], op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Zone: [%s], Coordinates: (%d, %d), Sprites: %d", g.currentZone.Name, g.currentPos.x, g.currentPos.y, g.currentZone.TotalSprites))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return game.currentZone.ScreenWidth, game.currentZone.ScreenHeight
}

func createSpriteImages(game *Game, playerSprite *ebiten.Image) []*ebiten.Image {
	var images []*ebiten.Image

	SpriteWidth, SpriteHeight := zones.CalculateFrameDimensions(game.currentZone.SpriteSheetWidth, game.currentZone.SpriteSheetHeight, game.currentZone.SpriteColumns, game.currentZone.SpriteRows)
	SpriteColumns := game.currentZone.SpriteColumns
	spriteRows := game.currentZone.SpriteRows
	totalSprites := SpriteColumns * spriteRows
	game.currentZone.TotalSprites = totalSprites

	for i := 0; i < totalSprites; i++ {

		// Calculate the position of the current sprite
		spriteX := (i % SpriteColumns) * SpriteWidth
		spriteY := (i / SpriteColumns) * SpriteHeight

		// Create a new image for the current sprite
		subImage := playerSprite.SubImage(image.Rect(spriteX, spriteY, spriteX+SpriteWidth, spriteY+SpriteHeight))
		spriteImage := ebiten.NewImageFromImage(subImage)
		images = append(images, spriteImage)
	}

	return images
}

func createImage(filePath string, game *Game) *ebiten.Image {

	img, _, err := ebitenutil.NewImageFromFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var fileName = strings.Split(filePath, "/")[2]

	parts := strings.Split(fileName, "_")
	sheetName := parts[0]
	game.currentZone.SheetName = sheetName
	sheetWidthAndHeight := strings.Split(parts[1], "x")

	var sheetWidth, _ = strconv.Atoi(sheetWidthAndHeight[0])
	var sheetHeight, _ = strconv.Atoi(sheetWidthAndHeight[1])

	game.currentZone.SpriteSheetWidth = sheetWidth
	game.currentZone.SpriteSheetHeight = sheetHeight

	columnsAndRows := strings.Split(strings.Split(parts[2], ".")[0], "x")
	var columns, _ = strconv.Atoi(columnsAndRows[0])
	var rows, _ = strconv.Atoi(columnsAndRows[1])

	game.currentZone.SpriteColumns = columns
	game.currentZone.SpriteRows = rows

	fmt.Printf("\nFilename: %s, sheetName: %s, spriteSheetWidth: %d, spriteSheetHeight: %d, columns: %d, rows %d\n", fileName, sheetName, sheetWidth, sheetHeight, columns, rows)

	return img
}

func createBackgroundImage(filePath string, game *Game) *ebiten.Image {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	finalImage := ebiten.NewImageFromImage(img)

	// Extract the file name from the relative path. Results in "{ZONE}_{SCREEN_WIDTH}x{SCREEN_HEIGHT}_{TILE_SIZE}.png"
	var fileName = strings.Split(filePath, "/")[2]

	// Split fileName into separate parts based on underscore delimiter. Produces array with the three elements: "{ZONE}", "{SCREEN_WIDTH}x{SCREEN_HEIGHT}" and "{TILE_SIZE}.png"
	var parts = strings.Split(fileName, "_")

	// Assign zone from fileName
	var zoneName = parts[0]

	// Assign
	var dimensions = strings.Split(parts[1], "x")
	var screenWidth, _ = strconv.Atoi(dimensions[0])
	var screenHeight, _ = strconv.Atoi(dimensions[1])

	var tileSize, _ = strconv.Atoi(strings.Split(parts[2], ".")[0])

	game.currentZone.Name = zoneName
	game.currentZone.ScreenWidth = screenWidth
	game.currentZone.ScreenHeight = screenHeight
	game.currentZone.TileSize = tileSize

	fmt.Printf("\nFilepath: %s", filePath)
	fmt.Printf("\nZone Name: %s\n", zoneName)
	fmt.Printf("\nscreenWidth: %d, screenHeight: %d, tileSize: %d\n", screenWidth, screenHeight, tileSize)

	return finalImage
}

func (p *Position) ClearChanged() {
	p.justChanged = false
	p.previousChange = false
}

func (p *Position) Increment(x, y int) {
	oldX := p.x
	oldY := p.y
	p.previousX = oldX
	p.previousY = oldY
	p.x += x * game.currentZone.TileSize
	p.y += y * game.currentZone.TileSize
	p.justChanged = true
	p.previousChange = true

	fmt.Printf("Moved from (%d, %d) to (%d, %d)\n", oldX, oldY, p.x, p.y)
}

func main() {
	ebiten.SetWindowSize(game.currentZone.ScreenWidth, game.currentZone.ScreenHeight)
	ebiten.SetWindowTitle("Rattle RPG")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
