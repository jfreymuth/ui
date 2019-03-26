package gldraw

import (
	"image"

	"github.com/jfreymuth/ui/draw"
)

type IconLookup interface {
	// IconSize should return the next smaller supported icon size.
	IconSize(int) int
	// DrawIcon draws an icon. The size of the image can be assumed to be a value returned by IconSize.
	DrawIcon(*image.Alpha, string)
}

type iconContext struct {
	iconLookup IconLookup
	iconSheets map[int]*iconSheet
}

type iconSheet struct {
	*iconContext
	size  int
	img   image.Alpha
	tex   *Texture
	icons map[string]int
}

func (c *iconSheet) load(icon string) int {
	if n, ok := c.icons[icon]; ok {
		return n
	}
	c.iconLookup.DrawIcon(&c.img, icon)
	n := len(c.icons)
	x, y := n%16, n/16
	c.tex.UpdateAlpha(&c.img, image.Pt(x*c.size, y*c.size))
	c.icons[icon] = n
	return n
}

func (c *Context) getIconSheet(size int) *iconSheet {
	size = c.iconLookup.IconSize(size)
	if s, ok := c.iconSheets[size]; ok {
		return s
	}
	s := c.newIconSheet(size)
	if c.iconSheets == nil {
		c.iconSheets = make(map[int]*iconSheet)
	}
	c.iconSheets[size] = s
	return s
}

func (c *Context) newIconSheet(size int) *iconSheet {
	c.prepare(0)
	return &iconSheet{
		iconContext: &c.iconContext,
		size:        size,
		img:         *image.NewAlpha(image.Rect(0, 0, size, size)),
		tex:         NewTextureAlphaEmpty(size*16, size*16),
		icons:       make(map[string]int),
	}
}

func (c *Context) drawIcon(r, clip image.Rectangle, icon string, color draw.Color) {
	if c.iconLookup == nil || icon == "" {
		return
	}
	s := r.Dx()
	if r.Dy() < s {
		s = r.Dy()
	}
	sheet := c.getIconSheet(s)
	s = sheet.size
	dx := (r.Dx() - s) / 2
	dy := (r.Dy() - s) / 2
	r = draw.XYWH(r.Min.X+dx, r.Min.Y+dy, s, s)
	n := sheet.load(icon)
	x, y := float32(n%16), float32(n/16)
	c.prepare(sheet.tex.tex)
	c.rect(r, clip, [4]float32{x / 16, y / 16, (x + 1) / 16, (y + 1) / 16}, color)
}
