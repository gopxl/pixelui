package pixelui

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/inkyblackness/imgui-go/v4"
)

type Clipboard struct {
	win *opengl.Window
}

func (c Clipboard) Text() (text string, err error) {
	text = c.win.ClipboardText()
	return
}

func (c Clipboard) SetText(value string) {
	c.win.SetClipboardText(value)
}

func (ui *UI) initIO() {
	ui.io.SetDisplaySize(IVec(ui.win.Bounds().Size()))
	ui.io.SetClipboard(Clipboard{win: ui.win})

	for k, v := range keyMap {
		ui.io.KeyMap(v, int(k))
	}

	ui.win.SetButtonCallback(func(win *opengl.Window, button pixel.Button, action pixel.Action) {
		if !button.IsKeyboardButton() {
			return
		}

		switch action {
		case pixel.Press:
			ui.io.KeyPress(int(button))
		case pixel.Release:
			ui.io.KeyRelease(int(button))
		}
	})

	ui.io.SetBackendFlags(imgui.BackendFlagsHasMouseCursors | imgui.BackendFlagsHasSetMousePos)

	ui.cursors[imgui.MouseCursorArrow] = opengl.CreateStandardCursor(opengl.ArrowCursor)
	ui.cursors[imgui.MouseCursorTextInput] = opengl.CreateStandardCursor(opengl.IBeamCursor)
	ui.cursors[imgui.MouseCursorHand] = opengl.CreateStandardCursor(opengl.HandCursor)
	ui.cursors[imgui.MouseCursorResizeEW] = opengl.CreateStandardCursor(opengl.HResizeCursor)
	ui.cursors[imgui.MouseCursorResizeNS] = opengl.CreateStandardCursor(opengl.VResizeCursor)
}

// prepareIO tells imgui.io about our current io state.
func (ui *UI) prepareIO() {
	ui.io.SetDisplaySize(IVec(ui.win.Bounds().Size()))

	ui.io.AddMouseWheelDelta(float32(ui.win.MouseScroll().X), float32(ui.win.MouseScroll().Y))
	mouse := ui.matrix.Unproject(ui.win.MousePosition())
	ui.io.SetMousePosition(imgui.Vec2{X: float32(mouse.X), Y: float32(mouse.Y)})

	ui.io.SetMouseButtonDown(0, ui.win.Pressed(pixel.MouseButtonLeft))
	ui.io.SetMouseButtonDown(1, ui.win.Pressed(pixel.MouseButtonRight))
	ui.io.SetMouseButtonDown(2, ui.win.Pressed(pixel.MouseButtonMiddle))

	ui.io.AddInputCharacters(ui.win.Typed())
	ui.updateKeyMod()

	c, has := ui.cursors[imgui.MouseCursor()]
	if !has {
		c = ui.cursors[imgui.MouseCursorArrow]
	}
	ui.win.SetCursor(c)
}

// updateKeyMod tells imgui.io where to find our key modifiers
func (ui *UI) updateKeyMod() {
	ui.io.KeyCtrl(int(pixel.KeyLeftControl), int(pixel.KeyRightControl))
	ui.io.KeyShift(int(pixel.KeyLeftShift), int(pixel.KeyRightShift))
	ui.io.KeyAlt(int(pixel.KeyLeftAlt), int(pixel.KeyRightAlt))
	ui.io.KeySuper(int(pixel.KeyLeftSuper), int(pixel.KeyRightSuper))
}

// inputWant is a helper for determining what type a button is: keyboard/mouse
func (ui *UI) inputWant(button pixel.Button) bool {
	switch button {
	case pixel.MouseButton1, pixel.MouseButton2, pixel.MouseButton3, pixel.MouseButton4, pixel.MouseButton5, pixel.MouseButton6, pixel.MouseButton7, pixel.MouseButton8:
		return ui.io.WantCaptureMouse()
	}
	return ui.io.WantCaptureKeyboard()
}

// MouseScroll returns the mouse scroll amount if imgui does not want the mouse
//
//	(if mouse is not hovering an imgui element)
func (ui *UI) MouseScroll() pixel.Vec {
	if ui.io.WantCaptureMouse() {
		return pixel.ZV
	}

	return ui.win.MouseScroll()
}

// JustPressed returns true if imgui hasn't handled the button and the button was just pressed
func (ui *UI) JustPressed(button pixel.Button) bool {
	return !ui.inputWant(button) && ui.win.JustPressed(button)
}

// JustPressed returns true if imgui hasn't handled the button and the button was just released
func (ui *UI) JustReleased(button pixel.Button) bool {
	return !ui.inputWant(button) && ui.win.JustReleased(button)
}

// JustPressed returns true if imgui hasn't handled the button and the button is pressed
func (ui *UI) Pressed(button pixel.Button) bool {
	return !ui.inputWant(button) && ui.win.Pressed(button)
}

// Repeated returns true if imgui hasn't handled the button and the button was repeated
func (ui *UI) Repeated(button pixel.Button) bool {
	return !ui.inputWant(button) && ui.win.Repeated(button)
}

// KeyCtrl returns true if either left or right control is pressed
func (ui *UI) KeyCtrl() bool {
	return ui.win.Pressed(pixel.KeyLeftControl) || ui.win.Pressed(pixel.KeyRightControl)
}

// KeyCtrl returns true if either left or right shift is pressed
func (ui *UI) KeyShift() bool {
	return ui.win.Pressed(pixel.KeyLeftShift) || ui.win.Pressed(pixel.KeyRightShift)
}

// KeyCtrl returns true if either left or right alt is pressed
func (ui *UI) KeyAlt() bool {
	return ui.win.Pressed(pixel.KeyLeftAlt) || ui.win.Pressed(pixel.KeyRightAlt)
}

// KeyCtrl returns true if either left or right super (windows key) is pressed
func (ui *UI) KeySuper() bool {
	return ui.win.Pressed(pixel.KeyLeftSuper) || ui.win.Pressed(pixel.KeyRightSuper)
}

var (
	keyMap = map[pixel.Button]int{
		pixel.KeyTab:       imgui.KeyTab,
		pixel.KeyLeft:      imgui.KeyLeftArrow,
		pixel.KeyRight:     imgui.KeyRightArrow,
		pixel.KeyUp:        imgui.KeyUpArrow,
		pixel.KeyDown:      imgui.KeyDownArrow,
		pixel.KeyPageUp:    imgui.KeyPageUp,
		pixel.KeyPageDown:  imgui.KeyPageDown,
		pixel.KeyHome:      imgui.KeyHome,
		pixel.KeyEnd:       imgui.KeyEnd,
		pixel.KeyInsert:    imgui.KeyInsert,
		pixel.KeyDelete:    imgui.KeyDelete,
		pixel.KeyBackspace: imgui.KeyBackspace,
		pixel.KeySpace:     imgui.KeySpace,
		pixel.KeyEnter:     imgui.KeyEnter,
		pixel.KeyEscape:    imgui.KeyEscape,
		pixel.KeyA:         imgui.KeyA,
		pixel.KeyC:         imgui.KeyC,
		pixel.KeyV:         imgui.KeyV,
		pixel.KeyX:         imgui.KeyX,
		pixel.KeyY:         imgui.KeyY,
		pixel.KeyZ:         imgui.KeyZ,
	}
)
