package pixelui

import "github.com/inkyblackness/imgui-go/v4"

// Color converts the given 8-bit r,g,b components to a imgui.Vec4 for color arguments
func Color(r, g, b uint8) imgui.Vec4 {
	return ColorA(r, g, b, 255)
}

// Color converts the given 8-bit r,g,b,a components to a imgui.Vec4 for color arguments
func ColorA(r, g, b, a uint8) imgui.Vec4 {
	var scale float32 = 255
	return imgui.Vec4{
		X: float32(r) / scale,
		Y: float32(g) / scale,
		Z: float32(b) / scale,
		W: float32(a) / scale,
	}
}
