package model

type Level struct {
	Width  int
	Height int
	Player *Player
	World  *[][]Tile
	Boxes  *[]*Position
	Next   *Level
	time   int64
}

func NewLevel(width int, height int, player *Player, world *[][]Tile, boxes *[]*Position) *Level {
	return &Level{Width: width, Height: height, Player: player, World: world, Boxes: boxes}
}

func (level *Level) IsSolved() bool {
	for _, box := range *level.Boxes {
		if tile := level.GetTileAtPosition(box); tile != DropZone {
			return false
		}
	}

	return true
}

func (level *Level) MoveTo(direction Direction) bool {
	target := direction.Project(level.Player.Position)
	if level.CanMoveTo(target, direction) {
		level.Player.Position = target
		if box := level.GetBoxAtPosition(target); box != nil {
			boxTarget := direction.Project(box)
			box.X = boxTarget.X
			box.Y = boxTarget.Y
		}
		level.Player.Direction = direction
		return true
	} else if level.Player.Direction != direction {
		level.Player.Direction = direction
		return false
	}

	return false
}

func (level *Level) CanMoveTo(target *Position, direction Direction) bool {
	if level.IsBoxAtPosition(target) {
		boxTarget := direction.Project(target)
		return !level.IsBoxAtPosition(boxTarget) && !level.GetTileAtPosition(boxTarget).BlocksMovement()
	} else {
		return level.IsPositionInsideWorld(target) && !level.GetTileAtPosition(target).BlocksMovement()
	}
}

func (level *Level) IsBoxAtPosition(position *Position) bool {
	return level.IsBoxAtCoordinates(position.X, position.Y)
}

func (level *Level) IsBoxAtCoordinates(x int, y int) bool {
	return level.GetBoxAtCoordinates(x, y) != nil
}

func (level *Level) GetBoxAtPosition(position *Position) *Position {
	return level.GetBoxAtCoordinates(position.X, position.Y)
}

func (level *Level) GetBoxAtCoordinates(x int, y int) *Position {
	for _, box := range *level.Boxes {
		if box.X == x && box.Y == y {
			return box
		}
	}

	return nil
}

func (level *Level) IsPositionInsideWorld(position *Position) bool {
	return position.Y < level.Height && position.Y >= 0 && position.X < level.Width && position.X >= 0
}

func (level *Level) GetTileAtPosition(position *Position) Tile {
	return (*level.World)[position.Y][position.X]
}

func (level *Level) GetTileAtCoordinates(x int, y int) Tile {
	return (*level.World)[y][x]
}

func (level *Level) IncrementTime(inc int64) {
	level.time += inc
}

func (level *Level) GetElapsedTime() int64 {
	return level.time
}

func (level *Level) setNextLevel(next *Level) {
	if level.Next == nil {
		level.Next = next
	} else {
		level.Next.setNextLevel(next)
	}
}
