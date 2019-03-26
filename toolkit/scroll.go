package toolkit

import (
	"image"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type scrollBar struct {
	*ScrollView
	Max, Size, Value float32
	Vertical         bool
	grab             int
}

func (s *scrollBar) SetValue(v float32) {
	if v < 0 {
		s.Value = 0
	} else if v > s.Max {
		s.Value = s.Max
	} else {
		s.Value = v
	}
}

func (s *scrollBar) PreferredSize(fonts draw.FontLookup) (int, int) {
	if s.Vertical {
		return 15, 0
	} else {
		return 0, 15
	}
}

func (s *scrollBar) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	l := 0
	if s.Vertical {
		l = h
	} else {
		l = w
	}
	b := int(s.Size / (s.Max + s.Size) * float32(l))
	if b < 30 {
		b = 30
	} else if b > l {
		b = l
	}
	o := int(s.Value / s.Max * float32(l-b))
	if state.MouseButtonDown(ui.MouseLeft) {
		m := 0
		if s.Vertical {
			m = state.MousePos().Y
		} else {
			m = state.MousePos().X
		}
		if s.grab < 0 {
			if m >= o && m < o+b {
				s.grab = m - o
			} else {
				s.grab = b / 2
			}
		}
		c := float32(m-s.grab) / float32(l-b)
		if c < 0 {
			s.Value = 0
		} else if c > 1 {
			s.Value = s.Max
		} else {
			s.Value = c * s.Max
		}
	} else {
		s.grab = -1
	}
	o = int(s.Value / s.Max * float32(l-b))
	g.Outline(draw.WH(w, h), s.Theme.Color("border"))
	if s.Vertical {
		g.Fill(draw.XYWH(1, o, w-2, b), s.Theme.Color("scrollBar"))
	} else {
		g.Fill(draw.XYWH(o, 1, b, h-2), s.Theme.Color("scrollBar"))
	}
}

type ScrollView struct {
	content ui.Component
	Theme   *Theme
	sh, sv  scrollBar
	viewport
}

func NewScrollView(content ui.Component) *ScrollView {
	s := &ScrollView{content: content, Theme: DefaultTheme}
	s.sh = scrollBar{s, 1, 0, 0, false, -1}
	s.sv = scrollBar{s, 1, 0, 0, true, -1}
	s.viewport.ScrollView = s
	return s
}

func (s *ScrollView) SetTheme(theme *Theme) {
	s.Theme = theme
	SetTheme(s.content, theme)
}

func (s *ScrollView) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := s.content.PreferredSize(fonts)
	return w + 15, h
}
func (s *ScrollView) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	s.w, s.h = s.content.PreferredSize(g.FontLookup)
	showH, showV := false, false
	if s.w > w {
		showH = true
		h -= 15
	}
	if s.h > h {
		showV = true
		w -= 15
	}
	if !showH && s.w > w {
		showH = true
		h -= 15
	}
	vv, hv := s.sv.Value, s.sh.Value
	if showH {
		s.sh.Max = float32(s.w - w)
		s.sh.Size = float32(w)
		s.sh.SetValue(s.sh.Value)
		state.UpdateChild(g, draw.XYWH(0, h, w, 15), &s.sh)
	}
	if showV {
		s.sv.Max = float32(s.h - h)
		s.sv.Size = float32(h)
		s.sv.SetValue(s.sv.Value)
		state.UpdateChild(g, draw.XYWH(w, 0, 15, h), &s.sv)
	}
	state.UpdateChild(g, draw.WH(w, h), &s.viewport)
	if r, ok := state.GetVisibilityRequest(); ok {
		if r.Dx() <= w {
			if r.Min.X < 0 {
				s.sh.SetValue(s.sh.Value + float32(r.Min.X))
			} else if r.Max.X > w {
				s.sh.SetValue(s.sh.Value + float32(r.Max.X-w))
			}
		}
		if r.Dy() <= h {
			if r.Min.Y < 0 {
				if r.Max.Y < h {
					s.sv.SetValue(s.sv.Value + float32(r.Min.Y))
				}
			} else if r.Max.Y > h {
				s.sv.SetValue(s.sv.Value + float32(r.Max.Y-h))
			}
		}
	}
	if scroll := state.Scroll(); scroll != (image.Point{}) {
		if showH && (scroll.X > 0 && s.sh.Value > 0 || scroll.X < 0 && s.sh.Value < s.sh.Max) ||
			showV && (scroll.Y > 0 && s.sv.Value > 0 || scroll.Y < 0 && s.sv.Value < s.sv.Max) {
			state.ConsumeScroll()
		}
		s.sh.SetValue(s.sh.Value - float32(scroll.X*45))
		s.sv.SetValue(s.sv.Value - float32(scroll.Y*45))
	}
	if s.sv.Value != vv || s.sh.Value != hv {
		state.RequestUpdate()
	}
}

type viewport struct {
	w, h int
	*ScrollView
}

func (viewport) PreferredSize(fonts draw.FontLookup) (int, int) { return 0, 0 }
func (v *viewport) Update(g *draw.Buffer, s *ui.State) {
	x := -int(v.sh.Value)
	y := -int(v.sv.Value)
	w, h := g.Size()
	if w >= v.w {
		x, v.w = 0, w
	}
	if h >= v.h {
		y, v.h = 0, h
	}
	s.UpdateChild(g, draw.XYWH(x, y, v.w, v.h), v.content)
}
