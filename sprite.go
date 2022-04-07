package pixelui

import (
	"image"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"unsafe"

	"github.com/inkyblackness/imgui-go/v4"

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

func (ui *UI) nextID() int {
	return int(atomic.AddInt32(&ui.lastID, 1))
}

func (ui *UI) AddSprite(name string, sprite *pixel.Sprite) imgui.TextureID {
	pic := sprite.Picture().(*pixel.PictureData)
	frame := sprite.Frame()
	newPic := pixel.MakePictureData(pixel.R(0, 0, frame.W(), frame.H()))
	i := 0
	for y := frame.Min.Y; y < frame.Max.Y; y++ {
		for x := frame.Min.X; x < frame.Max.X; x++ {
			newPic.Pix[i] = pic.Pix[pic.Index(pixel.V(float64(x), float64(y)))]
			i++
		}
	}
	id := ui.nextID()
	ui.packer.Insert(id, newPic.Image())
	return imgui.TextureID(id)
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
