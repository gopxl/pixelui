package pixelui

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type widget interface {
	Draw(t pixel.Target, pos pixel.Vec)
	Bounds() pixel.Rect
}

// TEXT WIDGETS

type widgetText struct {
	text string
}

func (w *widgetText) Draw(t pixel.Target, offset pixel.Vec) {
	ui.font.Clear()
	ui.font.Color = CurrentStyle().Text.Color
	fmt.Fprint(ui.font, w.text)
	ui.font.Draw(t, pixel.IM.Moved(offset))
}

func (w *widgetText) Bounds() pixel.Rect {
	r := ui.font.BoundsOf(w.text)
	return r.Moved(r.Min.Scaled(-1))
}

func Text(text string) {
	withCurrentWindow(func(w *window) {
		w.push(&widgetText{
			text: text,
		})
	})
}

func Textf(f string, a ...interface{}) {
	Text(fmt.Sprintf(f, a...))
}

// BUTTON WIDGETS

type buttonState uint8

const (
	btnStateHover buttonState = iota + 1
	btnStatePressed
	btnStateReleased
)

var (
	btnCanvas *pixelgl.Canvas
)

type widgetButton struct {
	bounds  pixel.Rect
	state   buttonState
	content widget
}

func (w *widgetButton) Bounds() pixel.Rect {
	return w.bounds
}

func (w *widgetButton) Draw(t pixel.Target, offset pixel.Vec) {
	if !rectEq(btnCanvas.Bounds(), w.bounds) {
		btnCanvas.SetBounds(w.bounds)
	}
	switch w.state {
	case btnStateHover:
		btnCanvas.Clear(CurrentStyle().Button.Hover)
	case btnStatePressed:
		btnCanvas.Clear(CurrentStyle().Button.Pressed)
	default:
		btnCanvas.Clear(CurrentStyle().Button.Background)
	}
	w.content.Draw(btnCanvas, pixel.ZV)

	btnCanvas.Draw(t, pixel.IM.Moved(w.bounds.Center()).Moved(offset))
}

func button(id string, content widget) bool {
	if btnCanvas == nil {
		btnCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 1, 1))
	}

	btn := &widgetButton{
		content: content,
		bounds:  content.Bounds(),
	}

	withCurrentWindow(func(w *window) {
		offset := w.push(btn)
		if mouseIn(btn.bounds.Moved(offset)) {
			if consumeHeld(pixelgl.MouseButtonLeft) {
				btn.state = btnStatePressed
			} else if consumeReleased(pixelgl.MouseButtonLeft) {
				btn.state = btnStateReleased
			} else {
				btn.state = btnStateHover
			}
		}
	})

	return btn.state == btnStateReleased
}

func Button(id string) bool {
	return button(id, &widgetText{
		text: id,
	})
}

// IMAGE WIDGETS

type widgetImage struct {
	id     int
	scale  float64
	sprite *pixel.Sprite
}

func (w *widgetImage) Bounds() pixel.Rect {
	return rectRecenter(w.sprite.Frame())
}

func (w *widgetImage) Draw(t pixel.Target, offset pixel.Vec) {
	w.sprite.Draw(t, pixel.IM.Scaled(w.sprite.Frame().Min, w.scale).Moved(offset))
}

func Image(id int, scale float64) {
	if ui.packer == nil {
		panic("Cannot use the Image object without providing a sprite Packer to pixelui")
	}

	img := &widgetImage{
		id:     id,
		scale:  scale,
		sprite: pixel.NewSprite(ui.imgPic, ui.packer.BoundsOf(id)),
	}

	withCurrentWindow(func(w *window) {
		w.push(img)
	})
}
