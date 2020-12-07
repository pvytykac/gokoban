package model

import (
	"testing"
)

// rigorous testing ensues lol

func TestIsSolved(t *testing.T) {
	level := NewLevel(5, 3, NewPlayer(NewPosition(1, 1), North), &[][]Tile{}, &[]*Position{})

	got := level.IsSolved()
	if got == false {
		t.Errorf("solved %t; wanted: true", got)
	}
}
