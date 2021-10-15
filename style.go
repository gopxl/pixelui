package pixelui

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type windowStyle struct {
	Background color.Color
}

type buttonStyle struct {
	Background color.Color
	Hover      color.Color
	Pressed    color.Color
}

type textStyle struct {
	Color color.Color
}

type Style struct {
	Window windowStyle
	Button buttonStyle
	Text   textStyle
}

var (
	defaultStyle = Style{
		Window: windowStyle{
			Background: color.White,
		},
		Button: buttonStyle{
			Background: colornames.Salmon,
			Hover:      colornames.Lightpink,
			Pressed:    colornames.Red,
		},
		Text: textStyle{
			Color: color.Black,
		},
	}
)

func DefaultStyle() Style {
	return defaultStyle
}

func CurrentStyle() Style {
	return ui.styleStack[len(ui.styleStack)-1]
}

func PushStyle(style Style) {
	ui.styleStack = append(ui.styleStack, style)
}

func PopStyle() {
	if len(ui.styleStack) == 1 {
		panic("About to Pop the default style")
	}

	ui.styleStack = ui.styleStack[0 : len(ui.styleStack)-1]
}
