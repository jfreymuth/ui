package toolkit

import (
	"bufio"
	"image"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type TextArea struct {
	Editable       bool
	Theme          *Theme
	Font, font     draw.Font
	text           []string
	cx             int
	w, h           int
	cursor         cursor
	selectionStart cursor
	changedLine    int
	ll, sll, slll  int
	lastX, lastY   int
	state          byte
	scr            bool
	changed        bool
	popup          Menu
}

func NewTextArea() *TextArea {
	t := &TextArea{Theme: DefaultTheme, Font: DefaultTheme.Font("inputText"), h: -1, text: []string{""}, Editable: true}
	t.popup = *NewPopupMenu("")
	t.popup.AddItem("Cut", t.Cut)
	t.popup.AddItem("Copy", t.Copy)
	t.popup.AddItem("Paste", t.Paste)
	return t
}

type cursor struct {
	line, col int
}

func (t *TextArea) SetTheme(theme *Theme) {
	t.Theme = theme
	t.popup.SetTheme(theme)
}
func (t *TextArea) Text() string { return strings.Join(t.text, "\n") }
func (t *TextArea) SetText(text string) {
	t.text = strings.Split(strings.Replace(text, "\t", "    ", -1), "\n")
	t.textReplaced()
}

func (t *TextArea) Lines() []string {
	return append(make([]string, 0, len(t.text)), t.text...)
}

func (t *TextArea) SetLines(l []string) {
	if len(l) == 0 {
		t.text = []string{""}
	} else {
		t.text = append(t.text[:0], l...)
	}
	t.textReplaced()
}

func (t *TextArea) Append(text string) {
	t.text = append(t.text, strings.Split(strings.Replace(text, "\t", "    ", -1), "\n")...)
}

func (t *TextArea) SetTextFromReader(r io.Reader) error {
	t.text = nil
	return t.AppendFromReader(r)
}

func (t *TextArea) AppendFromReader(r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		t.text = append(t.text, strings.Replace(sc.Text(), "\t", "    ", -1))
	}
	t.text = append(t.text, strings.Replace(sc.Text(), "\t", "    ", -1))
	t.textReplaced()
	return sc.Err()
}

func (t *TextArea) WriteTextTo(w io.Writer) error {
	text := t.text
	if t.text[len(t.text)-1] == "" {
		text = t.text[:len(t.text)-1]
	}
	for _, l := range text {
		_, err := io.WriteString(w, l)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TextArea) ReplaceSelection(text string) {
	s1, _ := t.selection()
	t.insert(text)
	t.selectionStart = s1
}

func (t *TextArea) textReplaced() {
	t.cursor = cursor{len(t.text) - 1, len(t.text[len(t.text)-1])}
	t.selectionStart = t.cursor
	t.h = -1
	t.cx = -1
	t.changed = true
}

func (t *TextArea) Changed() bool {
	c := t.changed
	t.changed = false
	return c
}

func (t *TextArea) PreferredSize(fonts draw.FontLookup) (int, int) {
	t.measure(nil, fonts)
	w, h := t.w, t.h*len(t.text)
	if w < 200 {
		w = 200
	}
	if len(t.text) < 3 {
		h = t.h * 3
	}
	return w + 4, h + 4
}

func (t *TextArea) Update(g *draw.Buffer, state *ui.State) {
	state.DisableTabFocus()
	state.SetCursor(ui.CursorText)
	t.handleKeyEvents(state, g.FontLookup)
	t.hanldeMouseEvents(state, g.FontLookup)
	t.measure(state, g.FontLookup)
	w, h := g.Size()
	m := g.FontLookup.Metrics(t.font)

	g.Fill(draw.WH(w, h), t.Theme.Color("inputBackground"))
	if t.Editable && state.HasKeyboardFocus() {
		g.Outline(draw.WH(w, h), t.Theme.Color("border"))
	}
	s1, s2 := t.selection()

	x := 2
	y := 2 + s1.line*t.h
	cx, cy := -1, -1
	x += int(m.Advance(t.text[s1.line][:s1.col]))
	if state.HasKeyboardFocus() && s1 == t.cursor {
		if t.cx < 0 {
			t.cx = x
			t.scr = true
		}
		t.scroll(state, x)
		cx, cy = x, y
	}
	if s1.line == s2.line {
		adv := m.Advance(t.text[s1.line][s1.col:s2.col])
		if s1.col != s2.col {
			if state.HasKeyboardFocus() {
				g.Fill(draw.XYWH(x, y, int(adv), t.h), t.Theme.Color("selection"))
			} else {
				g.Fill(draw.XYWH(x, y, int(adv), t.h), t.Theme.Color("selectionInactive"))
			}
		}
		x += int(adv)
	} else {
		var color draw.Color
		if state.HasKeyboardFocus() {
			color = t.Theme.Color("selection")
		} else {
			color = t.Theme.Color("selectionInactive")
		}
		g.Fill(draw.XYXY(x, y, w-2, y+t.h), color)
		g.Fill(draw.XYWH(2, 2+(s1.line+1)*t.h, w-4, (s2.line-s1.line-1)*t.h), color)
		x, y = 2, 2+s2.line*t.h
		adv := m.Advance(t.text[s2.line][:s2.col])
		g.Fill(draw.XYWH(x, y, int(adv), t.h), color)
		x += int(adv)
	}
	if state.HasKeyboardFocus() && s2 == t.cursor {
		if t.cx < 0 {
			t.cx = int(x)
			t.scr = true
		}
		t.scroll(state, int(x))
		cx, cy = x, y
	}
	for i := range t.text {
		g.Text(draw.XYWH(2, 2+i*t.h, w-4, t.h), t.text[i], t.Theme.Color("inputText"), t.font)
	}
	if t.Editable && cx >= 0 && state.Blink() {
		g.Fill(draw.XYWH(cx-1, cy, 2, t.h), t.Theme.Color("inputText"))
	}
}

func (t *TextArea) hanldeMouseEvents(state *ui.State, fonts draw.FontLookup) {
	mouse := state.MousePos()
	drag, drop := state.DraggedContent()
	if drag, ok := drag.(string); ok {
		t.cursor = t.getCursor(fonts, mouse)
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
		c := t.getCursor(fonts, mouse)
		if t.state == tfDrag {
			if !t.inSelection(c) {
				state.InitiateDrag(t.SelectedText())
				t.insert("")
			}
		} else if t.state == tfIdle && state.ClickCount() == 1 && t.inSelection(c) {
			t.state = tfDrag
		} else {
			if t.lastX != mouse.X || t.lastY != mouse.Y {
				t.cursor = c
				t.cx = mouse.X
				t.scr = true
				if t.state == tfIdle {
					t.state = tfSelect
					t.selectionStart = t.cursor
				}
			} else if t.state == tfIdle {
				switch state.ClickCount() % 3 {
				case 1:
					t.selectionStart, t.cursor = c, c
				case 2:
					t.selectionStart.line, t.cursor.line = c.line, c.line
					t.selectionStart.col, t.cursor.col = text.FindWord(t.text[c.line], c.col)
					t.state = tfSelect
				case 0:
					t.selectLine()
					t.state = tfSelect
				}
			}
			state.SetBlink()
			t.lastX, t.lastY = mouse.X, mouse.Y
		}
	} else {
		if t.state == tfDrag {
			t.cursor = t.getCursor(fonts, mouse)
			t.selectionStart = t.cursor
			state.SetBlink()
		}
		t.state = tfIdle
	}
	if state.MouseButtonDown(ui.MouseRight) {
		c := t.getCursor(fonts, mouse)
		if !t.inSelection(c) && c != t.cursor && c != t.selectionStart {
			t.cursor, t.selectionStart = c, c
		}
		t.popup.OpenPopupMenu(mouse, state, fonts)
		state.InitiateDrag(ui.MenuDrag)
	}
}

func (t *TextArea) handleKeyEvents(state *ui.State, fonts draw.FontLookup) {
	if text := state.TextInput(); text != "" {
		t.insert(text)
		state.SetBlink()
	}
	for _, k := range state.KeyPresses() {
		switch k {
		case ui.KeyLeft:
			t.cursor = t.prev(state)
			t.cx = -1
		case ui.KeyRight:
			t.cursor = t.next(state)
			t.cx = -1
		case ui.KeyUp:
			if t.cursor.line > 0 {
				t.cursor.line--
				t.cursor.col = t.findPosition(fonts, t.cursor.line, t.cx)
			} else {
				t.cursor.col = 0
			}
			t.scr = true
		case ui.KeyDown:
			if t.cursor.line < len(t.text)-1 {
				t.cursor.line++
				t.cursor.col = t.findPosition(fonts, t.cursor.line, t.cx)
			} else {
				t.cursor.col = len(t.text[t.cursor.line])
			}
			t.scr = true
		case ui.KeyHome:
			t.cursor.col = 0
			t.cx = -1
		case ui.KeyEnd:
			t.cursor.col = len(t.text[t.cursor.line])
			t.cx = -1
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
			t.insert("\n")
		case ui.KeyTab:
			t.insert("\t")
		case ui.KeyMenu:
			t.popup.OpenPopupMenu(image.Pt(t.cx, (t.cursor.line+1)*t.h), state, fonts)
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

func (t *TextArea) SelectAll(state *ui.State) {
	t.selectionStart = cursor{0, 0}
	t.cursor = cursor{len(t.text) - 1, len(t.text[len(t.text)-1])}
}

func (t *TextArea) Cut(state *ui.State) {
	state.SetClipboardString(t.SelectedText())
	t.insert("")
}

func (t *TextArea) Copy(state *ui.State) {
	state.SetClipboardString(t.SelectedText())
}

func (t *TextArea) Paste(state *ui.State) {
	t.insert(state.ClipboardString())
}

func (t *TextArea) scroll(state *ui.State, x int) {
	if t.scr {
		state.RequestVisible(draw.XYWH(x+2-t.h*2, (t.cursor.line-1)*t.h+2, t.h*4, t.h*3))
		t.scr = false
	}
}

func (t *TextArea) measure(state *ui.State, fonts draw.FontLookup) {
	if t.h < 0 || t.Font != t.font {
		t.font = t.Font
		m := fonts.Metrics(t.font)
		t.h = m.LineHeight()
		t.w = 200
		t.slll = 0
		for i, l := range t.text {
			wf := int(m.Advance(l))
			w := int(wf)
			if w > t.w {
				t.slll = t.w
				t.sll = t.ll
				t.w = w
				t.ll = i
			} else if w > t.slll {
				t.slll = w
				t.sll = i
			}
		}
		t.changedLine = -1
		if state != nil {
			state.RequestUpdate()
		}
	} else if t.changedLine != -1 {
		m := fonts.Metrics(t.font)
		cll := int(m.Advance(t.text[t.changedLine]))
		if t.changedLine == t.ll {
			if cll >= t.slll {
				t.w = cll
			} else {
				t.h = -1
			}
		} else if t.changedLine == t.sll {
			if cll < t.sll {
				t.h = -1
			} else if cll > t.w {
				t.ll, t.sll = t.sll, t.ll
				t.slll = t.w
				t.w = cll
			}
		} else {
			if cll > t.slll {
				if cll > t.w {
					t.slll = t.w
					t.sll = t.ll
					t.w = cll
					t.ll = t.changedLine
				} else {
					t.slll = cll
					t.sll = t.changedLine
				}
			}
		}
		t.changedLine = -1
		t.measure(state, fonts)
	}
}

func (t *TextArea) lineChanged(l int) {
	if t.changedLine == -1 {
		t.changedLine = l
	} else if t.changedLine != l {
		t.h = -1
	}
}

func (t *TextArea) getCursor(fonts draw.FontLookup, p image.Point) cursor {
	line := (p.Y - 2) / t.h
	if line < 0 {
		return cursor{0, 0}
	} else if line >= len(t.text) {
		return cursor{len(t.text) - 1, len(t.text[len(t.text)-1])}
	}
	col := t.findPosition(fonts, line, p.X)
	return cursor{line, col}
}

func (t *TextArea) findPosition(fonts draw.FontLookup, line int, x int) int {
	return fonts.Metrics(t.font).Index(t.text[line], float32(x-2))
}

func (t *TextArea) selection() (cursor, cursor) {
	if t.selectionStart.line < t.cursor.line {
		return t.selectionStart, t.cursor
	}
	if t.selectionStart.line == t.cursor.line && t.selectionStart.col < t.cursor.col {
		return t.selectionStart, t.cursor
	}
	return t.cursor, t.selectionStart
}

func (t *TextArea) inSelection(c cursor) bool {
	s1, s2 := t.selection()
	if c.line < s1.line {
		return false
	} else if c.line == s1.line && c.col <= s1.col {
		return false
	} else if c.line > s2.line {
		return false
	} else if c.line == s2.line && c.col >= s2.col {
		return false
	}
	return true
}

func (t *TextArea) SelectedText() string {
	s1, s2 := t.selection()
	if s1.line == s2.line {
		return t.text[s1.line][s1.col:s2.col]
	}
	if s1.line+1 == s2.line {
		return t.text[s1.line][s1.col:] + "\n" + t.text[s2.line][:s2.col]
	}
	return t.text[s1.line][s1.col:] + "\n" + strings.Join(t.text[s1.line+1:s2.line], "\n") + "\n" + t.text[s2.line][:s2.col]
}

func (t *TextArea) next(state *ui.State) cursor {
	c := t.cursor
	line := t.text[c.line]
	if c.col == len(line) {
		if c.line == len(t.text)-1 {
			return c
		}
		return cursor{c.line + 1, 0}
	}
	if state.HasModifiers(ui.Control) {
		return cursor{c.line, text.NextWord(line, c.col)}
	} else {
		_, size := utf8.DecodeRuneInString(t.text[c.line][c.col:])
		return cursor{c.line, c.col + size}
	}
}

func (t *TextArea) prev(state *ui.State) cursor {
	c := t.cursor
	line := t.text[c.line]
	if c.col == 0 {
		if c.line == 0 {
			return c
		}
		return cursor{c.line - 1, len(t.text[c.line-1])}
	}
	if state.HasModifiers(ui.Control) {
		return cursor{c.line, text.PreviousWord(line, c.col)}
	} else {
		_, size := utf8.DecodeLastRuneInString(t.text[c.line][:c.col])
		return cursor{c.line, c.col - size}
	}
}

func (t *TextArea) selectLine() {
	t.selectionStart.col = 0
	if t.cursor.line == len(t.text)-1 {
		t.cursor.col = len(t.text[t.cursor.line])
	} else {
		t.cursor.col = 0
		t.cursor.line++
	}
}

func (t *TextArea) insert(s string) {
	if !t.Editable {
		return
	}
	s = strings.Replace(s, "\t", "    ", -1)
	s1, s2 := t.selection()
	if s == "" && s1 == s2 {
		return
	}
	if strings.Contains(s, "\n") {
		lines := strings.Split(s, "\n")
		ll := len(lines) - 1
		t.cursor.col = len(lines[ll])
		lines[ll] = lines[ll] + t.text[s2.line][s2.col:]
		t.text[s1.line] = t.text[s1.line][:s1.col] + lines[0]
		t.text = append(t.text[:s1.line+1], append(lines[1:], t.text[s2.line+1:]...)...)
		t.cursor.line = s1.line + len(lines) - 1
		t.h = -1
	} else {
		t.text[s1.line] = t.text[s1.line][:s1.col] + s + t.text[s2.line][s2.col:]
		t.text = append(t.text[:s1.line+1], t.text[s2.line+1:]...)
		t.cursor = cursor{s1.line, s1.col + len(s)}
		if s1.line == s2.line {
			t.lineChanged(s1.line)
		} else {
			t.h = -1
		}
	}
	t.changed = true
	t.selectionStart = t.cursor
	t.cx = -1
}

// Reader returns an io.Reader that will read from the contents of the text area.
// Changes to the text area's content after this method is called will not affect the returned Reader.
func (t *TextArea) Reader() io.Reader {
	return &reader{t.Lines(), 0}
}

type reader struct {
	lines []string
	pos   int
}

func (r *reader) Read(p []byte) (int, error) {
	if len(r.lines) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.lines[0][r.pos:])
	r.pos += n
	p = p[n:]
	if len(p) > 0 && r.pos == len(r.lines[0]) {
		r.lines = r.lines[1:]
		r.pos = 0
		p[0] = '\n'
		n++
	}
	return n, nil
}
