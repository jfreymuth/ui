package toolkit

import (
	"image"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type Tab struct {
	Title   string
	Content ui.Component
	Close   func(*ui.State, int)
	title   text.Text
}

type TabContainer struct {
	Tabs     []Tab
	Selected int
	Changed  func(*ui.State, int)
	Theme    *Theme
	close    tabCloseButton
}

func NewTabContainer() *TabContainer {
	t := &TabContainer{Theme: DefaultTheme}
	t.close.t = t
	return t
}

func (t *TabContainer) SetTheme(theme *Theme) {
	t.Theme = theme
	for _, t := range t.Tabs {
		SetTheme(t.Content, theme)
	}
}

func (t *TabContainer) AddTab(title string, content ui.Component) {
	t.Tabs = append(t.Tabs, Tab{Title: title, Content: content})
}

func (t *TabContainer) AddClosableTab(title string, content ui.Component, close func(*ui.State, int)) {
	if close == nil {
		close = func(state *ui.State, i int) { t.CloseTab(i) }
	}
	t.Tabs = append(t.Tabs, Tab{Title: title, Content: content, Close: close})
}

func (t *TabContainer) CloseTab(i int) {
	if i < 0 || i >= len(t.Tabs) {
		return
	}
	t.Tabs = append(t.Tabs[:i], t.Tabs[i+1:]...)
}

func (t *TabContainer) PreferredSize(fonts draw.FontLookup) (int, int) {
	hw, hh := 20, 0
	cw, ch := 0, 0
	for _, tab := range t.Tabs {
		w, h := tab.title.Size(tab.Title, t.Theme.Font("title"), fonts)
		hw += w + 20
		if tab.Close != nil {
			hw += h
		}
		hh = h
		tw, th := tab.Content.PreferredSize(fonts)
		if tw > cw {
			cw = tw
		}
		if th > ch {
			ch = th
		}
	}
	if cw+20 > hw {
		hw = cw + 20
	}
	return hw, hh + 20 + ch + 20
}

func (t *TabContainer) Update(g *draw.Buffer, state *ui.State) {
	if len(t.Tabs) == 0 {
		return
	}
	if t.Selected < 0 {
		t.Selected = 0
	} else if t.Selected >= len(t.Tabs) {
		t.Selected = len(t.Tabs) - 1
	}
	w, h := g.Size()
	hh, _ := t.Tabs[0].title.Size(t.Tabs[0].Title, t.Theme.Font("title"), g.FontLookup)
	contentArea := draw.XYXY(10, hh+10, w-10, h-10)
	x := 0
	mouse := state.MousePos()
	close := -1
	for i, tab := range t.Tabs {
		hw, _ := tab.title.Size(tab.Title, t.Theme.Font("title"), g.FontLookup)
		w := hw + 20
		if tab.Close != nil {
			w += hh
		}
		rect := draw.XYWH(x+10, 10, w, hh)
		if i != t.Selected && mouse.In(rect) && state.MouseClick(ui.MouseLeft) {
			t.Selected = i
			if t.Changed != nil {
				t.Changed(state, i)
			}
			state.RequestUpdate()
			return
		}
		if i == t.Selected {
			g.Shadow(rect.Add(image.Pt(2, 2)), t.Theme.Color("shadow"), 10)
			g.Shadow(contentArea.Add(image.Pt(2, 2)), t.Theme.Color("shadow"), 10)
			g.Fill(rect, t.Theme.Color("background"))
			if tab.Close != nil {
				t.close.click = false
				state.UpdateChild(g, draw.XYWH(x+hw+30, 10, hh, hh), &t.close)
				if t.close.click {
					close = i
				}
			}
		} else if i != 0 && i != t.Selected+1 {
			g.Fill(draw.XYWH(x+10, 10, 1, hh), t.Theme.Color("veil"))
		}
		tab.title.DrawLeft(g, draw.XYWH(x+20, 0, hw, hh+20), tab.Title, t.Theme.Font("title"), t.Theme.Color("title"))
		x += w
	}
	g.Fill(contentArea, t.Theme.Color("background"))
	state.UpdateChild(g, contentArea, t.Tabs[t.Selected].Content)
	if close >= 0 {
		t.Tabs[close].Close(state, close)
		state.RequestUpdate()
	}
}

type tabCloseButton struct {
	t     *TabContainer
	click bool
}

func (*tabCloseButton) PreferredSize(draw.FontLookup) (int, int) { return 0, 0 }
func (b *tabCloseButton) Update(g *draw.Buffer, state *ui.State) {
	color := b.t.Theme.Color("buttonText")
	if state.IsHovered() {
		color = b.t.Theme.Color("buttonFocused")
	}
	g.Icon(draw.WH(g.Size()), "close", color)
	if state.MouseClick(ui.MouseLeft) {
		b.click = true
	}
}
