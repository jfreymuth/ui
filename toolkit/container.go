package toolkit

import (
	"image"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type Container struct {
	Center, Top, Bottom, Left, Right ui.Component
}

func (c *Container) SetTheme(theme *Theme) {
	SetTheme(c.Top, theme)
	SetTheme(c.Left, theme)
	SetTheme(c.Center, theme)
	SetTheme(c.Right, theme)
	SetTheme(c.Bottom, theme)
}

func (c *Container) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := 0, 0
	if c.Center != nil {
		w, h = c.Center.PreferredSize(fonts)
	}
	if c.Left != nil {
		cw, ch := c.Left.PreferredSize(fonts)
		w += cw
		if ch > h {
			h = ch
		}
	}
	if c.Right != nil {
		cw, ch := c.Right.PreferredSize(fonts)
		w += cw
		if ch > h {
			h = ch
		}
	}
	if c.Top != nil {
		cw, ch := c.Top.PreferredSize(fonts)
		if cw > w {
			w = cw
		}
		h += ch
	}
	if c.Bottom != nil {
		cw, ch := c.Bottom.PreferredSize(fonts)
		if cw > w {
			w = cw
		}
		h += ch
	}
	return w, h
}

func (c *Container) Update(g *draw.Buffer, state *ui.State) {
	x, y := 0, 0
	w, h := g.Size()
	if c.Top != nil {
		_, ch := c.Top.PreferredSize(g.FontLookup)
		if ch > h {
			ch = h
		}
		state.UpdateChild(g, draw.XYWH(x, y, w, ch), c.Top)
		y += ch
		h -= ch
	}
	var bottom, right image.Rectangle
	if c.Bottom != nil {
		_, ch := c.Bottom.PreferredSize(g.FontLookup)
		if ch > h {
			ch = h
		}
		bottom = draw.XYWH(x, y+h-ch, w, ch)
		h -= ch
	}
	if c.Left != nil {
		cw, _ := c.Left.PreferredSize(g.FontLookup)
		if cw > w {
			cw = w
		}
		state.UpdateChild(g, draw.XYWH(x, y, cw, h), c.Left)
		x += cw
		w -= cw
	}
	if c.Right != nil {
		cw, _ := c.Right.PreferredSize(g.FontLookup)
		if cw > w {
			cw = w
		}
		right = draw.XYWH(x+w-cw, y, cw, h)
		w -= cw
	}
	if c.Center != nil {
		state.UpdateChild(g, draw.XYWH(x, y, w, h), c.Center)
	}
	if c.Right != nil {
		state.UpdateChild(g, right, c.Right)
	}
	if c.Bottom != nil {
		state.UpdateChild(g, bottom, c.Bottom)
	}
}
