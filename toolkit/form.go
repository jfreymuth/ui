package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type Form struct {
	Theme  *Theme
	Fields []FormField
	w, h   int
}

type FormField struct {
	Content ui.Component
	Label   string
	label   text.Text
}

func NewForm() *Form {
	return &Form{Theme: DefaultTheme}
}

func (f *Form) SetTheme(theme *Theme) {
	f.Theme = theme
	for _, c := range f.Fields {
		SetTheme(c.Content, theme)
	}
}

func (f *Form) AddField(label string, c ui.Component) {
	f.Fields = append(f.Fields, FormField{Label: label, Content: c})
}

func (f *Form) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := 0, 5
	f.w, f.h = 0, 0
	for _, ff := range f.Fields {
		cw, ch := ff.Content.PreferredSize(fonts)
		lw, lh := ff.label.Size(ff.Label, f.Theme.Font("text"), fonts)
		if ch < lh {
			ch = lh
		}
		h += ch + 5
		if cw > w {
			w = cw
		}
		if lw > f.w {
			f.w = lw
		}
		if lh > f.h {
			f.h = lh
		}
	}
	return w + f.w + 15, h
}

func (f *Form) Update(g *draw.Buffer, state *ui.State) {
	f.measure(g.FontLookup)
	w, _ := g.Size()
	y := 5
	for _, ff := range f.Fields {
		_, ch := ff.Content.PreferredSize(g.FontLookup)
		if ch < f.h {
			ch = f.h
		}
		ff.label.DrawRight(g, draw.XYWH(5, y, f.w, ch), ff.Label, f.Theme.Font("text"), f.Theme.Color("text"))
		state.UpdateChild(g, draw.XYXY(f.w+10, y, w-5, y+ch), ff.Content)
		y += ch + 5
	}
}

func (f *Form) measure(fonts draw.FontLookup) {
	f.w, f.h = 0, 0
	for _, ff := range f.Fields {
		w, h := ff.label.Size(ff.Label, f.Theme.Font("text"), fonts)
		if w > f.w {
			f.w = w
		}
		if h > f.h {
			f.h = h
		}
	}
}
