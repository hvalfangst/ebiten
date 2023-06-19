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

var grassImage *ebiten.Image
var playerImage *ebiten.Image
var bambooImage *ebiten.Image
var fireImage *ebiten.Image

var roofImage *ebiten.Image
var doorImage *ebiten.Image

var waterImages []*ebiten.Image

var game *Game

var fireSound beep.StreamSeekCloser

const (
	TileTypeGrass    = "G"
	TileTypeBamboo   = "B"
	TileTypeFire     = "F"
	TileTypeWater    = "W"
	TileTypeRoof     = "R"
	TileTypeDoor     = "D"
	FireTileDuration = 5 * time.Second
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

func (g *Game) LoadTiles(filename string) error {
	tiles, err := loadTilesFromFile(filename)
	if err != nil {
		return err
	}

	g.tiles = tiles

	return nil
}

func init() {
	var err error
	grassImage, _, err = ebitenutil.NewImageFromFile("grass_2.png")
	playerImage, _, err = ebitenutil.NewImageFromFile("player.png")
	bambooImage, _, err = ebitenutil.NewImageFromFile("bamboo.png")

	roofImage, _, err = ebitenutil.NewImageFromFile("roof.png")
	doorImage, _, err = ebitenutil.NewImageFromFile("door.png")
	fireImage, _, err = ebitenutil.NewImageFromFile("fire.png")

	waterImage1, _, err := ebitenutil.NewImageFromFile("water.png")
	waterImage2, _, err := ebitenutil.NewImageFromFile("water_2.png")
	waterImages = []*ebiten.Image{waterImage1, waterImage2}

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
	SCREEN_SIZE       = 500
	TILE_SIZE         = 25
	NUM_TILES         = SCREEN_SIZE / TILE_SIZE
	PLAYER_SIZE       = TILE_SIZE
	TOTAL_GRASS_TILES = 400
)

type Game struct {
	tiles           [][]string
	currentPos      Position
	waterImageIndex int
}

func (g *Game) Update() error {
	g.HandleKeyPress()

	if g.currentPos.justChanged {
		if g.currentPos.x <= 200 && g.currentPos.y <= 100 && !seaColored {
			// Color the sea by changing all the "W" tiles to "F" (fire)
			for i := range g.tiles {
				for j := range g.tiles[i] {
					if g.tiles[i][j] == TileTypeWater {
						g.tiles[i][j] = TileTypeFire
					}
				}
			}
			seaColored = true

			err := fireSound.Seek(0)
			if err != nil {
				return nil
			}
			speaker.Play(fireSound)
		}

	}

	return nil
}

// ... (previous code)

func loadFireSound() {
	f, err := os.Open("alarm.wav")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	fireSound = streamer
}

func (g *Game) HandleKeyPress() {
	//oldX, oldY := g.currentPos.x, g.currentPos.y

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.currentPos.Increment(0, -1)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.currentPos.Increment(0, 1)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.currentPos.Increment(-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.currentPos.Increment(1, 0)
	}

	//newX, newY := g.currentPos.x, g.currentPos.y
	//log.Printf("Moved from (%d, %d) to (%d, %d)\n", oldX, oldY, newX, newY)
	g.currentPos.justChanged = true
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, tileRow := range g.tiles {
		for j, tileType := range tileRow {
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
			}

		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(TILE_SIZE)/float64(playerImage.Bounds().Dx()), float64(TILE_SIZE)/float64(playerImage.Bounds().Dy()))
	op.GeoM.Translate(float64(g.currentPos.x), float64(g.currentPos.y))
	screen.DrawImage(playerImage, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Coordinates: (%d, %d)", g.currentPos.x, g.currentPos.y))

	log.Printf("Coordinates (%d, %d) \n", g.currentPos.x, g.currentPos.y)

	// Additional check to draw colored sea if player is on the specified coordinate
	if g.currentPos.x <= 200 && g.currentPos.y <= 100 && seaColored {
		for i := 0; i < NUM_TILES; i++ {
			for j := 0; j < NUM_TILES; j++ {
				if g.tiles[i][j] == TileTypeFire {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
					screen.DrawImage(fireImage, op)
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
		tiles:      tiles,
		currentPos: NewPosition(),
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

	loadFireSound()

	// Load tiles from a file
	err := game.LoadTiles("map.txt")
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
	}()

	// Handle Ctrl+C signal to gracefully exit the program
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	<-quit
	fmt.Println("Exiting...")
}
