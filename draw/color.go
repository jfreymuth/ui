package draw

// Color is a 32 bit RGBA color with premultiplied alpha.
type Color [4]byte

// RGBA creates a color from float RGBA values.
// The color components should be in range [0,1] regardless of the alpha value;
// this constructor will multiply them with the alpha value.
func RGBA(r, g, b, a float32) Color {
	r *= a
	g *= a
	b *= a
	var c Color
	if r >= 1 {
		c[0] = 255
	} else if r > 0 {
		c[0] = byte(r * 255)
	}
	if g >= 1 {
		c[1] = 255
	} else if g > 0 {
		c[1] = byte(g * 255)
	}
	if b >= 1 {
		c[2] = 255
	} else if b > 0 {
		c[2] = byte(b * 255)
	}
	if a >= 1 {
		c[3] = 255
	} else if a > 0 {
		c[3] = byte(a * 255)
	}
	return c
}

// Gray creates an opaque gray color from an intensity in range [0,1], where 0 means black and 1 means white.
func Gray(i float32) Color {
	b := byte(i * 255)
	if i > 1 {
		b = 255
	} else if i < 0 {
		b = 0
	}
	return Color{b, b, b, 255}
}

// R returns the color's red component
func (c Color) R() float32 { return float32(c[0]) / 255 }

// R returns the color's green component
func (c Color) G() float32 { return float32(c[1]) / 255 }

// R returns the color's blue component
func (c Color) B() float32 { return float32(c[2]) / 255 }

// R returns the color's alpha component
func (c Color) A() float32 { return float32(c[3]) / 255 }

// Blend mixes two colors. f must be in range [0,1]
func Blend(c, d Color, f float32) Color {
	r, g, b, a := c.R(), c.G(), c.B(), c.A()
	return Color{byte((r + (d.R()-r)*f) * 255), byte((g + (d.G()-g)*f) * 255), byte((b + (d.B()-b)*f) * 255), byte((a + (d.A()-a)*f) * 255)}
}

var (
	Black       = Color{0, 0, 0, 255}
	White       = Color{255, 255, 255, 255}
	Transparent = Color{0, 0, 0, 0}
)
