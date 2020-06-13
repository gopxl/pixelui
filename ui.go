package pixelui

import (
	"C"

	"github.com/faiface/pixel"
	"github.com/inkyblackness/imgui-go"
)
import (
	"image/color"
	"time"
	"unsafe"

	"github.com/faiface/pixel/pixelgl"
)

// UI Stores the state of the pixelui UI
type UI struct {
	tris    *pixel.TrianglesData
	batch   *pixel.Batch
	context *imgui.Context
	io      imgui.IO
	fonts   imgui.FontAtlas
	timer   time.Time
}

// NewUI Creates the UI and setups up its internal structures
func NewUI(context *imgui.Context) *UI {
	ui := &UI{
		context: context,
	}

	ui.tris = pixel.MakeTrianglesData(0)
	ui.batch = pixel.NewBatch(ui.tris, nil)

	ui.io = imgui.CurrentIO()
	ui.io.SetDisplaySize(imgui.Vec2{X: 1920, Y: 1080})

	ui.fonts = ui.io.Fonts()
	ui.fonts.AddFontDefault()
	ui.fonts.TextureDataAlpha8()

	return ui
}

// NewFrame Call this at the beginning of the frame to tell the UI that the frame has started
func (ui *UI) NewFrame() {
	ui.timer = time.Now()
}

// update Handles general update type things and handle inputs. Called from ui.Draw.
func (ui *UI) update(win *pixelgl.Window, matrix pixel.Matrix) {
	ui.io.SetDeltaTime(float32(time.Since(ui.timer).Seconds()))

	mouse := matrix.Unproject(win.MousePosition())
	ui.io.SetMousePosition(imgui.Vec2{X: float32(mouse.X), Y: float32(mouse.Y)})

	ui.io.SetMouseButtonDown(0, win.Pressed(pixelgl.MouseButtonLeft))
	ui.io.SetMouseButtonDown(1, win.Pressed(pixelgl.MouseButtonMiddle))
	ui.io.SetMouseButtonDown(2, win.Pressed(pixelgl.MouseButtonRight))
	ui.io.AddMouseWheelDelta(float32(win.MouseScroll().X), float32(win.MouseScroll().Y))
}

// Draw Draws the imgui UI to the Pixel Window
func (ui *UI) Draw(win *pixelgl.Window) {
	// imgui draws things from top-left as 0,0 where Pixel draws from bottom-left as 0,0,
	//	for drawing and handling inputs, we need to "flip" imgui.
	matrix := pixel.IM.ScaledXY(win.Bounds().Center(), pixel.V(1, -1))
	ui.update(win, matrix)

	// Tell imgui to render and get the resulting draw data
	imgui.Render()
	data := imgui.RenderedDrawData()

	// In each command, there is a vertex buffer that holds all of the vertices to draw;
	// 	there's also an index buffer which stores the indices into the vertex buffer that should
	//	be draw together. The vertex buffer is shared between multiple commands.
	vertexSize, posOffset, _, colOffset := imgui.VertexBufferLayout()
	// vertexSize, posOffset, uvOffset, colOffset := imgui.VertexBufferLayout()
	for _, cmds := range data.CommandLists() {
		start, _ := cmds.VertexBuffer()
		idxStart, _ := cmds.IndexBuffer()

		for _, cmd := range cmds.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(cmds)
			} else {
				ui.tris.SetLen(cmd.ElementCount() + ui.tris.Len())
				jndex := 0
				for i := 0; i < cmd.ElementCount(); i++ {
					idx := unsafe.Pointer(uintptr(idxStart) + uintptr(i*imgui.IndexBufferLayout()))
					index := *(*C.ushort)(idx)
					ptr := unsafe.Pointer(uintptr(start) + (uintptr(int(index) * vertexSize)))
					pos := (*imgui.Vec2)(unsafe.Pointer(uintptr(ptr) + uintptr(posOffset)))
					// uv := (*imgui.Vec2)(unsafe.Pointer(uintptr(ptr) + uintptr(uvOffset)))
					col := (*uint32)(unsafe.Pointer(uintptr(ptr) + uintptr(colOffset)))

					position := imguiVecToPixelVec(*pos)
					color := imguiColorToPixelColor(*col)
					// uuvv := imguiVecToPixelVec(*uv)

					(*ui.tris)[jndex].Position = position
					// (*ui.tris)[jndex].Picture = uuvv
					(*ui.tris)[jndex].Color = color
					(*ui.tris)[jndex].Intensity = 1.0
					jndex++
				}
			}
		}

	}
	ui.batch.Dirty()

	win.SetMatrix(matrix)

	ui.batch.Draw(win)
	ui.tris.SetLen(0)

	win.SetMatrix(pixel.IM)
}

// imguiColorToPixelColor Converts the imgui color to a Pixel color
func imguiColorToPixelColor(c uint32) pixel.RGBA {
	return pixel.ToRGBA(color.RGBA{
		R: uint8((c >> 24) & 0xFF),
		G: uint8((c >> 16) & 0xFF),
		B: uint8((c >> 8) & 0xFF),
		A: uint8(c & 0xFF),
	})
}

// imguiVecToPixelVec Converts the imgui vector to a Pixel vector
func imguiVecToPixelVec(v imgui.Vec2) pixel.Vec {
	return pixel.V(float64(v.X), float64(v.Y))
}
