package component

import "github.com/rivo/tview"

type MyBox struct {
	*tview.Box
}

func NewMyBox(title string) *MyBox {
	return &MyBox{Box: tview.NewBox().SetTitle(title).SetBorder(true)}
}