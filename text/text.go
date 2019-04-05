package text

import (
	"image"

	"github.com/jfreymuth/ui/draw"
)

type Text struct {
	text    string
	font    draw.Font
	w, h, b int
}

func (t *Text) Size(text string, font draw.Font, fonts draw.FontLookup) (int, int) {
	t.measure(text, font, fonts)
	return t.w, t.h
}

func (t *Text) measure(text string, font draw.Font, fonts draw.FontLookup) {
	if text != t.text || font != t.font {
		t.text, t.font = text, font
		m := fonts.Metrics(font)
		t.w, t.h = int(m.Advance(text)), m.LineHeight()
		t.b = (m.Ascent() - m.Descent()) / 2
	}
}

func (t *Text) DrawLeft(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color) {
	t.measure(text, font, g.FontLookup)
	g.Add(draw.Text{Position: image.Pt(r.Min.X, (r.Min.Y+r.Max.Y)/2+t.b), Text: text, Font: font, Color: color})
}

func (t *Text) DrawCentered(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color) {
	w, _ := t.Size(text, font, g.FontLookup)
	g.Add(draw.Text{Position: image.Pt(r.Min.X+(r.Dx()-w)/2, (r.Min.Y+r.Max.Y)/2+t.b), Text: text, Font: font, Color: color})
}

func (t *Text) DrawRight(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color) {
	w, _ := t.Size(text, font, g.FontLookup)
	g.Add(draw.Text{Position: image.Pt(r.Max.X-w, (r.Min.Y+r.Max.Y)/2+t.b), Text: text, Font: font, Color: color})
}

func (t *Text) SizeIcon(text string, font draw.Font, icon string, gap int, fonts draw.FontLookup) (int, int) {
	t.measure(text, font, fonts)
	if icon != "" {
		if text == "" {
			return t.w + t.h, t.h
		}
		return t.w + t.h + gap, t.h
	}
	return t.w, t.h
}

func (t *Text) DrawLeftIcon(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color, icon string, gap int) {
	_, h := t.SizeIcon(text, font, icon, gap, g.FontLookup)
	if icon != "" {
		g.Add(draw.Icon{Rect: draw.XYWH(r.Min.X, r.Min.Y, h, r.Dy()), Icon: icon, Color: color})
		r.Min.X += gap + h
	}
	g.Add(draw.Text{Position: image.Pt(r.Min.X, (r.Min.Y+r.Max.Y)/2+t.b), Text: text, Font: font, Color: color})
}

func (t *Text) DrawCenteredIcon(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color, icon string, gap int) {
	w, h := t.Size(text, font, g.FontLookup)
	if icon != "" {
		if text == "" {
			gap = 0
		}
		g.Add(draw.Icon{Rect: draw.XYWH(r.Min.X+(r.Dx()-w-gap-h)/2, r.Min.Y, h, r.Dy()), Icon: icon, Color: color})
		r.Min.X += gap + h
	}
	g.Add(draw.Text{Position: image.Pt(r.Min.X+(r.Dx()-w)/2, (r.Min.Y+r.Max.Y)/2+t.b), Text: text, Font: font, Color: color})
}

func (t *Text) DrawRightIcon(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color, icon string, gap int) {
	w, h := t.Size(text, font, g.FontLookup)
	if icon != "" {
		g.Add(draw.Icon{Rect: draw.XYWH(r.Min.X+r.Dx()-w-gap-h, r.Min.Y, h, r.Dy()), Icon: icon, Color: color})
		r.Min.X += gap + h
	}
	g.Add(draw.Text{Position: image.Pt(r.Max.X-w, (r.Min.Y+r.Max.Y)/2+t.b), Text: text, Font: font, Color: color})
}
