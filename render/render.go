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
			t := time.Now().UnixNano()
			return color.RGBA{
				R: uint8(t % 255),
				G: uint8(t % 255),
				B: uint8(t % 255),
				A: 255,
			}
		},
	)

	// 设置渲染器大小
	renderer.Resize(fyne.NewSize(float32(canvasWidth), float32(canvasHeight)))

	// 包裹渲染器并设置固定大小
	rendererContainer := container.NewWithoutLayout(renderer)
	rendererContainer.Resize(fyne.NewSize(float32(canvasWidth), float32(canvasHeight)))

	// 使用 VBox 和 Spacer 将渲染器居中显示
	centeredContainer := container.NewVBox(

		container.NewHBox(
			layout.NewSpacer(), // 左侧 Spacer
			rendererContainer,  // 渲染器容器
			layout.NewSpacer(), // 右侧 Spacer
		),
		layout.NewSpacer(), // 下方 Spacer
	)

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
		centeredContainer, // 使用 VBox 居中容器
	))

	videoWindow.SetOnClosed(func() {
		parent.Show()
	})

	videoWindow.Show()
	parent.Hide()
}
