package main

import (
	"bufio"
	"fmt"
	"github.com/faiface/beep"
	_ "github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/color"
	_ "image/png"
	"log"
	"os"
	"os/signal"
	"strings"
	_ "syscall"
	"time"
)

var grassImage, playerImage, bambooImage, fireImage, roofImage *ebiten.Image
var doorImage, northImage, southImage, westImage, eastImage *ebiten.Image
var waterImages []*ebiten.Image
var alarmSound beep.StreamSeekCloser
var game *Game

const (
	TileTypeGrass  = "G"
	TileTypeBamboo = "B"
	TileTypeFire   = "F"
	TileTypeWater  = "W"
	TileTypeRoof   = "R"
	TileTypeDoor   = "D"

	TileTypeNorth = "NORTH"
	TileTypeWest  = "WEST"
	TileTypeEast  = "EAST"
	TileTypeSouth = "SOUTH"
)

var seaColored bool

func loadTilesFromFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tiles [][]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		tiles = append(tiles, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tiles, nil
}

func (g *Game) LoadTiles(filePath string) error {
	fmt.Printf("\n - - - Loading Tiles - - - \n\n")

	tiles, err := loadTilesFromFile(filePath)
	if err != nil {
		return err
	}

	var fileName = strings.Split(filePath, "/")[2]
	var zoneName = strings.Split(fileName, ".")[0]

	// Update game state with new tiles and name for zone
	g.currentZone = zoneName
	g.currentTiles = tiles

	fmt.Printf("\n - - - Finished Loaded Tiles: for Zone [%s] - - - \n\n", g.currentZone)

	return nil
}

func init() {
	var err error
	grassImage, _, err = ebitenutil.NewImageFromFile("assets/models/grass_2.png")
	playerImage, _, err = ebitenutil.NewImageFromFile("assets/models/rattle_front.png")
	bambooImage, _, err = ebitenutil.NewImageFromFile("assets/models/bamboo.png")

	roofImage, _, err = ebitenutil.NewImageFromFile("assets/models/roof.png")
	doorImage, _, err = ebitenutil.NewImageFromFile("assets/models/door.png")
	fireImage, _, err = ebitenutil.NewImageFromFile("assets/models/fire.png")

	waterImage1, _, err := ebitenutil.NewImageFromFile("assets/models/water.png")
	waterImage2, _, err := ebitenutil.NewImageFromFile("assets/models/water_2.png")
	waterImages = []*ebiten.Image{waterImage1, waterImage2}

	northImage, _, err = ebitenutil.NewImageFromFile("assets/models/north.png")
	southImage, _, err = ebitenutil.NewImageFromFile("assets/models/south.png")
	eastImage, _, err = ebitenutil.NewImageFromFile("assets/models/east.png")
	westImage, _, err = ebitenutil.NewImageFromFile("assets/models/west.png")

	if err != nil {
		log.Fatal(err)
	}
}

type Position struct {
	x           int
	y           int
	justChanged bool
}

func (x *Position) ClearChanged() {
	x.justChanged = false
}

func (a *Position) Increment(x int, y int) {
	a.x += x
	a.y += y
	a.justChanged = true
}

const (
	SCREEN_SIZE = 500
	TILE_SIZE   = 25
	NUM_TILES   = SCREEN_SIZE / TILE_SIZE
)

type Game struct {
	currentTiles    [][]string
	currentZone     string
	currentPos      Position
	waterImageIndex int
}

func (g *Game) Update() error {
	g.HandleKeyPress()
	if g.currentPos.justChanged {
		if g.currentZone == "island" {

			if g.currentPos.x <= 200 && g.currentPos.y <= 100 && !seaColored {

				//// Color the sea by changing all the "W" currentTiles to "F" (fire)
				//for i := range g.currentTiles {
				//	for j := range g.currentTiles[i] {
				//		if g.currentTiles[i][j] == TileTypeWater {
				//			g.currentTiles[i][j] = TileTypeFire
				//		}
				//	}
				//}
				//seaColored = true
				//
				//err := alarmSound.Seek(0)
				//if err != nil {
				//	return nil
				//}
				//speaker.Play(alarmSound)
			}
		}
	}
	return nil
}

func loadAlarmSound() {
	f, err := os.Open("assets/sounds/alarm.wav")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return
	}
	alarmSound = streamer
}

func (g *Game) HandleKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.currentPos.Increment(0, -1)
		g.currentTiles[18][18] = TileTypeNorth // Set the tile type to the north image
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.currentPos.Increment(0, 1)
		g.currentTiles[18][18] = TileTypeSouth // Set the tile type to the south image
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.currentPos.Increment(-1, 0)
		g.currentTiles[18][18] = TileTypeWest // Set the tile type to the west image
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.currentPos.Increment(1, 0)
		g.currentTiles[18][18] = TileTypeEast // Set the tile type to the east image
	}

	g.currentPos.justChanged = true
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < NUM_TILES; i++ {
		for j := 0; j < NUM_TILES; j++ {
			if i >= len(g.currentTiles) || j >= len(g.currentTiles[i]) {
				continue
			}

			tileType := g.currentTiles[i][j]

			switch tileType {
			case TileTypeGrass:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(grassImage, op)
			case TileTypeBamboo:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(bambooImage, op)
			case TileTypeFire:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(fireImage, op)
			case TileTypeWater:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(waterImages[g.waterImageIndex], op)
			case TileTypeRoof:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(roofImage, op)
			case TileTypeDoor:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(doorImage, op)

			// Add new cases for the north, south, west, and east images
			case TileTypeNorth:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(northImage, op)
			case TileTypeSouth:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(southImage, op)
			case TileTypeWest:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(westImage, op)
			case TileTypeEast:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
				screen.DrawImage(eastImage, op)
			}

		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(TILE_SIZE)/float64(playerImage.Bounds().Dx()), float64(TILE_SIZE)/float64(playerImage.Bounds().Dy()))
	op.GeoM.Translate(float64(g.currentPos.x), float64(g.currentPos.y))

	screen.DrawImage(playerImage, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Zone: [%s], Coordinates: (%d, %d)", g.currentZone, g.currentPos.x, g.currentPos.y))

	//log.Printf("Coordinates (%d, %d) \n", g.currentPos.x, g.currentPos.y)

	// Additional check to draw colored sea if player is on the specified coordinate
	if g.currentZone == "island" {
		if g.currentPos.x <= 200 && g.currentPos.y <= 100 && seaColored {
			for i := 0; i < NUM_TILES; i++ {
				for j := 0; j < NUM_TILES; j++ {
					if g.currentTiles[i][j] == TileTypeFire {
						op := &ebiten.DrawImageOptions{}
						op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
						screen.DrawImage(fireImage, op)
					}
				}
			}

		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_SIZE, SCREEN_SIZE
}

func NewPosition() Position {
	return Position{
		justChanged: true,
	}
}

func NewGame() *Game {
	tiles := make([][]string, NUM_TILES)
	for i := range tiles {
		tiles[i] = make([]string, NUM_TILES)
	}

	return &Game{
		currentTiles: tiles,
		currentPos:   NewPosition(),
	}
}

func main() {
	ebiten.SetWindowSize(SCREEN_SIZE, SCREEN_SIZE)
	ebiten.SetWindowTitle("Pondi Island RPG")

	playerTileX := NUM_TILES / 2
	playerTileY := NUM_TILES / 2
	playerPosX := playerTileX * TILE_SIZE
	playerPosY := playerTileY * TILE_SIZE

	game = NewGame()
	game.currentPos.x = playerPosX
	game.currentPos.y = playerPosY

	loadAlarmSound()

	// Load currentTiles from a file
	err := game.LoadTiles("assets/areas/island.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Set up the game update ticker
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Set up the water image switching ticker
	waterTicker := time.NewTicker(time.Second)
	defer waterTicker.Stop()

	// Create a channel to receive termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Goroutines to change color of water
	go func() {

		if game.currentZone == "island" {
			for {
				select {
				case <-ticker.C:
					if err := game.Update(); err != nil {
						log.Fatal(err)
					}
				case <-waterTicker.C:
					game.waterImageIndex = (game.waterImageIndex + 1) % len(waterImages)
				case <-quit:
					// Handle termination signal
					fmt.Println("\nReceived termination signal. Exiting...")
					os.Exit(0)
				}
			}
		}
	}() // X = 228, Y > 216 && Y < 245

	// Handle Ctrl+C signal to gracefully exit the program
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	<-quit
	fmt.Println("Exiting...")
}
