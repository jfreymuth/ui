package sdl

import (
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/impl/gldraw"
	"github.com/veandco/go-sdl2/sdl"
)

type Options struct {
	_ [0]byte
	// Title is the window's title.
	Title string
	// Width and Height set the window's size.
	// If either is 0, it will be replaced by the root component's preferred size.
	Width, Height int
	// Root is the root component of the window.
	// Must not be nil
	Root ui.Component
	// FontLookup should create a font lookup for the specified DPI setting. It will only be called once.
	// Must not be nil
	FontLookup func(dpi float32) gldraw.FontLookup
	//
	IconLookup gldraw.IconLookup
	// Init will be called once after the window is created. It can be used for any setup that requires access to a *ui.State.
	Init func(*ui.State)
	// Update will be called every time the application is updated.
	Update func(*ui.State)
	// Close
	Close func(*ui.State)
	// SDLInit will be called after the window is created, but before it is shown.
	// Can be used for any SDL-specific initialisation.
	SDLInit func(*sdl.Window)
}

// Show opens a window and blocks until it is closed.
func Show(opt Options) {
	runtime.LockOSThread()

	if opt.Root == nil {
		panic("ui: Root must not be nil")
	}
	if opt.FontLookup == nil {
		panic("ui: FontLookup must not be nil")
	}

	if opt.Close == nil {
		opt.Close = (*ui.State).Quit
	}

	sdl.Init(sdl.INIT_VIDEO)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	win, _ := sdl.CreateWindow(opt.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 0, 0, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE|sdl.WINDOW_HIDDEN)
	ctx, _ := win.GLCreateContext()
	win.GLMakeCurrent(ctx)
	sdl.GLSetSwapInterval(1)
	gl.InitWithProcAddrFunc(sdl.GLGetProcAddress)

	di, _ := win.GetDisplayIndex()
	dpi, _, _, err := sdl.GetDisplayDPI(di)
	if err != nil {
		dpi = 96
	}
	fonts := opt.FontLookup(dpi)
	var c gldraw.Context
	c.Init(fonts)
	c.SetIconLookup(opt.IconLookup)
	w, h := opt.Root.PreferredSize(fonts)
	if opt.Width != 0 {
		w = opt.Width
	}
	if opt.Height != 0 {
		h = opt.Height
	}
	if r, err := sdl.GetDisplayBounds(di); err == nil {
		if w == 0 {
			w = int(r.W) / 2
		} else if w > int(r.W) {
			w = int(r.W) * 7 / 8
		}
		if h == 0 {
			h = int(r.H) / 2
		} else if h > int(r.H) {
			h = int(r.H) * 7 / 8
		}
	}
	win.SetSize(int32(w), int32(h))
	if opt.SDLInit != nil {
		opt.SDLInit(win)
	}
	win.Show()

	state := &ui.BackendState{}
	state.SetWindowTitle(opt.Title)
	var g draw.Buffer
	g.FontLookup = fonts
	var cursor ui.Cursor
	cursorCache := make(map[ui.Cursor]*sdl.Cursor)
	clipboard, _ := sdl.GetClipboardText()
	state.SetClipboardString(clipboard)

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			sdl.PushEvent(&sdl.UserEvent{Type: sdl.USEREVENT, Code: 1})
			time.Sleep(500 * time.Millisecond)
			sdl.PushEvent(&sdl.UserEvent{Type: sdl.USEREVENT, Code: 2})
		}
	}()

	var grabButton uint8
	for !state.QuitRequested() {
		state.ResetEvents()
		var e sdl.Event
		if state.AnimationRequested() {
			e = sdl.PollEvent()
		} else {
			e = sdl.WaitEvent()
		}
		quitEvent := false
		for e != nil {
			switch e := e.(type) {
			case *sdl.QuitEvent:
				quitEvent = true
			case *sdl.MouseMotionEvent:
				state.SetMousePosition(int(e.X), int(e.Y))
			case *sdl.MouseWheelEvent:
				state.AddScroll(int(e.X), int(e.Y))
			case *sdl.MouseButtonEvent:
				if e.State == sdl.PRESSED {
					if grabButton == 0 {
						grabButton = e.Button
						sdl.CaptureMouse(true)
						state.GrabMouse()
					}
				} else if e.State == sdl.RELEASED {
					if grabButton == e.Button {
						grabButton = 0
						sdl.CaptureMouse(false)
						state.ReleaseMouse(ui.MouseButton(e.Button))
					}
				}
				state.SetMouseClicks(getClicks(e))
			case *sdl.KeyboardEvent:
				if e.State == sdl.PRESSED {
					state.AddKeyPress(ui.Key(e.Keysym.Scancode))
				}
			case *sdl.TextInputEvent:
				for i, c := range e.Text {
					if c == 0 {
						state.AddTextInput(string(e.Text[:i]))
						break
					}
				}
			case *sdl.WindowEvent:
				switch e.Event {
				case sdl.WINDOWEVENT_ENTER:
					state.SetHovered(true)
				case sdl.WINDOWEVENT_LEAVE:
					state.SetHovered(false)
				case sdl.WINDOWEVENT_SIZE_CHANGED:
					state.SetWindowSize(int(e.Data1), int(e.Data2))
				default:
				}
			case *sdl.UserEvent:
				switch e.Code {
				case 1:
					state.SetBlink(true)
				case 2:
					state.SetBlink(false)
				case 3:
					(<-funcs)(&state.State)
				}
			default:
			}
			e = sdl.PollEvent()
		}
		_, _, mb := sdl.GetMouseState()
		state.SetMouseButtons(ui.MouseButton(mb))
		state.SetModifiers(translateModifiers(sdl.GetModState()))

		w, h := win.GLGetDrawableSize()

		g.Reset(int(w), int(h))
		state.ResetRequests()
		if opt.Init == nil && opt.Update != nil {
			opt.Update(&state.State)
		}
		state.UpdateChild(&g, draw.WH(int(w), int(h)), opt.Root)
		if opt.Init != nil {
			// Call Init after the first update, so the state has it's root set correctly.
			// This is important if Init wants to show a dialog.
			opt.Init(&state.State)
			opt.Init = nil
		}
		if quitEvent {
			opt.Close(&state.State)
		}

		if state.RefocusRequested() {
			state.ReleaseMouse(0)
			state.UpdateChild(&g, draw.WH(int(w), int(h)), opt.Root)
			state.GrabMouse()
		}

		if state.UpdateRequested() {
			state.ResetEvents()
			g.Reset(int(w), int(h))
			state.UpdateChild(&g, draw.WH(int(w), int(h)), opt.Root)
			if state.UpdateRequested() {
				// If an application requests three updates in a row, wait one frame to prevent an infinite loop.
				state.RequestAnimation()
			}
		}

		g.Pop()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		c.Draw(int(w), int(h), g.All)

		if state.Cursor() != cursor {
			cursor = state.Cursor()
			setCursor(cursor, cursorCache)
		}
		if state.WindowTitle() != opt.Title {
			opt.Title = state.WindowTitle()
			win.SetTitle(opt.Title)
		}
		if state.ClipboardString() != clipboard {
			clipboard = state.ClipboardString()
			sdl.SetClipboardText(clipboard)
		} else {
			clipboard, _ = sdl.GetClipboardText()
			state.SetClipboardString(clipboard)
		}

		win.GLSwap()
	}
}

// Do queues a function for execution on the ui goroutine.
// Do should not be called before Show.
func Do(f func(*ui.State)) {
	funcs <- f
	sdl.PushEvent(&sdl.UserEvent{Type: sdl.USEREVENT, Code: 3})
}

var funcs = make(chan func(*ui.State), 1)

func translateModifiers(m sdl.Keymod) ui.Modifier {
	var mod ui.Modifier
	if m&sdl.KMOD_SHIFT != 0 {
		mod |= ui.Shift
	}
	if m&sdl.KMOD_CTRL != 0 {
		mod |= ui.Control
	}
	if m&sdl.KMOD_ALT != 0 {
		mod |= ui.Alt
	}
	if m&sdl.KMOD_GUI != 0 {
		mod |= ui.Super
	}
	if m&sdl.KMOD_CAPS != 0 {
		mod |= ui.CapsLock
	}
	if m&sdl.KMOD_NUM != 0 {
		mod |= ui.NumLock
	}
	return mod
}

func setCursor(c ui.Cursor, cache map[ui.Cursor]*sdl.Cursor) {
	if sdlc, ok := cache[c]; ok {
		sdl.SetCursor(sdlc)
	} else {
		switch c {
		case ui.CursorNormal:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)
		case ui.CursorText:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM)
		case ui.CursorHand:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND)
		case ui.CursorCrosshair:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_CROSSHAIR)
		case ui.CursorDisabled:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_NO)
		case ui.CursorWait:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_WAIT)
		case ui.CursorWaitBackground:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_WAITARROW)
		case ui.CursorMove:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZEALL)
		case ui.CursorResizeHorizontal:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZEWE)
		case ui.CursorResizeVertical:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENS)
		case ui.CursorResizeDiagonal:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENWSE)
		case ui.CursorResizeDiagonal2:
			sdlc = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENESW)
		default:
			return
		}
		cache[c] = sdlc
		sdl.SetCursor(sdlc)
	}
}

func getClicks(e *sdl.MouseButtonEvent) int {
	return int((*(*[2]uint8)(unsafe.Pointer(&e.State)))[1])
}
