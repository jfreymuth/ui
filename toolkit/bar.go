package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type Bar struct {
	Theme      *Theme
	Components []ui.Component
	Fill       int
}

func NewBar(fill int, c ...ui.Component) *Bar {
	return &Bar{Theme: DefaultTheme, Components: c, Fill: fill}
}

func (b *Bar) SetTheme(theme *Theme) {
	b.Theme = theme
	for _, c := range b.Components {
		SetTheme(c, theme)
	}
}

func (b *Bar) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := 0, 0
	for _, c := range b.Components {
		cw, ch := c.PreferredSize(fonts)
		w += cw
		if ch > h {
			h = ch
		}
	}
	return w, h
}

func (b *Bar) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	g.Fill(draw.WH(w, h), b.Theme.Color("altBackground"))
	x := 0
	if b.Fill >= len(b.Components) {
		for _, c := range b.Components {
			cw, _ := c.PreferredSize(g.FontLookup)
			state.UpdateChild(g, draw.XYWH(x, 0, cw, h), c)
			x += cw
		}
	} else if b.Fill < 0 {
		ws := make([]int, len(b.Components))
		for i := range ws {
			ws[i], _ = b.Components[b.Fill+1+i].PreferredSize(g.FontLookup)
			w -= ws[i]
		}
		x += w
		for i, c := range b.Components {
			cw := ws[i]
			state.UpdateChild(g, draw.XYWH(x, 0, cw, h), c)
			x += cw
		}
	} else {
		for _, c := range b.Components[:b.Fill] {
			cw, _ := c.PreferredSize(g.FontLookup)
			state.UpdateChild(g, draw.XYWH(x, 0, cw, h), c)
			x += cw
		}
		ws := make([]int, len(b.Components)-b.Fill-1)
		w -= x
		for i := range ws {
			ws[i], _ = b.Components[b.Fill+1+i].PreferredSize(g.FontLookup)
			w -= ws[i]
		}
		state.UpdateChild(g, draw.XYWH(x, 0, w, h), b.Components[b.Fill])
		x += w
		for i := b.Fill + 1; i < len(b.Components); i++ {
			c := b.Components[i]
			cw := ws[i-b.Fill-1]
			state.UpdateChild(g, draw.XYWH(x, 0, cw, h), c)
			x += cw
		}
	}
}
