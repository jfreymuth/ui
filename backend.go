package ui

import (
	"image"
	"time"

	"github.com/jfreymuth/ui/draw"
)

type BackendState struct {
	State
	last time.Time
}

func (s *BackendState) ResetEvents() {
	s.scroll = image.Point{}
	s.visible = image.Rectangle{}
	s.textInput = ""
	for _, k := range s.keyPresses {
		if k == KeyTab {
			s.focusNext = true
		}
	}
	s.keyPresses = s.keyPresses[:0]
	s.clicks = 0
	s.clickButtons = 0
	if s.drop {
		s.drag = nil
		s.drop = false
	}
	s.time = 0
}

func (s *BackendState) SetModifiers(m Modifier)       { s.modifiers = m }
func (s *BackendState) SetMousePosition(x, y int)     { s.mousePos = image.Pt(x, y) }
func (s *BackendState) SetMouseButtons(b MouseButton) { s.mouseButtons = b }
func (s *BackendState) SetMouseClicks(clicks int)     { s.clicks = clicks }
func (s *BackendState) SetHovered(h bool)             { s.hovered = h }
func (s *BackendState) AddScroll(x, y int)            { s.scroll = s.scroll.Add(image.Pt(x, y)) }
func (s *BackendState) AddKeyPress(k Key)             { s.keyPresses = append(s.keyPresses, k) }
func (s *BackendState) AddTextInput(text string)      { s.textInput += text }
func (s *BackendState) SetBlink(b bool)               { s.blink = b }
func (s *BackendState) SetWindowSize(w, h int)        { s.windowSize = draw.WH(w, h) }

func (s *BackendState) GrabMouse() {
	s.grabbed = s.hoveredC
}
func (s *BackendState) ReleaseMouse(b MouseButton) {
	if s.grabbed == s.hoveredC {
		s.clickButtons = b
	}
	s.grabbed = nil
	if s.drag != nil {
		s.drop = true
	}
}

func (state *BackendState) ResetRequests() {
	now := time.Now()
	if state.animation {
		state.time = float32(now.Sub(state.last).Seconds())
	}
	state.last = now
	state.update = false
	state.animation = false
	state.refocus = false
	state.cursor = CursorNormal
}

func (s *BackendState) Cursor() Cursor           { return s.cursor }
func (s *BackendState) WindowTitle() string      { return s.windowTitle }
func (s *BackendState) UpdateRequested() bool    { return s.update || s.focusNext }
func (s *BackendState) AnimationRequested() bool { return s.animation }
func (s *BackendState) RefocusRequested() bool   { return s.refocus }
func (s *BackendState) QuitRequested() bool      { return s.quit }
