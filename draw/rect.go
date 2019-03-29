package draw

import "image"

// WH creates a rectangle at (0,0) with the given size.
func WH(w, h int) image.Rectangle {
	return image.Rect(0, 0, w, h)
}

// XYXY creates a rectangle with the given minimum and maximum points.
func XYXY(x, y, x1, y1 int) image.Rectangle {
	return image.Rect(x, y, x1, y1)
}

// XYWH creates a rectangle at the given point with the given size.
func XYWH(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}
