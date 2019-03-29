package ui

import (
	"image"

	"github.com/jfreymuth/ui/draw"
)

// The State is a components connection to the application.
// It can be used to query input events and to request actions.
type State struct {
	current   Component
	disabled  bool
	hovered   bool            // true: mouse is inside the current component
	focusable bool            // true: current component has tried to receive keyboard events
	visible   image.Rectangle // visibility request, relative to bounds
	bounds    image.Rectangle // current component bounds

	mousePos      image.Point // mouse position, absolute
	hoveredC      Component   //
	grabbed       Component   // not nil: mouse was pressed and not yet released on component
	drag          interface{} // drag and drop, content
	drop          bool        // drag and drop, mouse just released
	focused       Component   // keyboard focus
	focusNext     bool
	lastFocusable Component
	mouseButtons  MouseButton // pressed mouse buttons
	clickButtons  MouseButton // mouse buttons released since last update
	clicks        int         // number of mouse clicks
	modifiers     Modifier
	scroll        image.Point // mouse wheel input
	textInput     string
	keyPresses    []Key
	root          Root
	cursor        Cursor
	windowTitle   string
	windowSize    image.Rectangle // always at (0,0)
	clipboard     string
	time          float32
	blink         bool

	// requests, set by component and read by backend
	update    bool
	animation bool
	refocus   bool
	quit      bool
}

// HasMouseFocus returns true if the current component receives mouse events, i.e. if it is hovered or grabbed.
func (s *State) HasMouseFocus() bool {
	return !s.disabled && s.drag == nil && (s.grabbed == s.current || (s.grabbed == nil && s.hovered))
}

// IsHovered returns true if the cursor is inside the current component.
func (s *State) IsHovered() bool {
	return !s.disabled && s.hovered && (s.grabbed == s.current || s.grabbed == nil)
}

// difference between HasMouseFocus and IsHovered: if the mouse is pressed and dragged outside the component,
// HasMouseFocus returns true, but IsHovered returns false

// DisableTabFocus disables focus cycling with the tab key for the current component.
func (s *State) DisableTabFocus() {
	if s.grabbed == s.current {
		s.focusNext = false
		s.focused = s.current
	}
	s.focusable = true
	if s.current == s.focused {
		s.focusNext = false
	}
}

// SetKeyboardFocus requests that the given component will receive keyboard events.
func (s *State) SetKeyboardFocus(c Component) {
	s.focused = c
	s.focusNext = false
	s.update = true
}

// SetKeyboardFocus requests that the next component will receive keyboard events.
func (s *State) FocusNext() {
	s.focusNext = true
}

// SetKeyboardFocus requests that the previous component will receive keyboard events.
func (s *State) FocusPrevious() {
	s.focused = s.lastFocusable
	s.update = true
}

// KeyboardFocus returns the component that currently receives keyboard events.
func (s *State) KeyboardFocus() Component {
	return s.focused
}

// HasKeyboardFocus returns true if the current component receives keyboard events.
func (s *State) HasKeyboardFocus() bool {
	if s.disabled {
		return false
	}
	if !s.focusable {
		s.focusable = true
		if s.current == s.focused {
			s.focusNext = false
			for i, k := range s.keyPresses {
				if k == KeyTab {
					if s.HasModifiers(Shift) {
						s.focused = s.lastFocusable
						s.update = true
					} else {
						s.focusNext = true
					}
					s.keyPresses = append(s.keyPresses[:i], s.keyPresses[i+1:]...)
					return false
				}
			}
		}
		if s.focused != s.current && (s.grabbed == s.current || s.focusNext) {
			s.focusNext = false
			s.focused = s.current
			s.update = true
			s.visible = draw.WH(s.bounds.Dx(), s.bounds.Dy())
		}
	}
	return s.focused == s.current
}

// MouseButtonDown returns true if a given mouse button is pressed.
func (s *State) MouseButtonDown(b MouseButton) bool {
	return s.HasMouseFocus() && s.mouseButtons&b == b
}

// MouseClick returns true if a given mouse button was clicked between the current and last update.
func (s *State) MouseClick(b MouseButton) bool {
	return s.HasMouseFocus() && s.clickButtons&b != 0
}

// ClickCount returns the number of consecutive mouse clicks.
func (s *State) ClickCount() int {
	return s.clicks
}

// HasModifiers returns true if the given modifiers are currently active.
func (s *State) HasModifiers(m Modifier) bool {
	return s.modifiers&m == m
}

// MousePos returns the position of the cursor relative to the current component.
func (s *State) MousePos() image.Point {
	return s.mousePos.Sub(s.bounds.Min)
}

// Scroll returns the amount of scrolling since the last update.
func (s *State) Scroll() image.Point {
	if !s.disabled && s.hovered {
		return s.scroll
	} else {
		return image.Point{}
	}
}

// ConsumeScroll notifies the State that the scroll amount has been used and should not be used by any other components.
func (s *State) ConsumeScroll() {
	if !s.disabled && s.hovered {
		s.scroll = image.Point{}
	}
}

// KeyPresses returns a list of key events that the current component should process.
func (s *State) KeyPresses() []Key {
	if s.HasKeyboardFocus() {
		k := s.keyPresses
		s.keyPresses = nil
		return k
	}
	return nil
}

// PeekKeyPresses returns a list of key events.
// Unlike KeyPresses, this method returns events even if they are not intended for the current component.
func (s *State) PeekKeyPresses() []Key {
	return s.keyPresses
}

// TextInput returns the string that would be generated by key inputs.
// Key presses that contributed to the text input will still appear in KeyPresses().
func (s *State) TextInput() string {
	if s.HasKeyboardFocus() {
		t := s.textInput
		s.textInput = ""
		return t
	} else {
		return ""
	}
}

// InitiateDrag starts a drag and drop gesture.
func (s *State) InitiateDrag(content interface{}) {
	if s.HasMouseFocus() {
		s.grabbed = nil
		s.drag = content
	}
}

// DraggedContent returns information abount a drag and drop gesture currently in progress.
// If there is no drag and drop gesture, or the mouse cursor is not above the current component, content will be nil.
// Otherwise, content will be the value passed to InitiateDrag.
// drop will be true if the mouse was released just before the current update.
func (s *State) DraggedContent() (content interface{}, drop bool) {
	if s.IsHovered() {
		return s.drag, s.drop
	}
	return nil, false
}

// Blink returns the current state of blinking elements.
// This is mostly indended for the cursor in text fields.
func (s *State) Blink() bool {
	return s.HasKeyboardFocus() && s.blink
}

// SetBlink requests that blinking elements should be visible.
func (s *State) SetBlink() {
	if s.HasKeyboardFocus() {
		s.blink = true
	}
}

func (s *State) SetRoot(r Root) {
	s.root = r
}

// OpenDialog displays the given component as a dialog.
// While a dialog is open, other components do not receive any events.
// Only one dialog may be open at a time.
func (s *State) OpenDialog(d Component) {
	if s.root != nil {
		s.root.OpenDialog(d)
		s.focused = d
		s.update = true
	}
}

// CloseDialog closes the currently open dialog, if any.
func (s *State) CloseDialog() {
	if s.root != nil {
		s.root.CloseDialog()
		s.update = true
	}
}

// OpenPopup displays the given component as a popup.
func (s *State) OpenPopup(bounds image.Rectangle, d Component) Popup {
	if s.root != nil {
		s.focused = d
		s.update = true
		return s.root.OpenPopup(bounds.Add(s.bounds.Min), d)
	}
	return nil
}

// ClosePopups closes all popups.
func (s *State) ClosePopups() {
	if s.root != nil {
		s.root.ClosePopups()
		s.update = true
	}
}

// HasPopups returns wether any popups are currently open.
func (s *State) HasPopups() bool {
	return s.root != nil && s.root.HasPopups()
}

// WindowBounds returns the bounds of the window relative to the current component's origin.
// This means the minimum point of the returned Rect will most likely be negative.
func (s *State) WindowBounds() image.Rectangle {
	return s.windowSize.Sub(s.bounds.Min)
}

// SetCursor sets the current compnents cursor style.
func (s *State) SetCursor(c Cursor) {
	if s.HasMouseFocus() {
		s.cursor = c
	}
}

// SetCursor sets the title of the window.
func (s *State) SetWindowTitle(title string) {
	s.windowTitle = title
}

// RequestVisible should be called to request that the current component should be scrolled, so that the given rectangle is visible.
func (s *State) RequestVisible(r image.Rectangle) {
	s.visible = r
}

// GetVisibilityRequest returns the rectangle passed to RequestVisible by a child of the current component.
// This method should be called after UpdateChild.
// The returned rectangle will be translated relative to the current component.
// If the second return value is false, RequestVisible was not called.
func (s *State) GetVisibilityRequest() (image.Rectangle, bool) {
	return s.visible, !s.visible.Empty()
}

// RequestUpdate requests that the ui should be updated again after the current update.
// This method will typically be called if an input event causes a change to the layout or to components that may have already been updated.
func (s *State) RequestUpdate() {
	s.update = true
}

// RequestAnimation requests that the ui should be updated again after a short amount of time.
func (s *State) RequestAnimation() {
	s.animation = true
}

// RequestRefocus requests that the component receiving mouse events should be determined again.
// This method should rarely be called by normal components.
func (s *State) RequestRefocus() {
	s.refocus = true
}

// AnimationSpeed returns the time since the last call to RequestAnimation in seconds.
func (s *State) AnimationSpeed() float32 {
	return s.time
}

// ClipboardString returns the contents of the clipboard.
func (s *State) ClipboardString() string {
	return s.clipboard
}

// SetClipboardString sets the contents of the clipboard.
func (s *State) SetClipboardString(c string) {
	s.clipboard = c
}

// Quit requests the application to close.
func (s *State) Quit() {
	s.quit = true
}

// DrawChild draws another component without passing on any input.
func (s *State) DrawChild(g *draw.Buffer, bounds image.Rectangle, c Component) {
	d := s.disabled
	this := s.current
	s.disabled = true
	s.current = c
	g.Push(bounds)
	c.Update(g, s)
	g.Pop()
	s.disabled = d
	s.current = this
}

// UpdateChild calls another component's Update method with a State that will deliver the correct events.
func (s *State) UpdateChild(g *draw.Buffer, bounds image.Rectangle, c Component) {
	g.Push(bounds)
	if s.disabled {
		c.Update(g, s)
		g.Pop()
		return
	}
	this := s.current
	h := s.hovered
	v := s.visible
	f := s.focusable
	s.visible = image.Rectangle{}
	s.hovered = h && s.mousePos.Sub(s.bounds.Min).In(bounds)
	if s.hovered {
		s.hoveredC = c
	}
	b := s.bounds
	s.bounds = bounds.Add(s.bounds.Min)
	s.current = c
	s.focusable = false
	if s.focused == s.current {
		s.focusNext = true
	}
	c.Update(g, s)
	if s.focusable {
		s.lastFocusable = c
	}
	s.current = this
	s.bounds = b
	s.hovered = h
	if s.visible.Empty() {
		s.visible = v
	} else {
		s.visible = s.visible.Add(bounds.Min)
	}
	s.focusable = f
	g.Pop()
}
