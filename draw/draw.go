package draw

import "image"

// A Command represents a single drawing operation.
type Command struct {
	// Style determines how the command should be interpreted. The meaning of other fields depends on this.
	Style  byte
	Color  Color
	Bounds image.Rectangle
	Clip   image.Rectangle
	Text   string
	Data   interface{}
}

// Constants for Command.Style
const (
	Fill byte = iota
	Outline
	Text   // Command.Data must be of type Font.
	Shadow // Command.Data is the shadow's radius and must be of type int.
	Icon
	ImageStatic  // Command.Data must be of type *image.RGBA
	ImageDynamic // Command.Data must be of type *image.RGBA
)

// A Buffer contains a list of commands.
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

// Reset clears the command list and sets the size of the drawing area.
func (b *Buffer) Reset(w, h int) {
	b.Commands = b.Commands[:0]
	b.stack = b.stack[:0]
	b.bounds = WH(w, h)
	b.clip = WH(w, h)
}

// Push constrains drawing to a rectangle.
// Subsequent operations will be translated and clipped.
// Subsequent calls to Size will return the rectangle's size.
func (b *Buffer) Push(r image.Rectangle) {
	b.stack = append(b.stack, b.state)
	r = r.Add(b.bounds.Min)
	b.clip = b.clip.Intersect(r)
	b.bounds = r
}

// Pop undoes the last call to Push.
func (b *Buffer) Pop() {
	l := len(b.stack) - 1
	if l >= 0 {
		b.state, b.stack = b.stack[l], b.stack[:l]
	}
}

// Size returns the current size of the drawing area.
func (b *Buffer) Size() (int, int) {
	return b.bounds.Dx(), b.bounds.Dy()
}

// Add adds a command to the buffer.
func (b *Buffer) Add(c Command) {
	c.Bounds = c.Bounds.Add(b.bounds.Min)
	c.Clip = c.Clip.Add(b.bounds.Min).Intersect(b.clip)
	if !c.Clip.Empty() {
		b.Commands = append(b.Commands, c)
	}
}

// Fill adds a command to fill an area with a solid color.
func (b *Buffer) Fill(r image.Rectangle, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Style: Fill})
	}
}

// Outline adds a command to outline a rectangle.
func (b *Buffer) Outline(r image.Rectangle, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Style: Outline})
	}
}

// Text adds a command to render text.
func (b *Buffer) Text(r image.Rectangle, text string, c Color, font Font) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Text: text, Data: font, Style: Text})
	}
}

// Shadow adds a command to draw a drop shadow.
func (b *Buffer) Shadow(r image.Rectangle, c Color, size int) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Data: size, Style: Shadow})
	}
}

// Icon adds a command to draw an icon.
func (b *Buffer) Icon(r image.Rectangle, id string, c Color) {
	if !b.clip.Empty() {
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Text: id, Style: Icon})
	}
}

// Image adds a command to draw an image.
// Setting static may improve performance if the image's contents do not change.
func (b *Buffer) Image(r image.Rectangle, i *image.RGBA, c Color, static bool) {
	if !b.clip.Empty() {
		style := ImageDynamic
		if static {
			style = ImageStatic
		}
		b.Commands = append(b.Commands, Command{Bounds: r.Add(b.bounds.Min), Clip: b.clip, Color: c, Style: style, Data: i})
	}
}

// SubImage adds a command to draw part of an image.
// Setting static may improve performance if the image's contents do not change.
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
