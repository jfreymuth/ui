package draw

import "image"

func WH(w, h int) image.Rectangle {
	return image.Rect(0, 0, w, h)
}
func XYXY(x, y, x1, y1 int) image.Rectangle {
	return image.Rect(x, y, x1, y1)
}
func XYWH(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}
