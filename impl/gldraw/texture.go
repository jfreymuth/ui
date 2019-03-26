package gldraw

import (
	"image"
	"image/draw"
	"io"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	c             *Context
	tex           uint32
	width, height int
	alpha         bool
}

func NewTexture(img *image.RGBA) *Texture       { t := new(Texture); t.init(img); return t }
func NewTextureEmpty(w, h int) *Texture         { t := new(Texture); t.initEmpty(w, h); return t }
func NewTextureAlpha(img *image.Alpha) *Texture { t := new(Texture); t.initAlpha(img); return t }
func NewTextureAlphaEmpty(w, h int) *Texture    { t := new(Texture); t.initAlphaEmpty(w, h); return t }

func (t *Texture) Bounds() image.Rectangle { return image.Rect(0, 0, t.width, t.height) }

func (t *Texture) initBase(w, h int) {
	t.width = w
	t.height = h
	gl.GenTextures(1, &t.tex)
	gl.BindTexture(gl.TEXTURE_2D, t.tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
}

func (t *Texture) initBaseAlpha(w, h int) {
	t.alpha = true
	t.initBase(w, h)
	swizzle := [...]int32{gl.RED, gl.RED, gl.RED, gl.RED}
	gl.TexParameteriv(gl.TEXTURE_2D, gl.TEXTURE_SWIZZLE_RGBA, &swizzle[0])
}

func (t *Texture) init(img *image.RGBA) {
	var old int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &old)
	t.initBase(img.Rect.Dx(), img.Rect.Dy())
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, int32(img.Stride/4))
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Dx()), int32(img.Rect.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.BindTexture(gl.TEXTURE_2D, uint32(old))
}

func (t *Texture) initEmpty(w, h int) {
	var old int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &old)
	t.initBase(w, h)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.BindTexture(gl.TEXTURE_2D, uint32(old))
}

func (t *Texture) initAlpha(img *image.Alpha) {
	var old int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &old)
	t.initBaseAlpha(img.Rect.Dx(), img.Rect.Dy())
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, int32(img.Stride))
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(img.Rect.Dx()), int32(img.Rect.Dy()), 0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.BindTexture(gl.TEXTURE_2D, uint32(old))
}

func (t *Texture) initAlphaEmpty(w, h int) {
	var old int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &old)
	t.initBaseAlpha(w, h)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(w), int32(h), 0, gl.RED, gl.UNSIGNED_BYTE, nil)
	gl.BindTexture(gl.TEXTURE_2D, uint32(old))
}

func (t *Texture) Update(img *image.RGBA, p image.Point) {
	if t.alpha {
		panic("wrong texture format")
	}
	if t.c != nil {
		if t.c.currentTexture == t.tex {
			t.c.buffer.flush()
		} else {
			t.c.prepare(t.tex)
		}
	}
	r := image.Rect(0, 0, t.width, t.height).Intersect(img.Rect.Add(p))
	if r.Empty() {
		return
	}
	xo, yo := 0, 0
	if p.X < 0 {
		xo = -p.X
	}
	if p.Y < 0 {
		yo = -p.Y
	}
	gl.BindTexture(gl.TEXTURE_2D, t.tex)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, int32(img.Stride/4))
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, int32(r.Min.X), int32(r.Min.Y), int32(r.Dx()), int32(r.Dy()), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix[(xo+yo*img.Stride)*4:]))
}

func (t *Texture) UpdateAlpha(img *image.Alpha, p image.Point) {
	if !t.alpha {
		panic("wrong texture format")
	}
	if t.c != nil {
		if t.c.currentTexture == t.tex {
			t.c.buffer.flush()
		} else {
			t.c.prepare(t.tex)
		}
	}
	r := image.Rect(0, 0, t.width, t.height).Intersect(img.Rect.Add(p))
	if r.Empty() {
		return
	}
	xo, yo := 0, 0
	if p.X < 0 {
		xo = -p.X
	}
	if p.Y < 0 {
		yo = -p.Y
	}
	gl.BindTexture(gl.TEXTURE_2D, t.tex)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, int32(img.Stride))
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, int32(r.Min.X), int32(r.Min.Y), int32(r.Dx()), int32(r.Dy()), gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix[(xo+yo*img.Stride):]))
}

func TextureFromFile(filename string) (*Texture, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	t, err := TextureFromReader(file)
	file.Close()
	return t, err
}

func TextureFromReader(r io.Reader) (*Texture, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	if rgba, ok := img.(*image.RGBA); ok {
		return NewTexture(rgba), nil
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Src)
	return NewTexture(rgba), nil
}
