package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gokoban/game"
	"gokoban/ui/component"
)

type GameView struct {
	*tview.Flex
	infoPanel *InfoPanel
	gamePanel *GamePanel
}

func newGameView(bindings *[]*component.KeyEvent) *GameView {
	infoPanel := newInfoPanel()
	gamePanel := newGamePanel()
	gamePanel.Focus(func(_ tview.Primitive){})

	gamePanel.Box.SetInputCapture(onKeyPressHandler(bindings))

	return &GameView{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(infoPanel, 5, 1, false).
			AddItem(gamePanel, 0, 7, true),
		infoPanel: infoPanel,
		gamePanel: gamePanel,
	}
}

func (view *GameView) SetGameState(state *game.GameState) {
	view.infoPanel.state = state
	view.gamePanel.state = state
}

func onKeyPressHandler(bindings *[]*component.KeyEvent) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		for _, binding := range *bindings {
			if binding.AppliesTo(event) {
				binding.Execute()
				return nil
			}
		}

		return event
	}
}
