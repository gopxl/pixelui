package pixelui

import "github.com/faiface/pixel"

func rect(x, y, w, h float64) pixel.Rect {
	return pixel.R(x, y, x+w, y+h)
}
