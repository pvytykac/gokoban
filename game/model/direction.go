package model

type Direction int

const (
	North = 0
	East  = 90
	South = 180
	West  = 270
)

func (dir Direction) ToString() string {
	switch dir {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		panic("unknown direction")
	}
}

func (dir Direction) Project(position *Position) *Position {
	var xOffset, yOffset int
	switch dir {
	case North:
		yOffset = -1
	case East:
		xOffset = 1
	case South:
		yOffset = 1
	case West:
		xOffset = -1
	default:
		panic("unknown direction")
	}

	return NewPosition(position.X + xOffset, position.Y + yOffset)
}
