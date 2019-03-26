package gldraw

import (
	"hash/crc32"
	"image"
)

type entry struct {
	t        *Texture
	lastUsed uint
	checksum uint32
}

func (c *Context) getImage(i *image.RGBA, static bool) uint32 {
	if e, ok := c.images[i]; ok {
		e.lastUsed = c.time
		if static {
			return e.t.tex
		} else {
			ch := checksum(i)
			if ch != e.checksum {
				e.t.Update(i, image.Point{})
				e.checksum = ch
			}
			return e.t.tex
		}
	}
	t := NewTexture(i)
	var ch uint32
	if !static {
		ch = checksum(i)
	}
	c.images[i] = &entry{t, c.time, ch}
	return t.tex
}

func checksum(i *image.RGBA) uint32 {
	return crc32.ChecksumIEEE(i.Pix[:i.Rect.Dx()*i.Rect.Dy()*4])
}
