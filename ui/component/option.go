package component

type Option struct {
	label       string
	description string
	shortcut    rune
	action      func()
}

func NewListOption(label string, description string, shortcut rune, action func()) *Option {
	return newOption(label, description, shortcut, action)
}

func NewModalOption(label string, shortcut rune, action func()) *Option {
	return newOption(label, "", shortcut, action)
}

func newOption(label string, description string, shortcut rune, action func()) *Option {
	return &Option{
		label:       label,
		description: description,
		shortcut:    shortcut,
		action:      action,
	}
}
