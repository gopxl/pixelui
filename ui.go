package pixelui

import (
	"C"
)
import (
	"github.com/dusk125/pixelutils/packer"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

// UI Stores the state of the pixelui UI
type UI struct {
	win    *pixelgl.Window
	canvas *pixelgl.Canvas
	packer *packer.Packer

	font *text.Text

	currentWin *window
	winStack   []*window
	styleStack []Style
}

var ui *UI

// pixelui.NewUI flags:
//	NO_DEFAULT_FONT: Do not load the default font during NewUI.
const (
	NO_DEFAULT_FONT uint8 = 1 << iota
)

// NewUI Creates the UI and setups up its internal structures
func Init(win *pixelgl.Window, flags uint8) {
	ui = &UI{
		win:        win,
		canvas:     pixelgl.NewCanvas(win.Bounds()),
		packer:     packer.New(),
		winStack:   make([]*window, 0),
		font:       text.New(pixel.ZV, text.NewAtlas(basicfont.Face7x13, text.ASCII)),
		styleStack: []Style{defaultStyle},
	}

	if flags&NO_DEFAULT_FONT == 0 {
	}
}

func GetPacker() *packer.Packer {
	return ui.packer
}

// NewFrame Call this at the beginning of the frame to tell the UI that the frame has started
func NewFrame() {
}

// Draw Draws the imgui UI to the Pixel Window
func Draw(win *pixelgl.Window, m pixel.Matrix) {
	ui.canvas.Clear(pixel.Alpha(0))

	for _, w := range ui.winStack {
		w.draw(ui.canvas)
	}

	ui.canvas.Draw(win, m.Moved(win.Bounds().Center()))
}
