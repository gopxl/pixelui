package pixelui

import (
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

func (ui *UI) AddSprite(sprite *pixel.Sprite) imgui.TextureID {
	id, _ := ui.packer.InsertV(sprite, packer.InsertFlipped)
	return imgui.TextureID(id)
}
