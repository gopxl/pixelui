package pixelui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/inkyblackness/imgui-go"
)

type Clipboard struct {
	win *pixelgl.Window
}

func (c Clipboard) Text() (text string, err error) {
	text = c.win.ClipboardText()
	return
}

func (c Clipboard) SetText(value string) {
	c.win.SetClipboardText(value)
}

// prepareIO tells imgui.io about our current io state.
func (ui *UI) prepareIO() {
	ui.io.SetDisplaySize(IVec(ui.win.Bounds().Size()))

	ui.io.AddMouseWheelDelta(float32(ui.win.MouseScroll().X), float32(ui.win.MouseScroll().Y))
	mouse := ui.matrix.Unproject(ui.win.MousePosition())
	ui.io.SetMousePosition(imgui.Vec2{X: float32(mouse.X), Y: float32(mouse.Y)})

	ui.io.SetMouseButtonDown(0, ui.win.Pressed(pixelgl.MouseButtonLeft))
	ui.io.SetMouseButtonDown(1, ui.win.Pressed(pixelgl.MouseButtonRight))
	ui.io.SetMouseButtonDown(2, ui.win.Pressed(pixelgl.MouseButtonMiddle))

	for _, key := range keys {
		if ui.win.Pressed(key) {
			ui.io.KeyPress(int(key))
		} else {
			ui.io.KeyRelease(int(key))
		}
		ui.updateKeyMod()
	}
	ui.io.AddInputCharacters(ui.win.Typed())
}

// updateKeyMod tells imgui.io where to find our key modifiers
func (ui *UI) updateKeyMod() {
	ui.io.KeyCtrl(int(pixelgl.KeyLeftControl), int(pixelgl.KeyRightControl))
	ui.io.KeyShift(int(pixelgl.KeyLeftShift), int(pixelgl.KeyRightShift))
	ui.io.KeyAlt(int(pixelgl.KeyLeftAlt), int(pixelgl.KeyRightAlt))
	ui.io.KeySuper(int(pixelgl.KeyLeftSuper), int(pixelgl.KeyRightSuper))
}

// inputWant is a helper for determining what type a button is: keyboard/mouse
func (ui *UI) inputWant(button pixelgl.Button) bool {
	switch button {
	case pixelgl.MouseButton1, pixelgl.MouseButton2, pixelgl.MouseButton3, pixelgl.MouseButton4, pixelgl.MouseButton5, pixelgl.MouseButton6, pixelgl.MouseButton7, pixelgl.MouseButton8:
		return ui.io.WantCaptureMouse()
	}
	return ui.io.WantCaptureKeyboard()
}

// MouseScroll returns the mouse scroll amount if imgui does not want the mouse
//	(if mouse is not hovering an imgui element)
func (ui *UI) MouseScroll() pixel.Vec {
	if ui.io.WantCaptureMouse() {
		return pixel.ZV
	}

	return ui.win.MouseScroll()
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

// KeyCtrl returns true if either left or right control is pressed
func (ui *UI) KeyCtrl() bool {
	return ui.win.Pressed(pixelgl.KeyLeftControl) || ui.win.Pressed(pixelgl.KeyRightControl)
}

// KeyCtrl returns true if either left or right shift is pressed
func (ui *UI) KeyShift() bool {
	return ui.win.Pressed(pixelgl.KeyLeftShift) || ui.win.Pressed(pixelgl.KeyRightShift)
}

// KeyCtrl returns true if either left or right alt is pressed
func (ui *UI) KeyAlt() bool {
	return ui.win.Pressed(pixelgl.KeyLeftAlt) || ui.win.Pressed(pixelgl.KeyRightAlt)
}

// KeyCtrl returns true if either left or right super (windows key) is pressed
func (ui *UI) KeySuper() bool {
	return ui.win.Pressed(pixelgl.KeyLeftSuper) || ui.win.Pressed(pixelgl.KeyRightSuper)
}

// setKeyMapping maps pixelgl buttons to imgui keys.
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
	pixelgl.KeyLeftShift,
	pixelgl.KeyLeftAlt,
	pixelgl.KeyLeftAlt,
	pixelgl.KeyLeftSuper,
	pixelgl.KeyRightShift,
	pixelgl.KeyRightAlt,
	pixelgl.KeyRightAlt,
	pixelgl.KeyRightSuper,
	pixelgl.KeyMenu,
	pixelgl.KeyLeftControl,
	pixelgl.KeyRightControl,
}
