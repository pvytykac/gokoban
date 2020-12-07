package model

type Player struct {
	Position *Position
	Direction Direction
}

func NewPlayer(position *Position, direction Direction) *Player {
	return &Player{Position: position, Direction: direction}
}
