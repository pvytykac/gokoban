package model

type Game struct {
	LevelIndex           int
	CurrentLevel         *Level
	time                 int64
}

func NewGame() *Game {
	return &Game{}
}

func (game *Game) AddLevel(level *Level) {
	if game.CurrentLevel == nil {
		game.CurrentLevel = level
		game.LevelIndex = 1
	} else {
		game.CurrentLevel.setNextLevel(level)
	}
}

func (game *Game) HasNextLevel() bool {
	return game.CurrentLevel.Next != nil
}

func (game *Game) NextLevel() {
	game.IncrementTime(game.CurrentLevel.time)
	game.CurrentLevel = game.CurrentLevel.Next
	game.LevelIndex++
}

func (game *Game) GetLevelCount() int {
	return game.GetRemainingLevelCount() + game.LevelIndex
}

func (game *Game) GetRemainingLevelCount() int {
	if game.CurrentLevel == nil {
		return 0
	}

	level := game.CurrentLevel
	cnt := 0
	for {
		if level.Next == nil {
			break
		}

		cnt++
		level = level.Next
	}

	return cnt
}

func (game *Game) IncrementTime(inc int64) {
	game.time += inc
}

func (game *Game) GetElapsedTime() int64 {
	return game.time + game.CurrentLevel.time
}
