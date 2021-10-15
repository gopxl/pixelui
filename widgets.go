package pixelui

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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

func (w *widgetText) Draw(t pixel.Target, pos pixel.Vec) {
	ui.font.Clear()
	ui.font.Color = CurrentStyle().Text.Color
	fmt.Fprint(ui.font, w.text)
	ui.font.Draw(t, pixel.IM.Moved(pos))
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
	btnStateHover    buttonState = iota
	btnStatePressed              = iota
	btnStateReleased             = iota
)

type widgetButton struct {
	bounds  pixel.Rect
	state   buttonState
	content widget
}

func (w *widgetButton) Bounds() pixel.Rect {
	return w.bounds
}

func (w *widgetButton) Draw(t pixel.Target, pos pixel.Vec) {
	im := imdraw.New(nil)
	im.Push(w.bounds.Min.Add(pos), w.bounds.Max.Add(pos))
	log.Println(w.bounds.Moved(pos))
	switch w.state {
	case btnStateHover:
		im.Color = CurrentStyle().Button.Hover
	case btnStatePressed:
		im.Color = CurrentStyle().Button.Pressed
	default:
		im.Color = CurrentStyle().Button.Background
	}
	im.Rectangle(0)
	im.Draw(t)

	w.content.Draw(t, pos)
}

func button(id string, content widget) bool {
	btn := &widgetButton{
		content: content,
		bounds:  content.Bounds(),
	}

	withCurrentWindow(func(w *window) {
		w.push(btn)
	})

	if mouseIn(btn.bounds) {
		if ui.win.Pressed(pixelgl.MouseButtonLeft) {
			btn.state = btnStatePressed
		} else if ui.win.JustReleased(pixelgl.MouseButtonLeft) {
			btn.state = btnStateReleased
		} else {
			btn.state = btnStateHover
		}
	}

	return btn.state == btnStateReleased
}

func Button(id string) bool {
	button(id, &widgetText{
		text: id,
	})
	return false
}
