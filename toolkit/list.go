package toolkit

import (
	"strings"
	"time"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type List struct {
	Theme    *Theme
	Items    []ListItem
	Selected int
	Changed  func(*ui.State, ListItem)
	Action   func(*ui.State, ListItem)

	grab    bool
	search  string
	searchT time.Time
}

type ListItem struct {
	Icon string
	Text string
	text text.Text
}

func NewList() *List {
	return &List{Theme: DefaultTheme}
}

func (l *List) SetTheme(theme *Theme) { l.Theme = theme }

func (l *List) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := 0, 0
	font := l.Theme.Font("text")
	for i := range l.Items {
		it := &l.Items[i]
		iw, ih := it.text.SizeIcon(it.Text, font, it.Icon, 3, fonts)
		if iw > w {
			w = iw
		}
		if ih > h {
			h = ih
		}
	}
	return w, h * len(l.Items)
}

func (l *List) Update(g *draw.Buffer, state *ui.State) {
	w, _ := g.Size()

	if len(l.Items) == 0 {
		return
	}
	_, h := l.Items[0].text.Size(l.Items[0].Text, l.Theme.Font("text"), g.FontLookup)

	mouse := state.MousePos()
	if state.MouseButtonDown(ui.MouseLeft) {
		if !l.grab {
			l.grab = true
			sel := mouse.Y / h
			if sel >= 0 && sel < len(l.Items) {
				if sel == l.Selected && state.ClickCount() == 2 {
					if l.Action != nil {
						l.Action(state, l.Items[sel])
						state.RequestUpdate()
					}
				} else {
					l.change(state, sel, h)
					state.ClosePopups()
				}
			}
		}
	} else {
		l.grab = false
		for _, k := range state.KeyPresses() {
			switch k {
			case ui.KeyUp:
				if l.Selected > 0 {
					l.change(state, l.Selected-1, h)
				}
			case ui.KeyDown:
				if l.Selected < len(l.Items)-1 {
					l.change(state, l.Selected+1, h)
				}
			case ui.KeySpace, ui.KeyEnter:
				if l.Action != nil {
					l.Action(state, l.Items[l.Selected])
					state.RequestUpdate()
				}
			}
		}
		if text := state.TextInput(); text != "" {
			now := time.Now()
			if now.Sub(l.searchT) > time.Second {
				l.search = ""
			}
			l.searchT = now
			l.search += text
			for i, item := range l.Items {
				ls := len(l.search)
				if len(item.Text) >= ls && strings.EqualFold(l.search, item.Text[:ls]) {
					l.change(state, i, h)
					break
				}
			}
		}
	}

	hov := state.IsHovered() && !state.MouseButtonDown(ui.MouseLeft)
	for i := range l.Items {
		item := &l.Items[i]
		x, y := 2, i*h
		r := draw.XYWH(x, y, w, h)
		if i == l.Selected {
			if state.HasKeyboardFocus() {
				g.Fill(r, l.Theme.Color("selection"))
			} else {
				g.Fill(r, l.Theme.Color("selectionInactive"))
			}
		} else if hov && mouse.In(r) {
			g.Fill(r, l.Theme.Color("buttonHovered"))
		}
		item.text.DrawLeftIcon(g, draw.XYXY(x, y, w-2, y+h), item.Text, l.Theme.Font("text"), l.Theme.Color("text"), item.Icon, 3)
	}
}

func (l *List) change(state *ui.State, i, h int) {
	l.Selected = i
	if l.Changed != nil {
		l.Changed(state, l.Items[l.Selected])
		state.RequestUpdate()
	}
	state.RequestVisible(draw.XYWH(0, i*h, 1, h))
}
