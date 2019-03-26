package draw

import "image"

type Command struct {
	Style  byte
	Color  Color
	Bounds image.Rectangle
	Clip   image.Rectangle
	Text   string
	Data   interface{}
}

const (
	Fill byte = iota
	Outline
	Text
	Shadow
	Icon
	ImageStatic
	ImageDynamic
)

type Buffer struct {
	Commands   []Command
	FontLookup FontLookup
	state
	stack []state
}

type state struct {
	clip   image.Rectangle
	bounds image.Rectangle
}

func (b *Buffer) Reset(w, h int) {
	b.Commands = b.Commands[:0]
	b.stack = b.stack[:0]
	b.bounds = WH(w, h)
	b.clip = WH(w, h)
}

func (b *Buffer) Push(r image.Rectangle) {
	b.stack = append(b.stack, b.state)
	r = r.Add(b.bounds.Min)
	b.clip = b.clip.Intersect(r)
	b.bounds = r
}

func (b *Buffer) Pop() {
	l := len(b.stack) - 1
	if l >= 0 {
		b.state, b.stack = b.stack[l], b.stack[:l]
	}
}

func (b *Buffer) Size() (int, int) {
	return b.bounds.Dx(), b.bounds.Dy()
}

func (b *Buffer) Add(c Command) {
	c.Bounds = c.Bounds.Add(b.bounds.Min)
	c.Clip = c.Clip.Add(b.bounds.Min).Intersect(b.clip)
	if !c.Clip.Empty() {
		b.Commands = append(b.Commands, c)
	}
}

func (b *Buffer) Fill(r image.Rectangle, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Style: Fill})
	}
}

func (b *Buffer) Outline(r image.Rectangle, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Style: Outline})
	}
}

func (b *Buffer) Text(r image.Rectangle, text string, c Color, font Font) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Text: text, Data: font, Style: Text})
	}
}

func (b *Buffer) Shadow(r image.Rectangle, c Color, size int) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Data: size, Style: Shadow})
	}
}

func (b *Buffer) Icon(r image.Rectangle, id string, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Text: id, Style: Icon})
	}
}

func (b *Buffer) Image(r image.Rectangle, i *image.RGBA, c Color, static bool) {
	if !b.clip.Empty() {
		style := ImageDynamic
		if static {
			style = ImageStatic
		}
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Style: style, Data: i})
	}
}

func (b *Buffer) SubImage(r image.Rectangle, i *image.RGBA, sub image.Rectangle, c Color, static bool) {
	if sub.Empty() || b.clip.Empty() {
		return
	}
	style := ImageDynamic
	if static {
		style = ImageStatic
	}
	clip := r.Add(b.bounds.Min)
	w := r.Dx() * i.Rect.Dx() / sub.Dx()
	h := r.Dy() * i.Rect.Dy() / sub.Dy()
	x := r.Min.X - sub.Min.X*r.Dx()/sub.Dx()
	y := r.Min.Y - sub.Min.Y*r.Dy()/sub.Dy()
	b.Commands = append(b.Commands, Command{Bounds: XYWH(b.bounds.Min.X+x, b.bounds.Min.Y+y, w, h), Clip: b.clip.Intersect(clip), Color: c, Style: style, Data: i})
}
