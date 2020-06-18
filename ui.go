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
	context *imgui.Context
	io      imgui.IO
	fonts   imgui.FontAtlas
	timer   time.Time
	pic     *pixel.PictureData
	picture pixel.TargetPicture
	canvas  *pixelgl.Canvas
}

// NewUI Creates the UI and setups up its internal structures
func NewUI(context *imgui.Context, win *pixelgl.Window) *UI {
	ui := &UI{
		context: context,
	}

	ui.io = imgui.CurrentIO()
	ui.io.SetDisplaySize(pixelVecToimguiVec(win.Bounds().Size()))

	ui.fonts = ui.io.Fonts()
	ui.fonts.AddFontDefault()
	f := ui.fonts.TextureDataAlpha8()

	ui.pic = pixel.MakePictureData(pixel.R(0, 0, float64(f.Width), float64(f.Height)))

	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			i := y*f.Width + x
			ptr := (*uint8)(unsafe.Pointer(uintptr(f.Pixels) + uintptr(i)))
			ui.pic.Pix[i] = color.RGBA{R: 0, G: 0, B: 0, A: *ptr}
		}
	}

	ui.canvas = pixelgl.NewCanvas(win.Canvas().Bounds())
	ui.canvas.SetComposeMethod(pixel.ComposeOver)
	ui.canvas.SetFragmentShader(`
	#version 330 core
	in vec4  vColor;
	in vec2  vTexCoords;
	in float vIntensity;

	out vec4 fragColor;

	uniform vec4 uColorMask;
	uniform vec4 uTexBounds;
	uniform sampler2D uTexture;
	void main() {
		fragColor *= vColor * texture(uTexture, vTexCoords).a;
	}
	`)

	ui.picture = ui.canvas.MakePicture(ui.pic)
	ui.setKeyMapping()

	return ui
}

// NewFrame Call this at the beginning of the frame to tell the UI that the frame has started
func (ui *UI) NewFrame() {
	ui.timer = time.Now()
	imgui.NewFrame()
	ui.canvas.Clear(pixel.Alpha(0))
}

// update Handles general update type things and handle inputs. Called from ui.Draw.
func (ui *UI) update(win *pixelgl.Window, matrix pixel.Matrix) {
	ui.io.SetDeltaTime(float32(time.Since(ui.timer).Seconds()))

	mouse := matrix.Unproject(win.MousePosition())
	ui.io.SetMousePosition(imgui.Vec2{X: float32(mouse.X), Y: float32(mouse.Y)})

	ui.io.SetMouseButtonDown(0, win.Pressed(pixelgl.MouseButtonLeft))
	ui.io.SetMouseButtonDown(1, win.Pressed(pixelgl.MouseButtonRight))
	ui.io.SetMouseButtonDown(2, win.Pressed(pixelgl.MouseButtonMiddle))
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
	vertexSize, posOffset, uvOffset, colOffset := imgui.VertexBufferLayout()
	indexSize := imgui.IndexBufferLayout()
	for _, cmds := range data.CommandLists() {
		var indexBufferOffset uintptr
		start, _ := cmds.VertexBuffer()
		idxStart, _ := cmds.IndexBuffer()

		for _, cmd := range cmds.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(cmds)
			} else {
				win.SetMatrix(matrix)

				tris := pixel.MakeTrianglesData(cmd.ElementCount())

				for i := 0; i < cmd.ElementCount(); i++ {
					idx := unsafe.Pointer(uintptr(idxStart) + indexBufferOffset)
					index := *(*C.ushort)(idx)
					ptr := unsafe.Pointer(uintptr(start) + (uintptr(int(index) * vertexSize)))
					pos := (*imgui.Vec2)(unsafe.Pointer(uintptr(ptr) + uintptr(posOffset)))
					uv := (*imgui.Vec2)(unsafe.Pointer(uintptr(ptr) + uintptr(uvOffset)))
					col := (*uint32)(unsafe.Pointer(uintptr(ptr) + uintptr(colOffset)))

					position := imguiVecToPixelVec(*pos)
					color := imguiColorToPixelColor(*col)
					uuvv := imguiVecToPixelVec(*uv)

					(*tris)[i].Position = position
					(*tris)[i].Picture = uuvv
					(*tris)[i].Color = pixel.ToRGBA(color)
					(*tris)[i].Intensity = 0
					indexBufferOffset += uintptr(indexSize)
				}

				ui.picture.Draw(ui.canvas.MakeTriangles(tris))
			}
		}
	}
	ui.canvas.Draw(win, pixel.IM.Moved(ui.canvas.Bounds().Center()))

	win.SetMatrix(pixel.IM)
}

// imguiColorToPixelColor Converts the imgui color to a Pixel color
func imguiColorToPixelColor(c uint32) color.RGBA {
	// ABGR -> RGBA
	return color.RGBA{
		A: uint8((c >> 24) & 0xFF),
		B: uint8((c >> 16) & 0xFF),
		G: uint8((c >> 8) & 0xFF),
		R: uint8(c & 0xFF),
	}
}

// imguiVecToPixelVec Converts the imgui vector to a Pixel vector
func imguiVecToPixelVec(v imgui.Vec2) pixel.Vec {
	return pixel.V(float64(v.X), float64(v.Y))
}

// imguiRectToPixelRect Converts the imgui rect to a Pixel rect
func imguiRectToPixelRect(r imgui.Vec4) pixel.Rect {
	return pixel.R(float64(r.X), float64(r.Y), float64(r.X+r.Z), float64(r.Y+r.W))
}

func pixelVecToimguiVec(v pixel.Vec) imgui.Vec2 {
	return imgui.Vec2{X: float32(v.X), Y: float32(v.Y)}
}

func (ui *UI) setKeyMapping() {
	keys := map[int]pixelgl.Button{
		imgui.KeyTab:        pixelgl.KeyTab,
		imgui.KeyLeftArrow:  pixelgl.KeyLeft,
		imgui.KeyRightArrow: pixelgl.KeyRight,
		imgui.KeyUpArrow:    pixelgl.KeyUp,
		imgui.KeyDownArrow:  pixelgl.KeyDown,
		imgui.KeyPageUp:     pixelgl.KeyPageUp,
		imgui.KeyPageDown:   pixelgl.KeyPageDown,
		imgui.KeyHome:       pixelgl.KeyHome,
		imgui.KeyEnd:        pixelgl.KeyEnd,
		imgui.KeyInsert:     pixelgl.KeyInsert,
		imgui.KeyDelete:     pixelgl.KeyDelete,
		imgui.KeyBackspace:  pixelgl.KeyBackspace,
		imgui.KeySpace:      pixelgl.KeySpace,
		imgui.KeyEnter:      pixelgl.KeyEnter,
		imgui.KeyEscape:     pixelgl.KeyEscape,
		imgui.KeyA:          pixelgl.KeyA,
		imgui.KeyC:          pixelgl.KeyC,
		imgui.KeyV:          pixelgl.KeyV,
		imgui.KeyX:          pixelgl.KeyX,
		imgui.KeyY:          pixelgl.KeyY,
		imgui.KeyZ:          pixelgl.KeyZ,
	}

	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	for imguiKey, nativeKey := range keys {
		ui.io.KeyMap(imguiKey, int(nativeKey))
	}
}
