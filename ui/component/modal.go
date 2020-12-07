package component

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ModalView struct {
	*tview.Modal
	options *[]*Option
}

func NewModal(text string, options *[]*Option) *ModalView {
	size := len(*options)
	buttons := make([]string, size)

	for ix, option := range *options {
		buttons[ix] = fmt.Sprintf("%s (%c)", option.label, option.shortcut)
	}

	modal := tview.NewModal().
		SetText(text).
		AddButtons(buttons).
		SetDoneFunc(onButtonSubmit(options))

	view := &ModalView{Modal: modal, options: options}

	modal.Box.SetInputCapture(view.onKeyPress)

	return view
}

func (view *ModalView) onKeyPress(event *tcell.EventKey) *tcell.EventKey {
	for _, option := range *view.options {
		if option.shortcut == event.Rune() {
			option.action()
			return nil
		}
	}

	return event
}

func onButtonSubmit(options *[]*Option) func(int, string) {
	return func(ix int, label string) {
		for _, option := range *options {
			if fmt.Sprintf("%s (%c)", option.label, option.shortcut) == label {
				option.action()
			}
		}
	}
}
