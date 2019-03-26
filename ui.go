package ui

import (
	"image"

	"github.com/jfreymuth/ui/draw"
)

type Component interface {
	PreferredSize(draw.FontLookup) (int, int)
	Update(*draw.Buffer, *State)
}

type Root interface {
	OpenDialog(Component)
	CloseDialog()
	OpenPopup(image.Rectangle, Component) Popup
	ClosePopups()
	HasPopups() bool
}

type Popup interface {
	Close()
	Closed() bool
}

type Modifier byte

const (
	Shift Modifier = 1 << iota
	Control
	Alt
	Super
	CapsLock
	NumLock
)

type MouseButton byte

const (
	MouseLeft MouseButton = 1 << iota
	MouseMiddle
	MouseRight
	MouseForward
	MouseBack
)

type Cursor byte

const (
	CursorNormal Cursor = iota
	CursorText
	CursorHand
	CursorCrosshair
	CursorDisabled
	CursorWait
	CursorWaitBackground
	CursorMove
	CursorResizeHorizontal
	CursorResizeVertical
	CursorResizeDiagonal
	CursorResizeDiagonal2
)

func HandleKeyboardShortcuts(state *State) {
	for _, k := range state.PeekKeyPresses() {
		if state.HasModifiers(Control) {
			switch k {
			case KeyA:
				if t, ok := state.KeyboardFocus().(interface{ SelectAll(*State) }); ok {
					t.SelectAll(state)
				}
			case KeyX:
				if t, ok := state.KeyboardFocus().(interface{ Cut(*State) }); ok {
					t.Cut(state)
				}
			case KeyC:
				if t, ok := state.KeyboardFocus().(interface{ Copy(*State) }); ok {
					t.Copy(state)
				}
			case KeyV:
				if t, ok := state.KeyboardFocus().(interface{ Paste(*State) }); ok {
					t.Paste(state)
				}
			}
		} else {
			switch k {
			case KeyEscape:
				state.ClosePopups()
			}
		}
	}
}

type menuDrag int

const MenuDrag menuDrag = 0
