package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type Frame struct {
	Theme   *Theme
	Content ui.Component
	Color   draw.Color
	Title   string
	Icon    string
	title   text.Text
}

func NewFrame(icon, title string, content ui.Component, color draw.Color) *Frame {
	return &Frame{Theme: DefaultTheme, Content: content, Color: color, Title: title, Icon: icon}
}

func (f *Frame) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := f.Content.PreferredSize(fonts)
	tw, th := f.title.SizeIcon(f.Title, f.Theme.Font("title"), f.Icon, 5, fonts)
	if w < tw {
		w = tw
	}
	return w + 20, h + th + 30
}

func (f *Frame) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	_, th := f.title.Size(f.Title, f.Theme.Font("title"), g.FontLookup)
	g.Shadow(draw.XYXY(12, 12, w-8, h-8), draw.RGBA(0, 0, 0, .5), 10)
	g.Fill(draw.XYXY(10, 10, w-10, th+20), f.Color)
	f.title.DrawCenteredIcon(g, draw.XYXY(15, 10, w-15, th+20), f.Title, f.Theme.Font("title"), f.Theme.Color("title"), f.Icon, 5)
	g.Fill(draw.XYXY(10, th+20, w-10, h-10), DefaultTheme.Color("background"))
	state.UpdateChild(g, draw.XYXY(10, th+20, w-10, h-10), f.Content)
}
