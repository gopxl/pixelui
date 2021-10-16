package pixelui

import "github.com/faiface/pixel"

func rect(x, y, w, h float64) pixel.Rect {
	return pixel.R(x, y, x+w, y+h)
}

func rectEq(r1, r2 pixel.Rect) bool {
	return r1.Min.Eq(r2.Min) && r1.Max.Eq(r2.Max)
}

func rectRecenter(r pixel.Rect) pixel.Rect {
	return r.Moved(r.Min.Scaled(-1))
}
