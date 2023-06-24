package zones

type Zone struct {
	Name              string
	TileSize          int
	ScreenWidth       int
	ScreenHeight      int
	SheetName         string
	SpriteSheetWidth  int
	SpriteSheetHeight int
	SpriteColumns     int
	SpriteRows        int
	TotalSprites      int
}

func CalculateFrameDimensions(SheetWidth int, SheetHeight int, NumColumns int, NumRows int) (int, int) {
	frameWidth := SheetWidth / NumColumns
	frameHeight := SheetHeight / NumRows
	return frameWidth, frameHeight
}
