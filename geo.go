package pixelui

import (
	"github.com/faiface/pixel"
	"github.com/inkyblackness/imgui-go/v4"
)

// imguiRectToPixelRect Converts the imgui rect to a Pixel rect
func imguiRectToPixelRect(r imgui.Vec4) pixel.Rect {
	return pixel.R(float64(r.X), float64(r.Y), float64(r.Z), float64(r.W))
}

// IVec converts a pixel vector to an imgui vector
func IVec(v pixel.Vec) imgui.Vec2 {
	return imgui.Vec2{X: float32(v.X), Y: float32(v.Y)}
}

// IV creates an imgui vector from the given points.
func IV(x, y float64) imgui.Vec2 {
	return imgui.Vec2{X: float32(x), Y: float32(y)}
}

// PV converts an imgui vector to a pixel vector
func PV(v imgui.Vec2) pixel.Vec {
	return pixel.V(float64(v.X), float64(v.Y))
}

// ProjectVec projects the vector by the UI's matrix (vertical flip)
// 	and returns that as a imgui vector
func ProjectVec(v pixel.Vec) imgui.Vec2 {
	return IVec(CurrentUI.matrix.Project(v))
}

// ProjectV creates a pixel vector and projects it using ProjectVec
func ProjectV(x, y float64) imgui.Vec2 {
	return ProjectVec(pixel.V(x, y))
}

// UnprojectV unprojects the vector by the UI's matrix (vertical flip)
// 	and returns that as a pixel vector
func UnprojectV(v imgui.Vec2) pixel.Vec {
	return CurrentUI.matrix.Unproject(PV(v))
}

// IZV returns an imgui zero vector
func IZV() imgui.Vec2 {
	return imgui.Vec2{X: 0, Y: 0}
}
