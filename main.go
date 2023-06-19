package main

import (
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
	"math/rand"
	"os"
	"os/signal"
	_ "syscall"
	"time"
)

var grassImage *ebiten.Image
var playerImage *ebiten.Image
var bambooImage *ebiten.Image
var fireImage *ebiten.Image
var game Game

const (
	FireTileDuration = 3 * time.Second
)

var fireSound beep.StreamSeekCloser

func loadFireSound() {
	f, err := os.Open("fire.wav")
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

func init() {
	var err error
	grassImage, _, err = ebitenutil.NewImageFromFile("grass.png")
	playerImage, _, err = ebitenutil.NewImageFromFile("player.png")
	bambooImage, _, err = ebitenutil.NewImageFromFile("bamboo.png")
	fireImage, _, err = ebitenutil.NewImageFromFile("fire.png")
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

const (
	TileTypeGrass = iota
	TileTypeBamboo
	TileTypeFire
)

type Game struct {
	tiles      [][]int
	currentPos Position
}

func (g *Game) Update() error {
	g.HandleKeyPress()

	if g.currentPos.justChanged {
		playerTileX := (NUM_TILES / 2) + g.currentPos.x
		playerTileY := (NUM_TILES / 2) + g.currentPos.y
		if playerTileX >= 0 && playerTileX < NUM_TILES && playerTileY >= 0 && playerTileY < NUM_TILES {
			if g.tiles[playerTileX][playerTileY] != TileTypeFire {
				g.tiles[playerTileX][playerTileY] = TileTypeGrass
			}
		}
		g.currentPos.justChanged = false
	}

	return nil
}

func (g *Game) HandleKeyPress() {
	oldX, oldY := g.currentPos.x, g.currentPos.y

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.currentPos.Increment(0, -1)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.currentPos.Increment(0, 1)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.currentPos.Increment(-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.currentPos.Increment(1, 0)
	}

	newX, newY := g.currentPos.x, g.currentPos.y
	log.Printf("Moved from (%d, %d) to (%d, %d)\n", oldX, oldY, newX, newY)
	g.currentPos.justChanged = true

	// Check if the player is moving onto a fire tile
	playerTileX := (NUM_TILES / 2) + g.currentPos.x
	playerTileY := (NUM_TILES / 2) + g.currentPos.y
	if playerTileX >= 0 && playerTileX < NUM_TILES && playerTileY >= 0 && playerTileY < NUM_TILES {
		if g.tiles[playerTileX][playerTileY] == TileTypeFire {
			g.tiles[playerTileX][playerTileY] = TileTypeGrass
		}
	}
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
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(TILE_SIZE)/float64(playerImage.Bounds().Dx()), float64(TILE_SIZE)/float64(playerImage.Bounds().Dy()))
	op.GeoM.Translate(float64(g.currentPos.x), float64(g.currentPos.y))
	screen.DrawImage(playerImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_SIZE, SCREEN_SIZE
}

func NewPosition() Position {
	return Position{
		justChanged: true,
	}
}

func NewGame() Game {
	tiles := make([][]int, NUM_TILES)
	for i := range tiles {
		tiles[i] = make([]int, NUM_TILES)
	}

	// Set some initial fire tiles
	tiles[2][2] = TileTypeFire
	tiles[4][6] = TileTypeFire
	tiles[7][3] = TileTypeFire

	return Game{
		tiles:      tiles,
		currentPos: NewPosition(),
	}
}

func main() {
	ebiten.SetWindowSize(SCREEN_SIZE, SCREEN_SIZE)
	ebiten.SetWindowTitle("Bamboo Forest RPG")

	playerTileX := NUM_TILES / 2
	playerTileY := NUM_TILES / 2
	playerPosX := playerTileX * TILE_SIZE
	playerPosY := playerTileY * TILE_SIZE

	game = NewGame()
	game.currentPos.x = playerPosX
	game.currentPos.y = playerPosY

	loadFireSound()

	// Create a goroutine to periodically spawn fire tiles
	go func() {
		for {
			// Randomly select a tile to spawn fire
			tileX := rand.Intn(NUM_TILES)
			tileY := rand.Intn(NUM_TILES)

			// Check if the tile is already occupied by fire
			if game.tiles[tileX][tileY] != TileTypeFire {
				game.tiles[tileX][tileY] = TileTypeFire
				// Play the fire sound
				fireSound.Seek(0)
				speaker.Play(fireSound)
			}

			time.Sleep(FireTileDuration)
		}
	}()

	// Handle Ctrl+C signal to gracefully exit the program
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}

	<-quit
	fmt.Println("Exiting...")
}
