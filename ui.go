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
type uiDef struct {
	win      *pixelgl.Window
	canvas   *pixelgl.Canvas
	packer   *packer.Packer
	imgPic   pixel.Picture
	imgBatch *pixel.Batch

	font *text.Text

	newframe   bool
	currentWin *window
	winStack   []*window
	styleStack []Style
}

var ui uiDef

// pixelui.NewUI flags:
//	NO_DEFAULT_FONT: Do not load the default font during NewUI.
const (
	NO_DEFAULT_FONT uint8 = 1 << iota
)

// NewUI Creates the UI and setups up its internal structures
func Init(win *pixelgl.Window, flags uint8) {
	ui = uiDef{
		win:        win,
		canvas:     pixelgl.NewCanvas(win.Bounds()),
		winStack:   make([]*window, 0),
		font:       text.New(pixel.ZV, text.NewAtlas(basicfont.Face7x13, text.ASCII)),
		styleStack: []Style{defaultStyle},
	}
}

func AddImagePacker(pack *packer.Packer) {
	ui.packer = pack
	ui.imgPic = pack.Picture()
	ui.imgBatch = pixel.NewBatch(&pixel.TrianglesData{}, ui.imgPic)
}

// NewFrame Call this at the beginning of the frame to tell the UI that the frame has started
func NewFrame() {
	ui.newframe = true

	for i := 0; i < int(pixelgl.KeyLast); i++ {
		b := pixelgl.Button(i)
		held[b] = ui.win.Pressed(b)
		pressed[b] = ui.win.JustPressed(b)
		released[b] = ui.win.JustReleased(b)
	}
}

// Draw Draws the imgui UI to the Pixel Window
func Draw(win *pixelgl.Window, m pixel.Matrix) {
	if !ui.newframe {
		panic("NewFrame must be called at the beginning of every frame")
	}

	ui.canvas.Clear(pixel.Alpha(0))

	for _, w := range ui.winStack {
		w.draw(ui.canvas)
	}

	ui.canvas.Draw(win, m.Moved(win.Bounds().Center()))
	ui.newframe = false
}
