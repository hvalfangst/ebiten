package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TileSize = 32

	//SpriteWidth   = 48 // Width of each sprite in pixels
	//SpriteHeight  = 72 // Height of each sprite in pixels
	//SpriteColumns = 12 // Number of columns in the sprite-sheet
)

var (
	backgroundImage, playerImage *ebiten.Image
	currentGame                  *Game
	spritesheetImage             *ebiten.Image
	spriteImages                 []*ebiten.Image
	currentSpriteIdx             int
)

func init() {
	currentGame = &Game{}
	backgroundImage = createBackgroundImage("assets/areas/reindal_960x640_32.png")
	//currentGame.currentDimensions = dimensions
	//currentGame.currentZone = zone

	// Load the sprite-sheet image
	spritesheetImage = createImage("assets/models/avatars_576x576_12x8.png")
	//9
	//frameWidt, frameHeight := CalculateFrameDimensions(576, 576, 12, 8)
	//fmt.Printf("Framewidt: %d, FrameHeight: %d", frameWidt, frameHeight)
	//

	// Create sprite images from the sprite-sheet
	spriteImages = createSpriteImages()
	//playerImage = createImage("assets/models/fire.png")
}

// Calculates pixels based on dimensions of a given file, the number of columns and number of rows
func CalculateFrameDimensions(SheetWidth int, SheetHeight int, NumColumns int, NumRows int) (int, int) {
	frameWidth := SheetWidth / NumColumns
	frameHeight := SheetHeight / NumRows
	return frameWidth, frameHeight
}

func createSpriteImages() []*ebiten.Image {
	var images []*ebiten.Image

	SpriteHeight := currentGame.spriteHeight
	SpriteWidth := currentGame.spriteWidth
	SpriteColumns := currentGame.spriteColumns

	// Calculate the number of sprites in the sprite-sheet
	spriteRows := spritesheetImage.Bounds().Dy() / SpriteHeight
	totalSprites := SpriteColumns * spriteRows
	currentGame.totalSprites = totalSprites

	for i := 0; i < totalSprites; i++ {
		// Calculate the position of the current sprite
		spriteX := (i % SpriteColumns) * SpriteWidth
		spriteY := (i / SpriteColumns) * SpriteHeight

		// Create a new image for the current sprite
		subImage := spritesheetImage.SubImage(image.Rect(spriteX, spriteY, spriteX+SpriteWidth, spriteY+SpriteHeight))
		spriteImage := ebiten.NewImageFromImage(subImage)

		images = append(images, spriteImage)
	}

	return images
}

func createImage(filePath string) *ebiten.Image {

	img, _, err := ebitenutil.NewImageFromFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var fileName = strings.Split(filePath, "/")[2]

	parts := strings.Split(fileName, "_")
	sheetName := parts[0]
	currentGame.sheetName = sheetName
	sheetWidthAndHeight := strings.Split(parts[1], "x")

	var sheetWidth, _ = strconv.Atoi(sheetWidthAndHeight[0])
	var sheetHeight, _ = strconv.Atoi(sheetWidthAndHeight[1])

	currentGame.spriteWidth = sheetWidth
	currentGame.spriteHeight = sheetHeight

	columnsAndRows := strings.Split(strings.Split(parts[2], ".")[0], "x")
	var columns, _ = strconv.Atoi(columnsAndRows[0])
	var rows, _ = strconv.Atoi(columnsAndRows[1])

	currentGame.spriteColumns = columns
	currentGame.spriteRows = rows

	fmt.Printf("\nFilename: %s, sheetName: %s, spriteWidth: %d, spriteHeight: %d, columns: %d, rows %d\n", fileName, sheetName, sheetWidth, sheetHeight, columns, rows)

	os.Exit(1)
	return img
}

func createBackgroundImage(filePath string) *ebiten.Image {
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

	currentGame.currentZone = zoneName
	currentGame.screenWidth = screenWidth
	currentGame.screenHeight = screenHeight
	currentGame.tileSize = tileSize

	fmt.Printf("\nFilepath: %s", filePath)
	fmt.Printf("\nZone Name: %s\n", zoneName)
	fmt.Printf("\nscreenWidth: %d, screenHeight: %d, tileSize: %d\n", screenWidth, screenHeight, tileSize)

	return finalImage
}

type Position struct {
	x              int
	y              int
	justChanged    bool
	previousX      int
	previousY      int
	previousChange bool
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
	p.x += x * TileSize
	p.y += y * TileSize
	p.justChanged = true
	p.previousChange = true

	fmt.Printf("Moved from (%d, %d) to (%d, %d)\n", oldX, oldY, p.x, p.y)
}

type Game struct {
	//currentTiles      [][]string
	currentZone       string
	currentPos        Position
	waterImageIndex   int
	currentDimensions string
	movementHandled   bool
	totalSprites      int
	sheetName         string
	spriteWidth       int
	spriteHeight      int
	spriteColumns     int
	spriteRows        int
	screenWidth       int
	screenHeight      int
	tileSize          int
}

func (g *Game) Update() error {
	g.HandleKeyPress()
	if g.currentPos.justChanged {
		if g.currentZone == "island" {

		}
	}
	return nil
}
func (g *Game) HandleKeyPress() {
	if !g.movementHandled {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.currentPos.Increment(0, -1)
			g.movementHandled = true
			currentSpriteIdx = 47
		} else if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.currentPos.Increment(0, 1)
			g.movementHandled = true
			currentSpriteIdx = 47
		} else if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.currentPos.Increment(-1, 0)
			g.movementHandled = true
			currentSpriteIdx = 47
		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.currentPos.Increment(1, 0)
			g.movementHandled = true
			currentSpriteIdx = 47
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

	//op = &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(float64(TileSize)/float64(playerImage.Bounds().Dx()), float64(TileSize)/float64(playerImage.Bounds().Dy()))
	//op.GeoM.Translate(float64(g.currentPos.x), float64(g.currentPos.y))
	//
	//screen.DrawImage(playerImage, op)

	//op.GeoM.Scale(float64(TileSize)/float64(spriteImages[currentSpriteIdx].Bounds().Dx()), float64(TileSize)/float64(spriteImages[currentSpriteIdx].Bounds().Dy()))
	op.GeoM.Translate(float64(g.currentPos.x), float64(g.currentPos.y))

	// Draw the current sprite image
	screen.DrawImage(spriteImages[currentSpriteIdx], op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Zone: [%s], Coordinates: (%d, %d), Sprites: %d", g.currentZone, g.currentPos.x, g.currentPos.y, g.totalSprites))

	// Draw other game elements on top of the background
	// ...

	// Draw debug information, if needed
	// ...
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return currentGame.screenWidth, currentGame.screenHeight
}

func main() {
	ebiten.SetWindowSize(currentGame.spriteWidth, currentGame.spriteHeight)
	ebiten.SetWindowTitle("Desert Game")

	if err := ebiten.RunGame(currentGame); err != nil {
		log.Fatal(err)
	}
}
