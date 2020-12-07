package model

type Position struct {
	X int
	Y int
}

func NewPosition(x int, y int) *Position {
	return &Position{X: x, Y: y}
}

func (position *Position) Matches(x int, y int) bool {
	return position.X == x && position.Y == y
}