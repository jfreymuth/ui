package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type Divider struct {
	First, Second ui.Component
	Vertical      bool
	Theme         *Theme
	pos           int
}

type dividerBar struct {
	*Divider
}

func NewVerticalDivider(top, bottom ui.Component) *Divider {
	return &Divider{top, bottom, true, DefaultTheme, -1}
}

func NewHorizontalDivider(left, right ui.Component) *Divider {
	return &Divider{left, right, false, DefaultTheme, -1}
}

func (d *Divider) SetTheme(theme *Theme) {
	d.Theme = theme
	SetTheme(d.First, theme)
	SetTheme(d.Second, theme)
}

func (d *Divider) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := d.First.PreferredSize(fonts)
	w2, h2 := d.Second.PreferredSize(fonts)
	if d.Vertical {
		if w2 > w {
			w = w2
		}
		h += h2 + 4
	} else {
		w += w2 + 4
		if h2 > h {
			h = h2
		}
	}
	return w, h
}

func (d *Divider) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	if d.Vertical {
		if d.pos == -1 {
			_, h1 := d.First.PreferredSize(g.FontLookup)
			_, h2 := d.Second.PreferredSize(g.FontLookup)
			d.pos = int(float32(h1+2) / float32(h1+h2+4) * float32(h))
		}
		if d.pos < 2 {
			d.pos = 2
		} else if d.pos > h-2 {
			d.pos = h - 2
		}
		state.UpdateChild(g, draw.WH(w, d.pos-2), d.First)
		state.UpdateChild(g, draw.XYXY(0, d.pos+2, w, h), d.Second)
		state.UpdateChild(g, draw.XYXY(0, d.pos-2, w, d.pos+2), dividerBar{d})
	} else {
		if d.pos == -1 {
			w1, _ := d.First.PreferredSize(g.FontLookup)
			w2, _ := d.Second.PreferredSize(g.FontLookup)
			d.pos = int(float32(w1+2) / float32(w1+w2+4) * float32(w))
		}
		if d.pos < 2 {
			d.pos = 2
		} else if d.pos > w-2 {
			d.pos = w - 2
		}
		state.UpdateChild(g, draw.WH(d.pos-2, h), d.First)
		state.UpdateChild(g, draw.XYXY(d.pos+2, 0, w, h), d.Second)
		state.UpdateChild(g, draw.XYXY(d.pos-2, 0, d.pos+2, h), dividerBar{d})
	}
}

func (d dividerBar) PreferredSize(fonts draw.FontLookup) (int, int) { return 0, 0 }

func (d dividerBar) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	g.Fill(draw.WH(w, h), d.Theme.Color("altBackground"))
	if d.Vertical {
		state.SetCursor(ui.CursorResizeVertical)
		if state.MouseButtonDown(ui.MouseLeft) {
			y := state.MousePos().Y
			d.pos += y - 2
			state.RequestUpdate()
		}
	} else {
		state.SetCursor(ui.CursorResizeHorizontal)
		if state.MouseButtonDown(ui.MouseLeft) {
			x := state.MousePos().X
			d.pos += x - 2
			state.RequestUpdate()
		}
	}
	if d.pos < 0 {
		d.pos = 0
	}
}
