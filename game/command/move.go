package command

import "gokoban/game/model"

type MoveCommand struct {
	direction model.Direction
}

func NewMoveCommand(direction model.Direction) *MoveCommand {
	return &MoveCommand{direction: direction}
}

func (cmd *MoveCommand) Apply(level *model.Level) bool {
	return level.MoveTo(cmd.direction)
}
