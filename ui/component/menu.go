package component

import (
	"github.com/rivo/tview"
)

type MenuView struct {
	*tview.List
	options *[]*Option
}

func NewMenuView(options *[]*Option) *MenuView {
	list := tview.NewList()

	for _, option := range *options {
		list = list.AddItem(option.label, option.description, option.shortcut, option.action)
	}

	return &MenuView{List: list, options: options}
}
