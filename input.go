package pixelui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func JustPressed(button pixelgl.Button) bool {
	return false
}

func JustReleased(button pixelgl.Button) bool {
	return false
}

func mouseIn(r pixel.Rect) bool {
	return r.Contains(ui.win.MousePosition())
}
