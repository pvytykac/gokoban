package game

import "gokoban/game/model"

type GameState struct {
	Level           int
	LevelCount      int
	ElapsedTime     int64
	LevelTime       int64
	PlayerDirection model.Direction
	World           [][]BoardTile
}

func NewGameRepresentation(game *model.Game) *GameState {
	level := game.CurrentLevel
	player := level.Player
	height := level.Height
	width := level.Width
	boxes := level.Boxes
	world := make([][]BoardTile, height)

	for y := 0; y < height; y++ {
		world[y] = make([]BoardTile, width)
		for x := 0; x < width; x++ {
			var tile BoardTile
			if player.Position.Matches(x, y) {
				tile = Player
			} else {
				tile = tileToBoardTile(level.GetTileAtCoordinates(x, y))
			}
			world[y][x] = tile
		}
	}

	for _, box := range *boxes {
		world[box.Y][box.X] = Box
	}

	return &GameState{
		Level:           game.LevelIndex,
		LevelCount:      game.GetLevelCount(),
		ElapsedTime:     game.GetElapsedTime(),
		LevelTime:       game.CurrentLevel.GetElapsedTime(),
		PlayerDirection: game.CurrentLevel.Player.Direction,
		World:           world,
	}
}

func tileToBoardTile(tile model.Tile) BoardTile {
	switch tile {
	case model.Wall:
		return Wall
	case model.Floor:
		return Floor
	case model.DropZone:
		return DropZone
	default:
		panic("unknown tile type")
	}
}

type BoardTile int

const (
	Wall     BoardTile = 0
	Floor    BoardTile = 1
	Player   BoardTile = 2
	Box      BoardTile = 3
	DropZone BoardTile = 4
)
