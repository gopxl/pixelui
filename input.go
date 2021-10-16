package pixelui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	held     [pixelgl.KeyLast + 1]bool
	pressed  [pixelgl.KeyLast + 1]bool
	released [pixelgl.KeyLast + 1]bool
)

func consumeHeld(button pixelgl.Button) bool {
	if held[button] {
		held[button] = false
		return true
	}
	return false
}

func consumePressed(button pixelgl.Button) bool {
	if pressed[button] {
		pressed[button] = false
		return true
	}
	return false
}

func consumeReleased(button pixelgl.Button) bool {
	if released[button] {
		released[button] = false
		return true
	}
	return false
}

func mouseIn(r pixel.Rect) bool {
	return r.Contains(ui.win.MousePosition())
}
