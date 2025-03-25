package render

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func CreateVideoWindow(a fyne.App, parent fyne.Window) {
	videoWindow := a.NewWindow("Video Renderer")
	videoWindow.Resize(fyne.NewSize(800, 600))

	// 固定渲染画布的大小
	canvasWidth := 640
	canvasHeight := 480

	renderer := canvas.NewRasterWithPixels(
		func(_, _, w, h int) color.Color {
			return color.RGBA{
				R: uint8(time.Now().UnixNano() % 255),
				G: uint8(time.Now().UnixNano() % 255),
				B: uint8(time.Now().UnixNano() % 255),
				A: 255,
			}
		},
	)
	// 包裹渲染器并设置固定大小
	rendererContainer := container.NewMax(renderer)
	rendererContainer.Resize(fyne.NewSize(float32(canvasWidth), float32(canvasHeight)))

	statusLabel := widget.NewLabel("Rendering at 60 FPS")
	stopButton := widget.NewButton("Stop", func() {
		videoWindow.Close()
	})

	go func() {
		ticker := time.NewTicker(time.Second / 60)
		defer ticker.Stop()

		for range ticker.C {
			renderer.Refresh()
		}
	}()

	videoWindow.SetContent(container.NewBorder(
		nil,
		container.NewHBox(statusLabel, layout.NewSpacer(), stopButton),
		nil,
		nil,
		rendererContainer, // 使用固定大小的容器
	))

	videoWindow.SetOnClosed(func() {
		parent.Show()
	})

	videoWindow.Show()
	parent.Hide()
}
