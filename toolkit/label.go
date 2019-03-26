package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type Label struct {
	Theme *Theme
	Text  string
	text  text.Text
}

func NewLabel(text string) *Label {
	return &Label{Theme: DefaultTheme, Text: text}
}

func (l *Label) SetTheme(theme *Theme) { l.Theme = theme }

func (l *Label) PreferredSize(fonts draw.FontLookup) (int, int) {
	return l.text.Size(l.Text, l.Theme.Font("text"), fonts)
}

func (l *Label) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	l.text.DrawLeft(g, draw.WH(w, h), l.Text, l.Theme.Font("text"), l.Theme.Color("text"))
}
