package text

import (
	"image"

	"github.com/jfreymuth/ui/draw"
)

type Text struct {
	text string
	font draw.Font
	w, h int
}

func (t *Text) Size(text string, font draw.Font, fonts draw.FontLookup) (int, int) {
	if text != t.text || font != t.font {
		t.text, t.font = text, font
		m := fonts.Metrics(font)
		t.w, t.h = int(m.Advance(text)), m.LineHeight()
	}
	return t.w, t.h
}

func (t *Text) DrawLeft(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color) {
	g.Text(r, text, color, font)
}

func (t *Text) DrawCentered(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color) {
	w, _ := t.Size(text, font, g.FontLookup)
	g.Text(draw.XYWH(r.Min.X+(r.Dx()-w)/2, r.Min.Y, w, r.Dy()), text, color, font)
}

func (t *Text) DrawRight(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color) {
	w, _ := t.Size(text, font, g.FontLookup)
	r.Min.X = r.Max.X - w
	g.Text(r, text, color, font)
}

func (t *Text) SizeIcon(text string, font draw.Font, icon string, gap int, fonts draw.FontLookup) (int, int) {
	if text != t.text || font != t.font {
		t.text, t.font = text, font
		m := fonts.Metrics(font)
		t.w, t.h = int(m.Advance(text)), m.LineHeight()
	}
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
		g.Icon(draw.XYWH(r.Min.X, r.Min.Y, h, r.Dy()), icon, color)
		r.Min.X += gap + h
	}
	g.Text(r, text, color, font)
}

func (t *Text) DrawCenteredIcon(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color, icon string, gap int) {
	w, h := t.Size(text, font, g.FontLookup)
	if icon != "" {
		if text == "" {
			gap = 0
		}
		g.Icon(draw.XYWH(r.Min.X+(r.Dx()-w-gap-h)/2, r.Min.Y, h, r.Dy()), icon, color)
		r.Min.X += gap + h
	}
	g.Text(draw.XYWH(r.Min.X+(r.Dx()-w)/2, r.Min.Y, w, r.Dy()), text, color, font)
}

func (t *Text) DrawRightIcon(g *draw.Buffer, r image.Rectangle, text string, font draw.Font, color draw.Color, icon string, gap int) {
	w, h := t.Size(text, font, g.FontLookup)
	if icon != "" {
		g.Icon(draw.XYWH(r.Min.X+r.Dx()-w-gap-h, r.Min.Y, h, r.Dy()), icon, color)
		r.Min.X += gap + h
	}
	g.Text(draw.XYWH(r.Min.X+r.Dx()-w, r.Min.Y, w, r.Dy()), text, color, font)
}
