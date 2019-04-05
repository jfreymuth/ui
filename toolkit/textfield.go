package toolkit

import (
	"unicode/utf8"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type TextField struct {
	Editable       bool
	Action         func(*ui.State, string)
	MinWidth       int
	Theme          *Theme
	Text           string
	text           text.Text
	cursor         int
	selectionStart int
	lastX          int
	state          byte
	anim           float32
}

func NewTextField() *TextField {
	return &TextField{Theme: DefaultTheme, Editable: true, MinWidth: 100}
}

func (t *TextField) SetTheme(theme *Theme) { t.Theme = theme }

func (t *TextField) SelectedText() string {
	s1, s2 := t.selection()
	return t.Text[s1:s2]
}

func (t *TextField) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := t.text.Size(t.Text, t.Theme.Font("inputText"), fonts)
	if w+6 < t.MinWidth {
		w = t.MinWidth - 6
	}
	return w + 6, h + 6
}

func (t *TextField) Update(g *draw.Buffer, state *ui.State) {
	m := g.FontLookup.Metrics(t.Theme.Font("inputText"))
	state.SetCursor(ui.CursorText)
	t.handleKeyEvents(state)
	t.handleMouseEvents(state, m)
	w, h := g.Size()
	_, th := t.text.Size(t.Text, t.Theme.Font("inputText"), g.FontLookup)

	animate(state, &t.anim, 8, t.Editable && state.HasKeyboardFocus())
	anim := int(t.anim * float32(w))
	g.Fill(draw.XYXY(0, (h-th)/2-3, anim, (h+th)/2+3), t.Theme.Color("inputBackground"))
	line := w - 6 - anim
	if line > 0 {
		g.Fill(draw.XYWH(3, (h+th)/2, line, 1), t.Theme.Color("inputText"))
	}
	x, y := float32(3), (h-th)/2
	s1, s2 := t.selection()
	x += m.Advance(t.Text[:s1])
	var cx int
	if s1 == t.cursor {
		cx = int(x)
	}
	if s1 != s2 {
		adv := m.Advance(t.Text[s1:s2])
		if state.HasKeyboardFocus() {
			g.Fill(draw.XYWH(int(x), y, int(adv), th), t.Theme.Color("selection"))
		} else {
			g.Fill(draw.XYWH(int(x), y, int(adv), th), t.Theme.Color("selectionInactive"))
		}
		x += adv
	}
	if s2 == t.cursor {
		cx = int(x)
	}
	if state.Blink() {
		g.Fill(draw.XYWH(cx-1, y, 2, th), t.Theme.Color("inputText"))
	}
	t.text.DrawLeft(g, draw.XYXY(3, y, int(x), y+th), t.Text, t.Theme.Font("inputText"), t.Theme.Color("inputText"))
}

func (t *TextField) handleMouseEvents(state *ui.State, m draw.FontMetrics) {
	mx := state.MousePos().X
	drag, drop := state.DraggedContent()
	if drag, ok := drag.(string); ok {
		t.cursor = m.Index(t.Text, float32(mx-3))
		t.selectionStart = t.cursor
		state.SetBlink()
		if drop {
			t.insert(drag)
		}
		t.state = tfIdle
		return
	}
	if !state.HasMouseFocus() {
		t.state = tfIdle
		return
	}
	if state.MouseButtonDown(ui.MouseLeft) {
		c := m.Index(t.Text, float32(mx-3))
		if t.state == tfDrag {
			if !t.inSelection(c) {
				state.InitiateDrag(t.SelectedText())
				t.insert("")
			}
		} else if t.state == tfIdle && state.ClickCount() == 1 && t.inSelection(c) {
			t.state = tfDrag
		} else {
			if t.lastX != mx {
				t.cursor = c
				if t.state == tfIdle {
					t.state = tfSelect
					t.selectionStart = t.cursor
				}
			} else if t.state == tfIdle {
				t.selectionStart = t.cursor
				t.state = tfSelect
				switch state.ClickCount() % 3 {
				case 1:
					t.selectionStart, t.cursor = c, c
				case 2:
					t.selectionStart, t.cursor = text.FindWord(t.Text, t.cursor)
					t.state = tfSelect
				case 0:
					t.SelectAll(state)
					t.state = tfSelect
				}
			}
			state.SetBlink()
			t.lastX = mx
		}
	} else {
		if t.state == tfDrag {
			t.cursor = m.Index(t.Text, float32(mx-3))
			t.selectionStart = t.cursor
			state.SetBlink()
		}
		t.state = tfIdle
	}
}

func (t *TextField) inSelection(c int) bool {
	s1, s2 := t.selection()
	return c > s1 && c < s2
}

func (t *TextField) handleKeyEvents(state *ui.State) {
	if text := state.TextInput(); text != "" {
		t.insert(text)
		state.SetBlink()
	}
	for _, k := range state.KeyPresses() {
		switch k {
		case ui.KeyLeft:
			t.cursor = t.prev(state)
		case ui.KeyRight:
			t.cursor = t.next(state)
		case ui.KeyHome:
			t.cursor = 0
		case ui.KeyEnd:
			t.cursor = len(t.Text)
		case ui.KeyBackspace:
			if t.Editable && t.cursor == t.selectionStart {
				t.cursor = t.prev(state)
			}
			t.insert("")
		case ui.KeyDelete:
			if t.Editable && t.cursor == t.selectionStart {
				t.cursor = t.next(state)
			}
			t.insert("")
		case ui.KeyEnter:
			t.TriggerAction(state)
			continue
		default:
			continue
		}
		if !state.HasModifiers(ui.Shift) {
			t.selectionStart = t.cursor
		}
		state.SetBlink()
	}
}

func (t *TextField) SelectAll(state *ui.State) {
	t.selectionStart = 0
	t.cursor = len(t.Text)
}

func (t *TextField) Cut(state *ui.State) {
	state.SetClipboardString(t.SelectedText())
	t.insert("")
}

func (t *TextField) Copy(state *ui.State) {
	state.SetClipboardString(t.SelectedText())
}

func (t *TextField) Paste(state *ui.State) {
	t.insert(state.ClipboardString())
}

func (t *TextField) TriggerAction(state *ui.State) {
	if t.Action != nil {
		t.Action(state, t.Text)
	}
}

func (t *TextField) selection() (int, int) {
	if t.selectionStart > len(t.Text) {
		t.selectionStart = len(t.Text)
	}
	if t.cursor > len(t.Text) {
		t.cursor = len(t.Text)
	}
	if t.selectionStart < t.cursor {
		return t.selectionStart, t.cursor
	}
	return t.cursor, t.selectionStart
}

func (t *TextField) next(state *ui.State) int {
	if t.cursor > len(t.Text) {
		return len(t.Text)
	}
	if state.HasModifiers(ui.Control) {
		return text.NextWord(t.Text, t.cursor)
	} else {
		_, size := utf8.DecodeRuneInString(t.Text[t.cursor:])
		return t.cursor + size
	}
}

func (t *TextField) prev(state *ui.State) int {
	if t.cursor > len(t.Text) {
		t.cursor = len(t.Text)
	}
	if state.HasModifiers(ui.Control) {
		return text.PreviousWord(t.Text, t.cursor)
	} else {
		_, size := utf8.DecodeLastRuneInString(t.Text[:t.cursor])
		return t.cursor - size
	}
}

func (t *TextField) insert(s string) {
	if !t.Editable {
		return
	}
	s1, s2 := t.selection()
	t.Text = t.Text[:s1] + s + t.Text[s2:]
	t.cursor = s1 + len(s)
	t.selectionStart = t.cursor
}

const (
	tfIdle = iota
	tfSelect
	tfDrag
)
