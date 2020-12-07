package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gokoban/game"
	"gokoban/ui/component"
)

type InfoPanel struct {
	*component.MyBox
	state *game.GameState
}

func newInfoPanel() *InfoPanel {
	return &InfoPanel{
		MyBox: component.NewMyBox("info"),
		state: &game.GameState{ElapsedTime: 0, LevelCount: 0, Level: 0, LevelTime: 0},
	}
}

func (view *InfoPanel) Draw(screen tcell.Screen) {
	view.Box.DrawForSubclass(screen, view)
	x, y, width, _ := view.GetInnerRect()

	lines := []string{
		fmt.Sprintf("Level: %d/%d", view.state.Level, view.state.LevelCount),
		fmt.Sprintf("Level Time: %s", toHumanReadableTime(view.state.LevelTime)),
		fmt.Sprintf("Total Time: %s", toHumanReadableTime(view.state.ElapsedTime)),
	}

	for ix, line := range lines {
		tview.Print(screen, line, x, y+ix, width, tview.AlignCenter, tcell.ColorYellow)
	}
}

func (view *InfoPanel) SetState(state *game.GameState) {
	view.state = state
}

func toHumanReadableTime(seconds int64) string {
	if seconds > 3600 {
		return fmt.Sprintf("%d:%02d:%02d", seconds / 3600, seconds % 3600 / 60, seconds % 3600 % 60)
	} else {
		return fmt.Sprintf("%02d:%02d", seconds / 60, seconds % 60)
	}
}