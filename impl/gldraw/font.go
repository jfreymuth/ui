package gldraw

import (
	"image"
	idraw "image/draw"

	"github.com/jfreymuth/ui/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type FontLookup interface {
	draw.FontLookup
	LoadFont(draw.Font) font.Face
}

type fontContext struct {
	fontLookup   FontLookup
	fontMap      map[draw.Font]*fontFace
	fontFaces    []*fontFace
	fontTextures []fontTexture
	nextGC       int
}

type fontFace struct {
	draw.Font
	f       font.Face
	height  int
	img     image.Alpha
	glyphs  map[glyphSpec]glyph
	texture *fontTexture
	line    *fontLine
	metrics font.Metrics
}

type fontTexture struct {
	lines *fontLine
	free  *fontLine
	img   *Texture
	space float32
}

type fontLine struct {
	draw.Font
	y, h int
	used uint
	x    int
	next *fontLine
}

type glyphSpec struct {
	r  rune
	sp byte
}

type glyph struct {
	v    [4]float32
	t    [4]float32
	adv  fixed.Int26_6
	line *fontLine
}

const fontTextureWidth = 1024
const fontTextureHeight = 1024
const fontGCThreshold = 10

func (c *fontContext) initFonts() {
	c.fontMap = make(map[draw.Font]*fontFace)
}

func (c *fontContext) getFontFace(s draw.Font) *fontFace {
	if f, ok := c.fontMap[s]; ok {
		return f
	}
	n := c.fontLookup.GetClosest(s)
	if f, ok := c.fontMap[n]; ok {
		c.fontMap[s] = f
		return f
	}
	f := &fontFace{Font: n, f: c.fontLookup.LoadFont(n), glyphs: make(map[glyphSpec]glyph)}
	f.metrics = f.f.Metrics()
	f.height = f.metrics.Ascent.Ceil() + f.metrics.Descent.Ceil() + 1
	f.img = *image.NewAlpha(image.Rect(0, 0, f.height*2, f.height))
	c.fontMap[n] = f
	c.fontFaces = append(c.fontFaces, f)
	return f
}

func (c *Context) getFontTexture(h int) *fontTexture {
	if len(c.fontTextures) > 0 {
		for i := range c.fontTextures {
			t := &c.fontTextures[i]
			if t.space < .8 {
				return t
			}
		}
		c.nextGC++
		if c.nextGC >= len(c.fontTextures) {
			c.nextGC = 0
		}
		i := c.nextGC
		for {
			t := &c.fontTextures[i]
			t.gc(c)
			if t.space < .8 {
				return t
			}
			i++
			if i == len(c.fontTextures) {
				i = 0
			}
			if i == c.nextGC {
				break
			}
		}
	}
	return c.newFontTexture(fontTextureWidth, fontTextureHeight)
}

func (c *Context) newFontTexture(w, h int) *fontTexture {
	t := NewTextureAlphaEmpty(w, h)
	t.c = c
	c.fontTextures = append(c.fontTextures, fontTexture{img: t})
	return &c.fontTextures[len(c.fontTextures)-1]
}

func (t *fontTexture) alloc(s draw.Font, h int, c *Context) *fontLine {
	if t.lines == nil {
		next := &fontLine{y: h, h: t.img.height - h}
		t.lines = &fontLine{Font: s, h: h, next: next}
		t.space += float32(h) / float32(t.img.height)
		return t.lines
	} else {
		l := t.lines
		for l.Name != "" || l.h < h {
			l = l.next
			if l == nil {
				//t.gc(c, 1)
				t.freeOld(c)
				return t.alloc(s, h, c)
			}
		}
		if l.h == h {
			l.Font = s
			t.space += float32(h) / float32(t.img.height)
			return l
		} else {
			new := &fontLine{y: l.y + h, h: l.h - h, next: l.next}
			l.Font = s
			l.h = h
			l.next = new
			t.space += float32(h) / float32(t.img.height)
			return l
		}
	}
}

func (t *fontTexture) gc(c *Context) {
	var p, l *fontLine = nil, t.lines
	for l != nil {
		if l.Name != "" && l.used < c.time-fontGCThreshold {
			l = t.remove(c, p, l)
		}
		p, l = l, l.next
	}
}

func (t *fontTexture) remove(c *Context, p *fontLine, l *fontLine) *fontLine {
	for _, face := range c.fontFaces {
		if face.Font == l.Font {
			if face.line == l {
				face.line = nil
			}
			i := 0
			for r, g := range face.glyphs {
				if g.line == l {
					delete(face.glyphs, r)
					i++
				}
			}
		}
	}
	t.space -= float32(l.h) / float32(t.img.height)
	if p != nil && p.Name == "" {
		if l.next != nil && l.next.Name == "" {
			p.h += l.h + l.next.h
			p.next = l.next.next
		} else {
			p.h += l.h
			p.next = l.next
		}
		return p
	} else {
		if l.next != nil && l.next.Name == "" {
			l.next.y -= l.h
			l.next.h += l.h
			if p != nil {
				p.next = l.next
			} else {
				t.lines = l.next
			}
			return l.next
		} else {
			l.Name = ""
			return l
		}
	}
}

func (t *fontTexture) freeOld(c *Context) {
	var best, bestP *fontLine
	var bestT uint = ^uint(0)
	p, l := t.lines, t.lines
	for l != nil {
		if l.Name != "" && l.used < bestT {
			best, bestP = l, p
			bestT = l.used
		}
		p, l = l, l.next
	}
	t.remove(c, bestP, best)
}

func (f *fontFace) load(r glyphSpec, c *Context) glyph {
	if g, ok := f.glyphs[r]; ok {
		g.line.used = c.time
		return g
	}
	dr, mask, maskp, adv, ok := f.f.Glyph(fixed.Point26_6{X: fixed.Int26_6(r.sp << (6 - f.subpixels())), Y: 0}, r.r)
	if !ok {
		err := glyph{}
		f.glyphs[r] = err
		return err
	}
	if f.line == nil || f.line.x+dr.Dx() > f.texture.img.width {
		if len(f.glyphs) == 0 {
			f.texture = c.getFontTexture(f.height)
		}
		f.line = f.texture.alloc(f.Font, f.height, c)
		f.line.x = 0
	}
	g := glyph{line: f.line, adv: adv}
	g.v = [4]float32{float32(dr.Min.X), float32(dr.Min.Y), float32(dr.Max.X), float32(dr.Max.Y)}
	rect := image.Rect(f.line.x, f.line.y, f.line.x+dr.Dx(), f.line.y+dr.Dy())
	w, h := float32(f.texture.img.width), float32(f.texture.img.height)
	x, y := float32(rect.Min.X)/w, float32(rect.Min.Y)/h
	g.t = [4]float32{x, y, x + (g.v[2]-g.v[0])/w, y + (g.v[3]-g.v[1])/h}
	f.line.x += dr.Dx() + 1

	if img, ok := mask.(*image.Alpha); ok {
		img.Rect = img.Rect.Intersect(f.img.Rect)
		f.texture.img.UpdateAlpha(img, rect.Min)
	} else {
		idraw.Draw(&f.img, f.img.Rect, mask, maskp, idraw.Src)
		f.texture.img.UpdateAlpha(&f.img, rect.Min)
	}
	f.glyphs[r] = g
	f.line.used = c.time
	return g
}

func (f *fontFace) advance(r rune) fixed.Int26_6 {
	if g, ok := f.glyphs[glyphSpec{r, 0}]; ok {
		return g.adv
	}
	adv, _ := f.f.GlyphAdvance(r)
	return adv
}

func (f *fontFace) subpixels() byte {
	if f.Size < 20 {
		return 2
	} else if f.Size < 30 {
		return 1
	}
	return 0
}

func (f *fontFace) write(s string, xf, yf float32, c *Context, scale float32, clip image.Rectangle, color draw.Color) (float32, string) {
	x, y := fixed.Int26_6(xf/scale*64), fixed.Int26_6(yf/scale*64)
	var last rune
	for _, r := range s {
		kern := f.f.Kern(last, r)
		x += kern
		var adv fixed.Int26_6
		{
			x, y := float32(x)/64, float32(y)/64
			var subp byte = 0
			xint := float32(int(x+100)) - 100
			xfrac := x - xint
			sp := float32(int(1) << f.subpixels())
			subp = byte(xfrac * sp)
			x = xint
			y = float32(int(y + .5))
			g := f.load(glyphSpec{r, subp}, c)
			adv = g.adv
			c.prepare(f.texture.img.tex)
			c.rect(image.Rect(int(x+g.v[0]), int(y+g.v[1]), int(x+g.v[2]), int(y+g.v[3])), clip, g.t, color)
		}
		x += adv
		last = r
	}
	return float32(x)/64*scale - xf, ""
}

func (c *Context) DebugShowFontTextures(s float32) {
	for i := range c.fontTextures {
		c.prepare(c.fontTextures[i].img.tex)
		c.rect(draw.XYWH(int(float32(i*fontTextureWidth)*s), 0, int(fontTextureWidth*s), int(fontTextureHeight*s)), draw.WH(2000, 2000), [4]float32{0, 0, 1, 1}, draw.Black)
	}
	c.buffer.flush()
}

func (f *fontFace) Ascent() int     { return f.metrics.Ascent.Ceil() }
func (f *fontFace) Descent() int    { return f.metrics.Descent.Ceil() }
func (f *fontFace) LineHeight() int { return f.metrics.Height.Ceil() }

func (f *fontFace) Advance(s string) float32 {
	var x fixed.Int26_6
	var last rune
	for _, r := range s {
		x += f.f.Kern(last, r)
		x += f.advance(r)
		last = r
	}
	return float32(x) / 64
}

func (f *fontFace) Index(s string, t float32) int {
	var x fixed.Int26_6
	var last rune
	for i, r := range s {
		x += f.f.Kern(last, r)
		adv := f.advance(r)
		if t < float32(x+adv/2)/64 {
			return i
		}
		x += adv
		last = r
	}
	return len(s)
}
