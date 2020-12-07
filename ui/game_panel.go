package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gokoban/game"
	"gokoban/game/model"
	"gokoban/ui/component"
)

type GamePanel struct {
	*component.MyBox
	state *game.GameState
}

func newGamePanel() *GamePanel {
	return &GamePanel{
		MyBox: component.NewMyBox("game"),
		state: &game.GameState{ElapsedTime: 0, LevelCount: 0, Level: 0, LevelTime: 0},
	}
}

func (view *GamePanel) Draw(screen tcell.Screen) {
	view.Box.DrawForSubclass(screen, view)
	x, y, width, _ := view.GetInnerRect()

	direction := view.state.PlayerDirection
	world := view.state.World

	for row := 0; row < len(world); row++ {
		line := ""
		for col := 0; col < len(world[row]); col++ {
			line += toString(world[row][col], direction)
		}
		tview.Print(screen, line, x, y + row, width, tview.AlignCenter, tcell.ColorYellow)
	}
}

func (view *GamePanel) SetState(state *game.GameState) {
	view.state = state
}

func toString(tile game.BoardTile, direction model.Direction) string {
	switch tile {
	case game.Wall:
		return "[gray]▓"
	case game.Floor:
		return " "
	case game.Box:
		return "[brown]▒"
	case game.Player:
		prefix := "[blue]"
		switch direction {
		case model.North:
			return prefix + "^"
		case model.East:
			return prefix + ">"
		case model.South:
			return prefix + "V"
		case model.West:
			return prefix + "<"
		default:
			return prefix + "0"
		}
	case game.DropZone:
		return "[red]░"
	default:
		return "/"
	}
}
