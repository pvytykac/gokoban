package model

type Tile int

const (
	Wall     = 0
	Floor    = 1
	DropZone = 3
)

func (tile Tile) ToString() string {
	switch tile {
	case Wall:
		return "Wall"
	case Floor:
		return "Floor"
	case DropZone:
		return "DropZone"
	default:
		panic("unknown tile")
	}
}

func (tile Tile) BlocksMovement() bool {
	return tile == Wall
}
