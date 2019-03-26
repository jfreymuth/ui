package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type Stack struct {
	Components []ui.Component
}

func NewStack(c ...ui.Component) *Stack {
	return &Stack{Components: c}
}

func (s *Stack) SetTheme(theme *Theme) {
	for _, c := range s.Components {
		SetTheme(c, theme)
	}
}

func (s *Stack) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := 0, 0
	for _, c := range s.Components {
		cw, ch := c.PreferredSize(fonts)
		if cw > w {
			w = cw
		}
		h += ch
	}
	return w, h
}

func (s *Stack) Update(g *draw.Buffer, state *ui.State) {
	w, _ := g.Size()
	y := 0
	for _, c := range s.Components {
		_, ch := c.PreferredSize(g.FontLookup)
		state.UpdateChild(g, draw.XYWH(0, y, w, ch), c)
		y += ch
	}
}
