package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type Padding struct {
	Content                  ui.Component
	Left, Right, Top, Bottom int
}

func NewPadding(c ui.Component, p int) *Padding {
	return &Padding{Content: c, Left: p, Right: p, Top: p, Bottom: p}
}

func (p *Padding) SetTheme(theme *Theme) {
	SetTheme(p.Content, theme)
}

func (p *Padding) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := p.Content.PreferredSize(fonts)
	return p.Left + w + p.Right, p.Top + h + p.Bottom
}

func (p *Padding) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	state.UpdateChild(g, draw.XYXY(p.Left, p.Top, w-p.Right, h-p.Bottom), p.Content)
}

type FixedSize struct {
	Content       ui.Component
	Width, Height int
}

func (f *FixedSize) SetTheme(theme *Theme) {
	SetTheme(f.Content, theme)
}

func (f *FixedSize) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := f.Width, f.Height
	if w == 0 || h == 0 {
		cw, ch := f.Content.PreferredSize(fonts)
		if w == 0 {
			w = cw
		}
		if h == 0 {
			h = ch
		}
	}
	return w, h
}

func (f *FixedSize) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	state.UpdateChild(g, draw.WH(w, h), f.Content)
}

type Shadow struct {
	Content ui.Component
	theme   *Theme
}

func NewShadow(c ui.Component) *Shadow {
	return &Shadow{c, DefaultTheme}
}

func (s *Shadow) Theme() *Theme { return s.theme }
func (s *Shadow) SetTheme(theme *Theme) {
	s.theme = theme
	SetTheme(s.Content, theme)
}

func (s *Shadow) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := s.Content.PreferredSize(fonts)
	return w + 20, h + 20
}

func (s *Shadow) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	g.Shadow(draw.XYXY(12, 12, w-8, h-8), s.theme.Color("shadow"), 10)
	g.Fill(draw.XYXY(10, 10, w-10, h-10), s.theme.Color("background"))
	state.UpdateChild(g, draw.XYXY(10, 10, w-10, h-10), s.Content)
}
