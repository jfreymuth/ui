package draw

// Font describes a font face.
type Font struct {
	Name string
	Size float32
}

// FontMetrics can be used to query the geometry of a font face.
type FontMetrics interface {
	Ascent() int
	Descent() int
	LineHeight() int
	Advance(s string) float32
	Index(s string, pos float32) int
}

// A FontLookup provides information about fonts.
type FontLookup interface {
	// GetClosest takes a font face description and returns the best match that is actually supported.
	// If a font has multiple names, it should return a consistent, canonical name.
	GetClosest(Font) Font
	// Metrics returns the FontMetrics for the specified font.
	// Implementations must not assume that the argument has been passed through GetClosest.
	Metrics(Font) FontMetrics
}
