package draw

import "image"

type Command interface{}

type Fill struct {
	Rect  image.Rectangle
	Color Color
}

type Outline struct {
	Rect  image.Rectangle
	Color Color
}

type Text struct {
	Position image.Point
	Text     string
	Font     Font
	Color    Color
}

type Shadow struct {
	Rect  image.Rectangle
	Color Color
	Size  int
}

type Icon struct {
	Rect  image.Rectangle
	Icon  string
	Color Color
}

type Image struct {
	Rect   image.Rectangle
	Image  *image.RGBA
	Color  Color
	Update bool
}

type CommandList struct {
	Commands []Command
	Offset   image.Point
	Clip     image.Rectangle
}

// A Buffer contains a list of commands.
type Buffer struct {
	Commands   []Command
	FontLookup FontLookup
	All        []CommandList
	state
	stack []state
}

type state struct {
	clip   image.Rectangle
	bounds image.Rectangle
}

// Reset clears the command list and sets the size of the drawing area.
func (b *Buffer) Reset(w, h int) {
	b.All = b.All[:0]
	b.Commands = nil
	b.stack = b.stack[:0]
	b.bounds = WH(w, h)
	b.clip = WH(w, h)
}

// Push constrains drawing to a rectangle.
// Subsequent operations will be translated and clipped.
// Subsequent calls to Size will return the rectangle's size.
func (b *Buffer) Push(r image.Rectangle) {
	b.flush()
	b.stack = append(b.stack, b.state)
	r = r.Add(b.bounds.Min)
	b.clip = b.clip.Intersect(r)
	b.bounds = r
}

// Pop undoes the last call to Push.
func (b *Buffer) Pop() {
	b.flush()
	l := len(b.stack) - 1
	if l >= 0 {
		b.state, b.stack = b.stack[l], b.stack[:l]
	}
}

// Size returns the current size of the drawing area.
func (b *Buffer) Size() (int, int) {
	return b.bounds.Dx(), b.bounds.Dy()
}

// Add adds commands to the buffer.
func (b *Buffer) Add(c ...Command) {
	b.Commands = append(b.Commands, c...)
}

// Fill adds a command to fill an area with a solid color.
func (b *Buffer) Fill(r image.Rectangle, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Fill{Rect: r, Color: c})
	}
}

// Outline adds a command to outline a rectangle.
func (b *Buffer) Outline(r image.Rectangle, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Outline{Rect: r, Color: c})
	}
}

// Text adds a command to render text.
func (b *Buffer) Text(p image.Point, text string, c Color, font Font) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Text{Position: p, Color: c, Text: text, Font: font})
	}
}

// Shadow adds a command to draw a drop shadow.
func (b *Buffer) Shadow(r image.Rectangle, c Color, size int) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Shadow{Rect: r, Color: c, Size: size})
	}
}

// Icon adds a command to draw an icon.
func (b *Buffer) Icon(r image.Rectangle, id string, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Icon{Rect: r, Color: c, Icon: id})
	}
}

// Image adds a command to draw an image.
// update should be true if the image has changed since it was last drawn.
func (b *Buffer) Image(r image.Rectangle, i *image.RGBA, c Color, update bool) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Image{Rect: r, Color: c, Image: i})
	}
}

// SubImage adds a command to draw part of an image.
// update should be true if the image has changed since it was last drawn.
func (b *Buffer) SubImage(r image.Rectangle, i *image.RGBA, sub image.Rectangle, c Color, update bool) {
	if sub.Empty() || b.clip.Empty() {
		return
	}
	w := r.Dx() * i.Rect.Dx() / sub.Dx()
	h := r.Dy() * i.Rect.Dy() / sub.Dy()
	x := r.Min.X - sub.Min.X*r.Dx()/sub.Dx()
	y := r.Min.Y - sub.Min.Y*r.Dy()/sub.Dy()
	b.Push(r)
	b.Commands = append(b.Commands, Image{Rect: XYWH(x, y, w, h), Color: c, Image: i, Update: update})
	b.Pop()
}

func (b *Buffer) flush() {
	if len(b.Commands) > 0 {
		b.All = append(b.All, CommandList{b.Commands, b.bounds.Min, b.clip})
		b.Commands = nil
	}
}
