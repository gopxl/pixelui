package pixelui

import (
	"image/color"
	"log"
	"unsafe"

	"github.com/inkyblackness/imgui-go"

	"github.com/faiface/pixel"
)

// loadFont parses the imgui font data and creates a pixel picture from it.
func (ui *UI) loadFont() {
	f := ui.fonts.TextureDataAlpha8()
	pic := pixel.MakePictureData(pixel.R(0, 0, float64(f.Width), float64(f.Height)))

	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			i := y*f.Width + x
			ptr := (*uint8)(unsafe.Pointer(uintptr(f.Pixels) + uintptr(i)))
			pic.Pix[i] = color.RGBA{R: 0, G: 0, B: 0, A: *ptr}
		}
	}

	ui.fontAtlas = ui.win.MakePicture(pic)
	id := "default-font"
	if err := ui.packer.Replace(id, pixel.NewSprite(pic, pic.Bounds())); err != nil {
		log.Fatalln(err)
	}
	ui.fontId = ui.packer.IdOf(id)
	ui.fonts.SetTextureID(imgui.TextureID(ui.fontId))
}

// loadDefaultFont loads the imgui default font if the user wants it.
func (ui *UI) loadDefaultFont() {
	ui.fonts.AddFontDefault()
	ui.loadFont()
}

// AddTTFFont loads the given font into imgui.
func (ui *UI) AddTTFFont(path string, size float32) {
	ui.fonts.AddFontFromFileTTF(path, size)
	ui.loadFont()
}
