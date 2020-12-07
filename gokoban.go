package main

import (
	"gokoban/game"
	"gokoban/ui"
)

func main() {
	controller := game.NewController()
	ui.NewUi(controller).Start()
}
