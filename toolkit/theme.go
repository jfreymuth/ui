package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type Theme struct {
	fonts  map[string]draw.Font
	colors map[string]draw.Color
	parent *Theme
}

func SetTheme(c ui.Component, theme *Theme) {
	if c, ok := c.(interface{ SetTheme(*Theme) }); ok {
		c.SetTheme(theme)
	}
}

func (t *Theme) New() *Theme {
	return &Theme{parent: t}
}

func (t *Theme) Font(name string) draw.Font {
	if f, ok := t.fonts[name]; ok {
		return f
	}
	if t.parent != nil {
		return t.parent.Font(name)
	}
	return draw.Font{Name: "default", Size: 12}
}

func (t *Theme) Color(name string) draw.Color {
	if c, ok := t.colors[name]; ok {
		return c
	}
	if t.parent != nil {
		return t.parent.Color(name)
	}
	return draw.Transparent
}

func (t *Theme) SetFont(name string, f draw.Font) {
	if t.fonts == nil {
		t.fonts = make(map[string]draw.Font)
	}
	t.fonts[name] = f
}

func (t *Theme) SetColor(name string, c draw.Color) {
	if t.colors == nil {
		t.colors = make(map[string]draw.Color)
	}
	t.colors[name] = c
}

var DefaultTheme = LightTheme
var LightTheme = &Theme{
	fonts: map[string]draw.Font{
		"text":       draw.Font{Name: "default", Size: 12},
		"title":      draw.Font{Name: "bold", Size: 12},
		"buttonText": draw.Font{Name: "bold", Size: 12},
		"inputText":  draw.Font{Name: "default", Size: 12},
	},
	colors: map[string]draw.Color{
		"background":           draw.White,
		"border":               draw.Black,
		"shadow":               draw.RGBA(0, 0, 0, .5),
		"veil":                 draw.RGBA(0, 0, 0, .2),
		"altBackground":        draw.Gray(.9),
		"text":                 draw.Black,
		"title":                draw.Black,
		"titleBackground":      draw.RGBA(.7, .75, 1, 1),
		"titleBackgroundError": draw.RGBA(1, .75, .7, 1),
		"buttonBackground":     draw.Transparent,
		"buttonHovered":        draw.RGBA(0, 0, 0, .2),
		"buttonText":           draw.Black,
		"buttonFocused":        draw.Gray(.3),
		"inputBackground":      draw.White,
		"inputText":            draw.Black,
		"selection":            draw.RGBA(.8, .85, 1, 1),
		"selectionInactive":    draw.Gray(.8),
		"scrollBar":            draw.RGBA(0, 0, 0, .3),
	},
}
var DarkTheme = &Theme{
	fonts: map[string]draw.Font{
		"text":       draw.Font{Name: "default", Size: 12},
		"title":      draw.Font{Name: "bold", Size: 12},
		"buttonText": draw.Font{Name: "bold", Size: 12},
		"inputText":  draw.Font{Name: "default", Size: 12},
	},
	colors: map[string]draw.Color{
		"background":           draw.Gray(.3),
		"border":               draw.Black,
		"shadow":               draw.Black,
		"veil":                 draw.RGBA(1, 1, 1, .2),
		"altBackground":        draw.Gray(.4),
		"text":                 draw.White,
		"title":                draw.White,
		"titleBackground":      draw.RGBA(.2, .25, .4, 1),
		"titleBackgroundError": draw.RGBA(.5, .05, 0, 1),
		"buttonBackground":     draw.Transparent,
		"buttonHovered":        draw.RGBA(1, 1, 1, .1),
		"buttonText":           draw.White,
		"buttonFocused":        draw.Gray(.8),
		"inputBackground":      draw.Gray(.3),
		"inputText":            draw.White,
		"selection":            draw.RGBA(.35, .4, .6, 1),
		"selectionInactive":    draw.Gray(.5),
		"scrollBar":            draw.RGBA(1, 1, 1, .3),
	},
}
