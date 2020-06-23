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

	"github.com/faiface/mainthread"
	"github.com/faiface/pixel/pixelgl"
)

const uiShader = `
#version 330 core
in vec4  vColor;
in vec2  vTexCoords;
in float vIntensity;

out vec4 fragColor;

uniform vec4 uColorMask;
uniform vec4 uTexBounds;
uniform sampler2D uTexture;
uniform vec4 uClipRect;

void main() {
	fragColor *= vColor * texture(uTexture, vTexCoords).a;
}
`

// UI Stores the state of the pixelui UI
type UI struct {
	win     *pixelgl.Window
	context *imgui.Context
	io      imgui.IO
	fonts   imgui.FontAtlas
	timer   time.Time
	pic     *pixel.PictureData
	picture pixel.TargetPicture
	shader  *pixelgl.GLShader
	matrix  pixel.Matrix
}

// NewUI Creates the UI and setups up its internal structures
func NewUI(win *pixelgl.Window) *UI {
	var context *imgui.Context
	mainthread.Call(func() {
		context = imgui.CreateContext(nil)
	})

	ui := &UI{
		win:     win,
		context: context,
	}

	ui.matrix = pixel.IM.ScaledXY(win.Bounds().Center(), pixel.V(1, -1))

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

	ui.shader = pixelgl.NewGLShader(uiShader)

	ui.picture = win.Canvas().MakePicture(ui.pic)
	ui.setKeyMapping()

	return ui
}

// Destroy cleans up the imgui context
func (ui *UI) Destroy() {
	ui.context.Destroy()
}

// NewFrame Call this at the beginning of the frame to tell the UI that the frame has started
func (ui *UI) NewFrame() {
	ui.timer = time.Now()
	ui.io.AddMouseWheelDelta(float32(ui.win.MouseScroll().X), float32(ui.win.MouseScroll().Y))
	mouse := ui.matrix.Unproject(ui.win.MousePosition())
	ui.io.SetMousePosition(imgui.Vec2{X: float32(mouse.X), Y: float32(mouse.Y)})

	ui.io.SetMouseButtonDown(0, ui.win.Pressed(pixelgl.MouseButtonLeft))
	ui.io.SetMouseButtonDown(1, ui.win.Pressed(pixelgl.MouseButtonRight))
	ui.io.SetMouseButtonDown(2, ui.win.Pressed(pixelgl.MouseButtonMiddle))
	ui.io.AddInputCharacters(ui.win.Typed())

	for _, key := range keys {
		if ui.win.Pressed(key) {
			ui.io.KeyPress(int(key))
		} else {
			ui.io.KeyRelease(int(key))
		}
		ui.updateKeyMod()
	}
	imgui.NewFrame()
}

func (ui *UI) mapModifier(lKey pixelgl.Button, rKey pixelgl.Button) (lResult int, rResult int) {
	if ui.win.Pressed(lKey) {
		lResult = 1
	}
	if ui.win.Pressed(rKey) {
		rResult = 1
	}
	return
}

func (ui *UI) updateKeyMod() {
	ui.io.KeyCtrl(ui.mapModifier(pixelgl.KeyLeftControl, pixelgl.KeyRightControl))
	ui.io.KeyShift(ui.mapModifier(pixelgl.KeyLeftShift, pixelgl.KeyRightShift))
	ui.io.KeyAlt(ui.mapModifier(pixelgl.KeyLeftAlt, pixelgl.KeyRightAlt))
	ui.io.KeySuper(ui.mapModifier(pixelgl.KeyLeftSuper, pixelgl.KeyRightSuper))
}

// update Handles general update type things and handle inputs. Called from ui.Draw.
func (ui *UI) update() {
	ui.io.SetDeltaTime(float32(time.Since(ui.timer).Seconds()))
}

// inputWant is a helper for determining what type a button is: keyboard/mouse
func (ui *UI) inputWant(button pixelgl.Button) bool {
	switch button {
	case pixelgl.MouseButton1, pixelgl.MouseButton2, pixelgl.MouseButton3, pixelgl.MouseButton4, pixelgl.MouseButton5, pixelgl.MouseButton6, pixelgl.MouseButton7, pixelgl.MouseButton8:
		return ui.io.WantCaptureMouse()
	}
	return ui.io.WantCaptureKeyboard()
}

// JustPressed returns true if imgui hasn't handled the button and the button was just pressed
func (ui *UI) JustPressed(button pixelgl.Button) bool {
	return !ui.inputWant(button) && ui.win.JustPressed(button)
}

// JustPressed returns true if imgui hasn't handled the button and the button was just released
func (ui *UI) JustReleased(button pixelgl.Button) bool {
	return !ui.inputWant(button) && ui.win.JustReleased(button)
}

// JustPressed returns true if imgui hasn't handled the button and the button is pressed
func (ui *UI) Pressed(button pixelgl.Button) bool {
	return !ui.inputWant(button) && ui.win.Pressed(button)
}

// Repeated returns true if imgui hasn't handled the button and the button was repeated
func (ui *UI) Repeated(button pixelgl.Button) bool {
	return !ui.inputWant(button) && ui.win.Repeated(button)
}

// Draw Draws the imgui UI to the Pixel Window
func (ui *UI) Draw(win *pixelgl.Window) {
	win.SetComposeMethod(pixel.ComposeOver)

	// imgui draws things from top-left as 0,0 where Pixel draws from bottom-left as 0,0,
	//	for drawing and handling inputs, we need to "flip" imgui.
	ui.update()

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
				win.SetMatrix(ui.matrix)
				tris := pixel.MakeTrianglesData(cmd.ElementCount())

				for i := 0; i < cmd.ElementCount(); i++ {
					idx := unsafe.Pointer(uintptr(idxStart) + indexBufferOffset)
					index := *(*uint16)(idx)
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

				clipRect := imguiRectToPixelRect(cmd.ClipRect())
				clipRect.Min = ui.matrix.Project(clipRect.Min)
				clipRect.Max = ui.matrix.Project(clipRect.Max)
				shaderTris := pixelgl.NewGLTriangles(ui.shader, tris)
				shaderTris.SetClipRect(clipRect)
				ui.picture.Draw(win.Canvas().MakeTriangles(shaderTris))
			}
		}
	}

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
	return pixel.R(float64(r.X), float64(r.Y), float64(r.Z), float64(r.W))
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

var keys = []pixelgl.Button{
	pixelgl.KeySpace,
	pixelgl.KeyApostrophe,
	pixelgl.KeyComma,
	pixelgl.KeyMinus,
	pixelgl.KeyPeriod,
	pixelgl.KeySlash,
	pixelgl.Key0,
	pixelgl.Key1,
	pixelgl.Key2,
	pixelgl.Key3,
	pixelgl.Key4,
	pixelgl.Key5,
	pixelgl.Key6,
	pixelgl.Key7,
	pixelgl.Key8,
	pixelgl.Key9,
	pixelgl.KeySemicolon,
	pixelgl.KeyEqual,
	pixelgl.KeyA,
	pixelgl.KeyB,
	pixelgl.KeyC,
	pixelgl.KeyD,
	pixelgl.KeyE,
	pixelgl.KeyF,
	pixelgl.KeyG,
	pixelgl.KeyH,
	pixelgl.KeyI,
	pixelgl.KeyJ,
	pixelgl.KeyK,
	pixelgl.KeyL,
	pixelgl.KeyM,
	pixelgl.KeyN,
	pixelgl.KeyO,
	pixelgl.KeyP,
	pixelgl.KeyQ,
	pixelgl.KeyR,
	pixelgl.KeyS,
	pixelgl.KeyT,
	pixelgl.KeyU,
	pixelgl.KeyV,
	pixelgl.KeyW,
	pixelgl.KeyX,
	pixelgl.KeyY,
	pixelgl.KeyZ,
	pixelgl.KeyLeftBracket,
	pixelgl.KeyBackslash,
	pixelgl.KeyRightBracket,
	pixelgl.KeyGraveAccent,
	pixelgl.KeyWorld1,
	pixelgl.KeyWorld2,
	pixelgl.KeyEscape,
	pixelgl.KeyEnter,
	pixelgl.KeyTab,
	pixelgl.KeyBackspace,
	pixelgl.KeyInsert,
	pixelgl.KeyDelete,
	pixelgl.KeyRight,
	pixelgl.KeyLeft,
	pixelgl.KeyDown,
	pixelgl.KeyUp,
	pixelgl.KeyPageUp,
	pixelgl.KeyPageDown,
	pixelgl.KeyHome,
	pixelgl.KeyEnd,
	pixelgl.KeyCapsLock,
	pixelgl.KeyScrollLock,
	pixelgl.KeyNumLock,
	pixelgl.KeyPrintScreen,
	pixelgl.KeyPause,
	pixelgl.KeyF1,
	pixelgl.KeyF2,
	pixelgl.KeyF3,
	pixelgl.KeyF4,
	pixelgl.KeyF5,
	pixelgl.KeyF6,
	pixelgl.KeyF7,
	pixelgl.KeyF8,
	pixelgl.KeyF9,
	pixelgl.KeyF10,
	pixelgl.KeyF11,
	pixelgl.KeyF12,
	pixelgl.KeyF13,
	pixelgl.KeyF14,
	pixelgl.KeyF15,
	pixelgl.KeyF16,
	pixelgl.KeyF17,
	pixelgl.KeyF18,
	pixelgl.KeyF19,
	pixelgl.KeyF20,
	pixelgl.KeyF21,
	pixelgl.KeyF22,
	pixelgl.KeyF23,
	pixelgl.KeyF24,
	pixelgl.KeyF25,
	pixelgl.KeyKP0,
	pixelgl.KeyKP1,
	pixelgl.KeyKP2,
	pixelgl.KeyKP3,
	pixelgl.KeyKP4,
	pixelgl.KeyKP5,
	pixelgl.KeyKP6,
	pixelgl.KeyKP7,
	pixelgl.KeyKP8,
	pixelgl.KeyKP9,
	pixelgl.KeyKPDecimal,
	pixelgl.KeyKPDivide,
	pixelgl.KeyKPMultiply,
	pixelgl.KeyKPSubtract,
	pixelgl.KeyKPAdd,
	pixelgl.KeyKPEnter,
	pixelgl.KeyKPEqual,
	pixelgl.KeyMenu,
}
