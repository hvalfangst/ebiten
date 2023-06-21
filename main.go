package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SCREEN_WIDTH  = 960
	SCREEN_HEIGHT = 640
	TILE_SIZE     = 32
)

var (
	backgroundImage, playerImage *ebiten.Image
	currentGame                  *Game
)

func init() {
	currentGame = &Game{}
	backgroundImage = createBackgroundImage("assets/areas/reindal_960x640.png")
	//currentGame.currentDimensions = dimensions
	//currentGame.currentZone = zone
	playerImage = createImage("assets/models/fire.png")
}

func createImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
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

	var fileName = strings.Split(filePath, "/")[2]
	var zoneName = strings.Split(fileName, "_")[0]
	var dimensionsWithExtension = strings.Split(fileName, "_")[1]
	var dimensions = strings.Split(dimensionsWithExtension, ".")[0]

	currentGame.currentZone = zoneName
	currentGame.currentDimensions = dimensions

	fmt.Printf("\nZone Name: %s\n", zoneName)
	fmt.Printf("\nDimensions: %s\n", dimensions)

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
	p.x += x * TILE_SIZE
	p.y += y * TILE_SIZE
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
		} else if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.currentPos.Increment(0, 1)
			g.movementHandled = true
		} else if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.currentPos.Increment(-1, 0)
			g.movementHandled = true
		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.currentPos.Increment(1, 0)
			g.movementHandled = true
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

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(TILE_SIZE)/float64(playerImage.Bounds().Dx()), float64(TILE_SIZE)/float64(playerImage.Bounds().Dy()))
	op.GeoM.Translate(float64(g.currentPos.x), float64(g.currentPos.y))

	screen.DrawImage(playerImage, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Zone: [%s], Coordinates: (%d, %d)", g.currentZone, g.currentPos.x, g.currentPos.y))

	// Draw other game elements on top of the background
	// ...

	// Draw debug information, if needed
	// ...
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Desert Game")

	if err := ebiten.RunGame(currentGame); err != nil {
		log.Fatal(err)
	}
}
