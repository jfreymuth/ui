package gldraw

import (
	"image"

	"github.com/jfreymuth/ui/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
)

type Context struct {
	buffer  buffer
	sbuffer sbuffer
	images  map[*image.RGBA]*entry
	empty   *Texture

	currentTexture uint32

	fontContext
	iconContext
	time uint
}

func (c *Context) Init(f FontLookup) {
	if f == nil {
		panic("gldraw: FontLookup must not be nil")
	}
	c.buffer.init(1024)
	c.sbuffer.init(512)
	c.images = make(map[*image.RGBA]*entry)
	c.empty = NewTextureAlpha(&image.Alpha{Pix: []uint8{255}, Stride: 1, Rect: image.Rect(0, 0, 1, 1)})
	c.empty.c = c
	c.fontLookup = f
	c.initFonts()
}

func (c *Context) FontLookup() draw.FontLookup {
	return c.fontLookup
}

func (c *Context) SetIconLookup(l IconLookup) {
	if l != nil {
		c.iconLookup = l
	}
}

func (c *Context) Draw(w, h int, cmd []draw.CommandList) {
	gl.Viewport(0, 0, int32(w), int32(h))
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	c.buffer.setScreenSize(w, h)
	c.sbuffer.setScreenSize(w, h)
	c.currentTexture = 0
	for _, l := range cmd {
		for _, cmd := range l.Commands {
			switch cmd := cmd.(type) {
			case draw.Fill:
				c.prepare(c.empty.tex)
				c.buffer.rect(cmd.Rect.Add(l.Offset).Intersect(l.Clip), m.Vec2{}, m.Vec2{}, cmd.Color)
			case draw.Outline:
				cmd.Rect = cmd.Rect.Add(l.Offset)
				r := cmd.Rect
				c.prepare(c.empty.tex)
				r.Max.X = r.Min.X + 1
				c.buffer.rect(r.Intersect(l.Clip), m.Vec2{}, m.Vec2{}, cmd.Color)
				r.Max.X, r.Min.X = cmd.Rect.Max.X, cmd.Rect.Max.X-1
				c.buffer.rect(r.Intersect(l.Clip), m.Vec2{}, m.Vec2{}, cmd.Color)
				r.Min.X, r.Max.Y = cmd.Rect.Min.X, r.Min.Y+1
				c.buffer.rect(r.Intersect(l.Clip), m.Vec2{}, m.Vec2{}, cmd.Color)
				r.Max.Y, r.Min.Y = cmd.Rect.Max.Y, cmd.Rect.Max.Y-1
				c.buffer.rect(r.Intersect(l.Clip), m.Vec2{}, m.Vec2{}, cmd.Color)
			case draw.Text:
				ff := c.getFontFace(cmd.Font)
				ff.write(cmd.Text, float32(l.Offset.X+cmd.Position.X), float32(l.Offset.Y+cmd.Position.Y), c, 1, l.Clip, cmd.Color)
			case draw.Shadow:
				c.buffer.flush()
				r := cmd.Rect.Add(l.Offset).Intersect(l.Clip.Inset(cmd.Size))
				c.sbuffer.rect(m.Vec2{float32(r.Min.X), float32(r.Min.Y)}, m.Vec2{float32(r.Max.X), float32(r.Max.Y)}, cmd.Rect.Add(l.Offset), float32(cmd.Size), cmd.Color)
			case draw.Icon:
				c.drawIcon(cmd.Rect.Add(l.Offset), l.Clip, cmd.Icon, cmd.Color)
			case draw.Image:
				t := c.getImage(cmd.Image, !cmd.Update)
				c.prepare(t)
				c.rect(cmd.Rect.Add(l.Offset), l.Clip, [4]float32{0, 0, 1, 1}, cmd.Color)
			}
		}
	}
	c.buffer.flush()
	c.sbuffer.flush()
	c.time++
}

func (c *Context) prepare(t uint32) {
	c.sbuffer.flush()
	if t != c.currentTexture {
		c.buffer.flush()
		gl.BindTexture(gl.TEXTURE_2D, t)
		c.currentTexture = t
	}
}

func (c *Context) rect(r, clip image.Rectangle, tr [4]float32, color draw.Color) {
	if clip.Min.X >= r.Max.X {
		return
	}
	if clip.Min.Y >= r.Max.Y {
		return
	}
	if clip.Max.X <= r.Min.X {
		return
	}
	if clip.Max.Y <= r.Min.Y {
		return
	}
	if clip.Min.X > r.Min.X {
		tr[0] += (tr[2] - tr[0]) * float32(clip.Min.X-r.Min.X) / float32(r.Dx())
		r.Min.X = clip.Min.X
	}
	if clip.Min.Y > r.Min.Y {
		tr[1] += (tr[3] - tr[1]) * (float32(clip.Min.Y-r.Min.Y) / float32(r.Dy()))
		r.Min.Y = clip.Min.Y
	}
	if clip.Max.X < r.Max.X {
		tr[2] -= (tr[2] - tr[0]) * (float32(r.Max.X-clip.Max.X) / float32(r.Dx()))
		r.Max.X = clip.Max.X
	}
	if clip.Max.Y < r.Max.Y {
		tr[3] -= (tr[3] - tr[1]) * (float32(r.Max.Y-clip.Max.Y) / float32(r.Dy()))
		r.Max.Y = clip.Max.Y
	}
	tmin := m.Vec2{tr[0], tr[1]}
	tmax := m.Vec2{tr[2], tr[3]}
	c.buffer.rect(r, tmin, tmax, color)
}
