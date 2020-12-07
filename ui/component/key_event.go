package component

import "github.com/gdamore/tcell/v2"

type KeyEvent struct {
	key    Key
	action func()
}

func NewKeyEvent(key Key, action func()) *KeyEvent {
	return &KeyEvent{key: key, action: action}
}

func (ke *KeyEvent) AppliesTo(event *tcell.EventKey) bool {
	rune := ke.key.ToRune()
	key := ke.key.ToKey()
	return (rune != 0 && rune == event.Rune()) || (key != 0 && key == event.Key())
}

func (ke *KeyEvent) Execute() {
	ke.action()
}

type Key string

const (
	Escape     Key = "Esc"
	Enter      Key = "Enter"
	ArrowLeft  Key = "Arrow Left"
	ArrowRight Key = "Arrow Right"
	ArrowUp    Key = "Arrow Up"
	ArrowDown  Key = "Arrow Down"
	Q          Key = "q"
	W          Key = "w"
	S          Key = "s"
	A          Key = "a"
	D          Key = "d"
	Z          Key = "z"

)

func (key Key) ToRune() rune {
	switch key {
	case Q:
		return 'q'
	case W:
		return 'w'
	case S:
		return 's'
	case A:
		return 'a'
	case D:
		return 'd'
	case Z:
		return 'z'
	}

	return rune(0)
}

func (key Key) ToKey() tcell.Key {
	switch key {
	case Escape:
		return tcell.KeyEscape
	case Enter:
		return tcell.KeyEnter
	case ArrowLeft:
		return tcell.KeyLeft
	case ArrowRight:
		return tcell.KeyRight
	case ArrowUp:
		return tcell.KeyUp
	case ArrowDown:
		return tcell.KeyDown
	}

	return tcell.Key(0)
}
