package pixelui

import "github.com/inkyblackness/imgui-go"

// Image is a helper for imgui.Image that looks up the sprite in the internal packed atlas.
func (ui *UI) Image(alias interface{}, scale float64) {
	id := ui.packer.IdOf(alias)
	sprite := ui.packer.BoundsOf(alias)
	imgui.Image(imgui.TextureID(id), IVec(sprite.Size().Scaled(scale)))
}

func (ui *UI) ImageButton(alias interface{}, scale float64) bool {
	id := ui.packer.IdOf(alias)
	sprite := ui.packer.BoundsOf(alias)

	return imgui.ImageButton(imgui.TextureID(id), IVec(sprite.Size().Scaled(scale)))
}
