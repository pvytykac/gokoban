package command

import "gokoban/game/model"

type Command interface {
	Apply(level *model.Level) bool
}
