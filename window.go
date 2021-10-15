package pixelui

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type WindowFlags uint

type window struct {
	id      string
	open    *bool
	flags   WindowFlags
	bounds  pixel.Rect
	canvas  *pixelgl.Canvas
	widgets []widget
}

func (w *window) push(wid widget) {
	w.widgets = append(w.widgets, wid)
}

func (w *window) draw(t pixel.Target) {
	w.canvas.Clear(pixel.Alpha(1))

	ui.font.Clear()
	ui.font.Color = colornames.Black
	fmt.Fprint(ui.font, w.id)
	ui.font.Draw(w.canvas, pixel.IM.Moved(pixel.V(w.bounds.W()/2-ui.font.Bounds().W()/2, w.bounds.H()-ui.font.LineHeight)))

	for i, wid := range w.widgets {
		wid.Draw(w.canvas, pixel.V(0, w.bounds.H()-(float64(i+2)*ui.font.LineHeight)))
	}

	w.widgets = w.widgets[:0]
	w.canvas.Draw(t, pixel.IM.Moved(w.bounds.Center()))
}

func findWindow(id string) *window {
	for _, win := range ui.winStack {
		if win.id == id {
			return win
		}
	}
	return nil
}

func withCurrentWindow(f func(w *window)) {
	var newWin = ui.currentWin == nil

	if newWin {
		Begin("")
	}

	f(ui.currentWin)

	if newWin {
		End()
	}
}

func BeginV(id string, open *bool, flags WindowFlags) bool {
	if ui.currentWin = findWindow(id); ui.currentWin == nil {
		ui.currentWin = &window{
			id:      id,
			open:    open,
			flags:   flags,
			bounds:  rect(0, 0, 100, 100),
			widgets: make([]widget, 0),
		}
		ui.currentWin.canvas = pixelgl.NewCanvas(ui.currentWin.bounds)

		ui.winStack = append(ui.winStack, ui.currentWin)
	} else {
		ui.currentWin.flags = flags
	}

	if ui.currentWin.open != nil {
		return *ui.currentWin.open
	}

	return true
}

func Begin(id string) bool {
	return BeginV(id, nil, 0)
}

func End() {
	if ui.currentWin == nil {
		panic("No current window, did you call 'End' too many times?")
	}

	if ui.currentWin.bounds.Contains(ui.win.MousePosition()) && ui.win.Pressed(pixelgl.MouseButtonLeft) {
		ui.currentWin.bounds = ui.currentWin.bounds.Moved(ui.win.MousePosition().Sub(ui.win.MousePreviousPosition()))
	}

	ui.currentWin = nil
}
