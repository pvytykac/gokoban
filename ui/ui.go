package ui

import (
	"github.com/rivo/tview"
	"gokoban/game"
	"gokoban/game/command"
	"gokoban/game/model"
	"gokoban/ui/component"
	"sync"
)

type Ui struct {
	application    *tview.Application
	controller     *game.Controller
	pages          *tview.Pages
	backFunction   func()
	channelLock    sync.Mutex
	stateChannel   <-chan *game.GameState
	commandChannel chan<- command.Command
	gameView       *GameView
}

func NewUi(controller *game.Controller) *Ui {
	return &Ui{
		application:    tview.NewApplication(),
		controller:     controller,
		pages:          nil,
		backFunction:   nil,
		channelLock:    sync.Mutex{},
		stateChannel:   nil,
		commandChannel: nil,
		gameView:       nil,
	}
}

func (ui *Ui) Start() {
	ui.pages = ui.buildUi()
	ui.application.SetRoot(ui.pages, true).
		EnableMouse(true)

	if err := ui.application.Run(); err != nil {
		panic(err)
	}
}

func (ui *Ui) exit() {
	ui.application.Stop()
}

func (ui *Ui) startGame() {
	ui.setGameChannels(ui.controller.Start())
	go ui.listenForStateChanges()
	ui.switchToPage(PlayPage, nil)
}

func (ui *Ui) resumeGame() {
	ui.setGameChannels(ui.controller.Resume())
	go ui.listenForStateChanges()
	ui.switchToPage(PlayPage, nil)
}

func (ui *Ui) restartLevel() {
	ui.controller.RestartLevel()
	ui.resumeGame()
}

func (ui *Ui) switchToMainMenu() {
	ui.switchToPage(MainMenuPage, nil)
}

func (ui *Ui) switchToGameMenu() {
	ui.switchToPage(GameMenuPage, nil)
	ui.setGameChannels(nil, nil)
	ui.controller.Pause()
}

func (ui *Ui) switchToQuitModal(backPage Page) func() {
	return func() {
		ui.switchToPage(QuitPage, func() { ui.switchToPage(backPage, nil) })
	}
}

func (ui *Ui) switchToPage(page Page, backFunction func()) {
	ui.pages.SwitchToPage(string(page))
	ui.backFunction = backFunction
}

func (ui *Ui) switchToPreviousPage() {
	if ui.backFunction != nil {
		ui.backFunction()
	}
}

func (ui *Ui) buildUi() *tview.Pages {
	mainMenuOptions := []*component.Option{
		component.NewListOption("New Game", "Starts a new game", 'n', ui.startGame),
		component.NewListOption("Quit", "Quits the application", 'q', ui.switchToQuitModal(MainMenuPage)),
	}

	gameMenuOptions := []*component.Option{
		component.NewListOption("Resume Game", "Resumes the current game", 'r', ui.resumeGame),
		component.NewListOption("Restart Level", "Restarts the current level", 'l', ui.restartLevel),
		component.NewListOption("Quit", "Quits the application", 'q', ui.switchToQuitModal(GameMenuPage)),
	}

	quitOptions := []*component.Option{
		component.NewModalOption("Yes", 'y', ui.exit),
		component.NewModalOption("No", 'n', ui.switchToPreviousPage),
	}

	winPageOptions := []*component.Option{
		component.NewModalOption("Yes!", 'y', ui.exit),
		component.NewModalOption("For the lack of other options, Yes!", 'y', ui.exit),
	}

	ui.gameView = newGameView(&[]*component.KeyEvent{
		component.NewKeyEvent(component.Escape, ui.switchToGameMenu),
		component.NewKeyEvent(component.ArrowUp, ui.sendGameCommand(command.NewMoveCommand(model.North))),
		component.NewKeyEvent(component.ArrowRight, ui.sendGameCommand(command.NewMoveCommand(model.East))),
		component.NewKeyEvent(component.ArrowDown, ui.sendGameCommand(command.NewMoveCommand(model.South))),
		component.NewKeyEvent(component.ArrowLeft, ui.sendGameCommand(command.NewMoveCommand(model.West))),
	})

	return tview.NewPages().
		AddPage(MainMenuPage, component.NewMenuView(&mainMenuOptions), true, true).
		AddPage(GameMenuPage, component.NewMenuView(&gameMenuOptions), true, false).
		AddPage(QuitPage, component.NewModal("Do you want to quit the application ?", &quitOptions), true, false).
		AddPage(PlayPage, ui.gameView, true, false).
		AddPage(WinPage, component.NewModal("You won!!1 Would you like to exit the application ?", &winPageOptions), true, false)
}

func (ui *Ui) sendGameCommand(command command.Command) func() {
	return func() {
		ui.getCommandChannel() <- command
	}
}

func (ui *Ui) getStateChannel() <-chan *game.GameState {
	ui.channelLock.Lock()
	defer ui.channelLock.Unlock()

	return ui.stateChannel
}

func (ui *Ui) getCommandChannel() chan<- command.Command {
	ui.channelLock.Lock()
	defer ui.channelLock.Unlock()

	return ui.commandChannel
}

func (ui *Ui) setGameChannels(stateChannel <-chan *game.GameState, commandChannel chan<- command.Command) {
	ui.channelLock.Lock()
	defer ui.channelLock.Unlock()

	if ui.stateChannel != nil && stateChannel != nil {
		panic("state channel is already set")
	}

	if ui.commandChannel != nil && commandChannel != nil {
		panic("command channel is already set")
	}

	ui.stateChannel = stateChannel
	ui.commandChannel = commandChannel
}

func (ui *Ui) listenForStateChanges() {
	for {
		channel := ui.getStateChannel()
		if channel == nil {
			return
		}

		state, more := <-channel
		if !more {
			return
		}

		if state == nil {
			ui.switchToPage(WinPage, nil)
		} else {
			ui.gameView.SetGameState(state)
			ui.application.Draw()
		}
	}
}

type Page string

const (
	MainMenuPage = "MainMenu"
	GameMenuPage = "GameMenu"
	QuitPage     = "Quit"
	PlayPage     = "Play"
	WinPage      = "Win"
)
