package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

func animate(state *ui.State, v *float32, speed float32, on bool) {
	if on {
		*v += state.AnimationSpeed() * speed
		if *v >= 1 {
			*v = 1
		} else {
			state.RequestAnimation()
		}
	} else {
		*v -= state.AnimationSpeed() * speed
		if *v <= 0 {
			*v = 0
		} else {
			state.RequestAnimation()
		}
	}
}

type Separator struct {
	Theme         *Theme
	Width, Height int
}

func NewSeparator(w, h int) *Separator {
	return &Separator{Theme: DefaultTheme, Width: w, Height: h}
}

func (s *Separator) PreferredSize(fonts draw.FontLookup) (int, int) {
	return s.Width, s.Height
}

func (s *Separator) Update(g *draw.Buffer, state *ui.State) {
	g.Fill(draw.WH(g.Size()), s.Theme.Color("veil"))
}
