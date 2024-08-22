package pixelui

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"unsafe"

	"github.com/inkyblackness/imgui-go/v4"
)

// loadFont parses the imgui font data and creates a pixel picture from it.
func (ui *UI) loadFont() {
	f := ui.fonts.TextureDataAlpha8()
	pic := image.NewRGBA(image.Rect(0, 0, f.Width, f.Height))

	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			i := y*f.Width + x
			ptr := (*uint8)(unsafe.Pointer(uintptr(f.Pixels) + uintptr(i)))
			pic.SetRGBA(x, y, color.RGBA{0, 0, 0, *ptr})
		}
	}

	ui.atlas.Clear(ui.group)
	ui.font = ui.group.AddImage(pic)
	ui.atlas.Pack()
	ui.fonts.SetTextureID(imgui.TextureID(ui.font.ID()))
}

// loadDefaultFont loads the imgui default font if the user wants it.
func (ui *UI) loadDefaultFont() {
	ui.fonts.AddFontDefault()
	ui.loadFont()
}

// AddTTFFont loads the given font into imgui.
func (ui *UI) AddTTFFont(path string, size float32) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("The font file: %s does not exist", path))
	}
	ui.fonts.AddFontFromFileTTF(path, size)
	ui.loadFont()
}
