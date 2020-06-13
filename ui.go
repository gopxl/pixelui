package pixelui

import (
	"C"
	"fmt"
	"unsafe"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/inkyblackness/imgui-go"
)
import "log"

type Ui struct {
	tris    *pixel.TrianglesData
	batch   *pixel.Batch
	context *imgui.Context
	io      imgui.IO
	fonts   imgui.FontAtlas
}

func NewUi(context *imgui.Context) *Ui {
	ui := &Ui{
		context: context,
	}

	ui.tris = pixel.MakeTrianglesData(0)
	ui.batch = pixel.NewBatch(ui.tris, nil)

	ui.io = imgui.CurrentIO()
	ui.io.SetDisplaySize(imgui.Vec2{X: 1920, Y: 1080})

	ui.fonts = ui.io.Fonts()
	ui.fonts.AddFontFromFileTTF("./test.ttf", 16)
	if err := ui.fonts.BuildWithFreeType(); err != nil {
		log.Fatal(err)
	}

	return ui
}

func (ui *Ui) Draw(win *pixelgl.Window) {
	imgui.Render()
	data := imgui.RenderedDrawData()

	vertexSize, posOffset, uvOffset, colOffset := imgui.VertexBufferLayout()
	for _, cmds := range data.CommandLists() {
		start, byteSize := cmds.VertexBuffer()
		idxStart, _ := cmds.IndexBuffer()

		ui.tris.SetLen(byteSize / vertexSize)
		for _, cmd := range cmds.Commands() {
			if cmd.HasUserCallback() {
				fmt.Println("callback")
			} else {
				jndex := 0
				for i := 0; i < cmd.ElementCount()/3; i++ {
					idx := unsafe.Pointer(uintptr(idxStart) + uintptr(i*imgui.IndexBufferLayout()))
					index := *(*C.ushort)(idx)
					ptr := unsafe.Pointer(uintptr(start) + (uintptr(int(index) * vertexSize)))
					pos := (*imgui.Vec2)(unsafe.Pointer(uintptr(ptr) + uintptr(posOffset)))
					uv := (*imgui.Vec2)(unsafe.Pointer(uintptr(ptr) + uintptr(uvOffset)))
					col := (*uint32)(unsafe.Pointer(uintptr(ptr) + uintptr(colOffset)))

					position := imguiVecToPixelVec(*pos)

					(*ui.tris)[jndex].Position = position
					(*ui.tris)[jndex].Picture = imguiVecToPixelVec(*uv)
					(*ui.tris)[jndex].Color = imguiColorToPixelColor(*col)
					(*ui.tris)[jndex].Intensity = 1.0
					jndex++
				}
			}
		}

	}
	ui.batch.Dirty()

	win.SetMatrix(pixel.IM.ScaledXY(win.Bounds().Center(), pixel.V(1, -1)))

	ui.batch.Draw(win)
	ui.tris.SetLen(0)

	win.SetMatrix(pixel.IM)
}

func imguiColorToPixelColor(c uint32) pixel.RGBA {
	return pixel.RGBA{
		R: float64((c&0xFF000000)>>24) / 256,
		G: float64((c&0x00FF0000)>>16) / 256,
		B: float64((c&0x0000FF00)>>8) / 256,
		A: float64(c&0x000000FF) / 256,
	}
}

func imguiVecToPixelVec(v imgui.Vec2) pixel.Vec {
	return pixel.V(float64(v.X), float64(v.Y))
}
