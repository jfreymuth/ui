package gofont

import (
	"strings"

	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/impl/gldraw"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

var (
	fonts map[string]entry
	fm    map[draw.Font]*metrics
)

type entry struct {
	realName string
	ttf      *truetype.Font
}

type lookup struct {
	dpi float64
}

func initFonts() {
	fonts = make(map[string]entry)
	ttf, _ := truetype.Parse(goregular.TTF)
	fonts["goregular"] = entry{"goregular", ttf}
	ttf, _ = truetype.Parse(goitalic.TTF)
	fonts["goitalic"] = entry{"goitalic", ttf}
	ttf, _ = truetype.Parse(gobold.TTF)
	fonts["gobold"] = entry{"gobold", ttf}
	ttf, _ = truetype.Parse(gobolditalic.TTF)
	fonts["gobolditalic"] = entry{"gobolditalic", ttf}
	ttf, _ = truetype.Parse(gomono.TTF)
	fonts["gomono"] = entry{"gomono", ttf}
	ttf, _ = truetype.Parse(gomonoitalic.TTF)
	fonts["gomonoitalic"] = entry{"gomonoitalic", ttf}
	ttf, _ = truetype.Parse(gomonobold.TTF)
	fonts["gomonobold"] = entry{"gomonobold", ttf}
	ttf, _ = truetype.Parse(gomonobolditalic.TTF)
	fonts["gomonobolditalic"] = entry{"gomonobolditalic", ttf}
}

func Lookup(dpi float32) gldraw.FontLookup {
	return &lookup{dpi: float64(dpi)}
}

func Add(name string, ttf *truetype.Font) {
	if fonts == nil {
		initFonts()
	}
	fonts[name] = entry{name, ttf}
}

func AddAlias(name string, alias string) {
	if fonts == nil {
		initFonts()
	}
	if e, ok := fonts[name]; ok {
		fonts[alias] = e
	}
}

func (f *lookup) GetClosest(font draw.Font) draw.Font {
	if fonts == nil {
		initFonts()
	}
	if e, ok := fonts[font.Name]; ok {
		font.Name = e.realName
	} else {
		font.Name = fallback(font.Name)
	}
	if font.Size < 5 {
		font.Size = 5
	} else if font.Size > 72 {
		font.Size = 72
	} else {
		font.Size = float32(int(font.Size*2+.5)) / 2
	}
	return font
}

func (f *lookup) LoadFont(font draw.Font) font.Face {
	if fonts == nil {
		initFonts()
	}
	e, ok := fonts[font.Name]
	if !ok {
		font = f.GetClosest(font)
		e = fonts[font.Name]
	}
	return truetype.NewFace(e.ttf, &truetype.Options{Size: float64(font.Size), DPI: f.dpi, GlyphCacheEntries: 1})
}

func (f *lookup) Metrics(font draw.Font) draw.FontMetrics {
	font = f.GetClosest(font)
	if m, ok := fm[font]; ok {
		return m
	}
	face := f.LoadFont(font)
	m := &metrics{face, face.Metrics()}
	if fm == nil {
		fm = make(map[draw.Font]*metrics)
	}
	fm[font] = m
	return m
}

func fallback(name string) (goName string) {
	if i := strings.Index(name, "mono"); i != -1 {
		goName += "mono"
		name = name[:i] + name[i+4:]
	}
	if i := strings.Index(name, "bold"); i != -1 {
		goName += "bold"
		name = name[:i] + name[i+4:]
	}
	if i := strings.Index(name, "italic"); i != -1 {
		goName += "italic"
		name = name[:i] + name[i+6:]
	}
	if goName == "" {
		return "goregular"
	}
	return "go" + goName
}

type metrics struct {
	font font.Face
	m    font.Metrics
}

func (m *metrics) Ascent() int     { return m.m.Ascent.Ceil() }
func (m *metrics) Descent() int    { return m.m.Descent.Ceil() }
func (m *metrics) LineHeight() int { return m.m.Height.Ceil() }

func (m *metrics) Advance(s string) float32 {
	var x fixed.Int26_6
	var last rune
	for _, r := range s {
		x += m.font.Kern(last, r)
		adv, _ := m.font.GlyphAdvance(r)
		x += adv
		last = r
	}
	return float32(x) / 64
}

func (m *metrics) Index(s string, t float32) int {
	var x fixed.Int26_6
	var last rune
	for i, r := range s {
		x += m.font.Kern(last, r)
		adv, _ := m.font.GlyphAdvance(r)
		if t < float32(x+adv/2)/64 {
			return i
		}
		x += adv
		last = r
	}
	return len(s)
}
