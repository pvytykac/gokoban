package game

import (
	"gokoban/game/command"
	"gokoban/game/model"
	"sync"
	"time"
)

type Controller struct {
	tickerLock     *sync.Mutex
	commandLock    *sync.Mutex
	stateLock      *sync.Mutex
	gameLock       *sync.Mutex
	game           *model.Game
	ticker         *time.Ticker
	stateChannel   chan *GameState
	commandChannel chan command.Command
}

func NewController() *Controller {
	return &Controller{
		tickerLock:     &sync.Mutex{},
		commandLock:    &sync.Mutex{},
		stateLock:      &sync.Mutex{},
		gameLock:       &sync.Mutex{},
		game:           nil,
		ticker:         nil,
		stateChannel:   nil,
		commandChannel: nil,
	}
}

func (controller *Controller) Start() (<-chan *GameState, chan<- command.Command) {
	if controller.IsStopped() {
		controller.setGame(model.LoadGameFiles())
		controller.startListeners()
	}

	return controller.getStateChannel(), controller.getCommandChannel()
}

func (controller *Controller) Stop() {
	if controller.IsStarted() {
		controller.setGame(nil)
		controller.stopListeners()
	}
}

func (controller *Controller) Pause() {
	if controller.IsRunning() {
		controller.stopListeners()
	}
}

func (controller *Controller) Resume() (<-chan *GameState, chan<- command.Command) {
	if controller.IsPaused() {
		controller.startListeners()
	}

	return controller.getStateChannel(), controller.getCommandChannel()
}

func (controller *Controller) RestartLevel() {
	controller.gameLock.Lock()


	restarted := model.LoadLevel(controller.game.LevelIndex)
	restarted.Next = controller.game.CurrentLevel.Next
	controller.game.CurrentLevel = restarted

	controller.gameLock.Unlock()
	controller.updateState()
}

func (controller *Controller) IsStarted() bool {
	return controller.getGame() != nil
}

func (controller *Controller) IsStopped() bool {
	return !controller.IsStarted()
}

func (controller *Controller) IsPaused() bool {
	return controller.IsStarted() && controller.getTicker() == nil
}

func (controller *Controller) IsRunning() bool {
	return controller.IsStarted() && controller.getTicker() != nil
}

func (controller *Controller) startListeners() {
	controller.setTicker(time.NewTicker(1 * time.Second))
	controller.setStateChannel(make(chan *GameState, 2))
	controller.setCommandChannel(make(chan command.Command, 2))

	go controller.startTickListener()
	go controller.startCommandListener()

	controller.updateState()
}

func (controller *Controller) stopListeners() {
	controller.getTicker().Stop()
	controller.setTicker(nil)

	close(controller.getStateChannel())
	controller.setStateChannel(nil)

	close(controller.getCommandChannel())
	controller.setCommandChannel(nil)
}

func (controller *Controller) updateState() {
	controller.gameLock.Lock()
	defer controller.gameLock.Unlock()

	channel := controller.getStateChannel()
	game := controller.game

	if channel != nil && game != nil {
		channel <- NewGameRepresentation(game)
		if game.CurrentLevel.IsSolved() {
			if game.HasNextLevel() {
				game.NextLevel()
			} else {
				channel <- nil
			}
		}
	}
}

func (controller *Controller) getGame() *model.Game {
	controller.gameLock.Lock()
	defer controller.gameLock.Unlock()

	return controller.game
}

func (controller *Controller) setGame(game *model.Game) {
	controller.gameLock.Lock()
	defer controller.gameLock.Unlock()

	if controller.game != nil && game != nil {
		panic("game is already set")
	}

	controller.game = game
}

func (controller *Controller) getTicker() *time.Ticker {
	controller.tickerLock.Lock()
	defer controller.tickerLock.Unlock()

	return controller.ticker
}

func (controller *Controller) setTicker(ticker *time.Ticker) {
	controller.tickerLock.Lock()
	defer controller.tickerLock.Unlock()

	if controller.ticker != nil && ticker != nil {
		panic("ticker is already set")
	}

	controller.ticker = ticker
}

func (controller *Controller) getCommandChannel() chan command.Command {
	controller.commandLock.Lock()
	defer controller.commandLock.Unlock()

	return controller.commandChannel
}

func (controller *Controller) setCommandChannel(commandChannel chan command.Command) {
	controller.commandLock.Lock()
	defer controller.commandLock.Unlock()

	if controller.commandChannel != nil && commandChannel != nil{
		panic("command channel is already set")
	}

	controller.commandChannel = commandChannel
}

func (controller *Controller) getStateChannel() chan *GameState {
	controller.stateLock.Lock()
	defer controller.stateLock.Unlock()

	return controller.stateChannel
}

func (controller *Controller) setStateChannel(stateChannel chan *GameState) {
	controller.stateLock.Lock()
	defer controller.stateLock.Unlock()

	if controller.stateChannel != nil && stateChannel != nil {
		panic("state channel is already set")
	}

	controller.stateChannel = stateChannel
}

func (controller *Controller) startTickListener() {
	for {
		ticker := controller.getTicker()
		if ticker == nil {
			return
		}

		_, more := <-ticker.C
		if !more {
			return
		}

		controller.game.CurrentLevel.IncrementTime(1)
		controller.updateState()
	}
}

func (controller *Controller) startCommandListener() {
	for {
		channel := controller.getCommandChannel()
		if channel == nil {
			return
		}

		cmd, more := <-channel
		if !more {
			return
		}

		game := controller.getGame()
		if game != nil && cmd.Apply(game.CurrentLevel) {
			controller.updateState()
		}
	}
}
