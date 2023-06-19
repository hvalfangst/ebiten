package main

import (
	"fmt"
	_ "image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var grassImage *ebiten.Image
var playerImage *ebiten.Image
var game Game

func init() {
	var err error
	grassImage, _, err = ebitenutil.NewImageFromFile("grass.png")
	playerImage, _, err = ebitenutil.NewImageFromFile("player.png")
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
	tiles      [][]float64
	currentPos Position
}

func (g *Game) Update() error {
	g.HandleKeyPress()

	if g.currentPos.justChanged {
		playerTileX := (NUM_TILES / 2) + g.currentPos.x
		playerTileY := (NUM_TILES / 2) + g.currentPos.y
		if playerTileX >= 0 && playerTileX < NUM_TILES && playerTileY >= 0 && playerTileY < NUM_TILES {
			g.tiles[playerTileX][playerTileY] = 1.0
		}
		g.currentPos.justChanged = false
	}

	return nil
}

func (g *Game) HandleKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		oldX, oldY := g.currentPos.x, g.currentPos.y
		g.currentPos.Increment(0, -1)
		newX, newY := g.currentPos.x, g.currentPos.y
		log.Printf("Moved from (%d, %d) to (%d, %d)\n", oldX, oldY, newX, newY)
		g.currentPos.justChanged = true
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		oldX, oldY := g.currentPos.x, g.currentPos.y
		g.currentPos.Increment(0, 1)
		newX, newY := g.currentPos.x, g.currentPos.y
		log.Printf("Moved from (%d, %d) to (%d, %d)\n", oldX, oldY, newX, newY)
		g.currentPos.justChanged = true
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		oldX, oldY := g.currentPos.x, g.currentPos.y
		g.currentPos.Increment(-1, 0)
		newX, newY := g.currentPos.x, g.currentPos.y
		log.Printf("Moved from (%d, %d) to (%d, %d)\n", oldX, oldY, newX, newY)
		g.currentPos.justChanged = true
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		oldX, oldY := g.currentPos.x, g.currentPos.y
		g.currentPos.Increment(1, 0)
		newX, newY := g.currentPos.x, g.currentPos.y
		log.Printf("Moved from (%d, %d) to (%d, %d)\n", oldX, oldY, newX, newY)
		g.currentPos.justChanged = true
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, tileRow := range g.tiles {
		for j := range tileRow {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i*TILE_SIZE), float64(j*TILE_SIZE))
			screen.DrawImage(grassImage, op)
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
	tiles := make([][]float64, NUM_TILES)
	for i := range tiles {
		tiles[i] = make([]float64, NUM_TILES)
	}

	return Game{
		tiles:      tiles,
		currentPos: NewPosition(),
	}
}

func main() {
	ebiten.SetWindowSize(SCREEN_SIZE, SCREEN_SIZE)
	ebiten.SetWindowTitle("Rattle RPG")

	playerTileX := NUM_TILES / 2
	playerTileY := NUM_TILES / 2
	playerPosX := playerTileX * TILE_SIZE
	playerPosY := playerTileY * TILE_SIZE

	game = NewGame()
	game.currentPos.x = playerPosX
	game.currentPos.y = playerPosY

	fmt.Println("Initial player coordinates:", game.currentPos.x, game.currentPos.y)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
