package pixelui

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type marginStyle struct {
	top, bottom, left, right float64
}

type windowStyle struct {
	Background color.Color
}

type buttonStyle struct {
	Background color.Color
	Hover      color.Color
	Pressed    color.Color
	margin     marginStyle
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
			Hover:      colornames.Aquamarine,
			Pressed:    colornames.Red,
			margin: marginStyle{
				top:    2,
				bottom: 2,
				left:   2,
				right:  2,
			},
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
