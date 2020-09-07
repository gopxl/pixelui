package pixelui

import (
	"image"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/dusk125/pixelutils/packer"
	"github.com/inkyblackness/imgui-go"

	"github.com/faiface/pixel"
)

const (
	WrappedNone = iota
	WrappedSprite
	WrappedBatch
	WrappedCanvas
)

type wrapper struct {
	Type  int
	Value interface{}
}

func Sprite(sprite *pixel.Sprite) imgui.TextureID {
	return imgui.TextureID(unsafe.Pointer(&wrapper{
		Type:  WrappedSprite,
		Value: sprite,
	}))
}

func (ui *UI) AddSprite(name string, sprite *pixel.Sprite) imgui.TextureID {
	if err := ui.packer.InsertV(name, sprite, packer.InsertFlipped); err != nil {
		log.Fatalln(err)
	}
	return imgui.TextureID(ui.packer.IdOf(name))
}

func (ui *UI) AddSpriteFromFile(path string) (id imgui.TextureID, sprite *pixel.Sprite) {
	return ui.AddSpriteFromFileV(filepath.Base(path[:len(path)-4]), path)
}

func (ui *UI) AddSpriteFromFileV(name, path string) (id imgui.TextureID, sprite *pixel.Sprite) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}

	data := pixel.PictureDataFromImage(img)
	sprite = pixel.NewSprite(data, data.Bounds())
	id = ui.AddSprite(name, sprite)

	return
}
