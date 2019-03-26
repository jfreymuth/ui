package draw

type Font struct {
	Name string
	Size float32
}

type FontMetrics interface {
	Ascent() int
	Descent() int
	LineHeight() int
	Advance(s string) float32
	Index(s string, pos float32) int
}

type FontLookup interface {
	GetClosest(Font) Font
	Metrics(Font) FontMetrics
}
