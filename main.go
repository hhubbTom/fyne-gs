package main

import (
	"fyne-gs/client"
	"fyne-gs/render"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Piongs Client")
	w.Resize(fyne.NewSize(800, 600)) // 设置窗口初始大小

	// 小标题
	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("Enter the server URL")

	// 下拉框
	codecSelect := widget.NewSelect([]string{"H.264", "VP8", "VP9"}, nil)
	codecSelect.SetSelected("H.264")

	fpsSelect := widget.NewSelect([]string{"30FPS", "60FPS", "120FPS"}, nil)
	fpsSelect.SetSelected("60FPS")

	// 自定义上下选择控件
	initialRateSpinBox, initialRateLabel := client.CreateSpinBox(5, 0, 100) // 初始值为 5，范围为 0 到 100
	maxRateSpinBox, maxRateLabel := client.CreateSpinBox(30, 0, 100)        // 初始值为 30，范围为 0 到 100

	// 创建第二个窗口的函数
	createSecondWindow := func() {
		client.CreateSecondWindow(a, w, serverEntry, codecSelect, fpsSelect, initialRateLabel, maxRateLabel)
	}

	// 创建视频渲染窗口
	createVideoWindow := func() {
		render.CreateVideoWindow(a, w)
	}

	// 按钮
	nextButton := widget.NewButton("Next", createVideoWindow)

	infoButton := widget.NewButton("Info", func() {
		client.HandleInfoButton(serverEntry, codecSelect, fpsSelect, initialRateLabel, maxRateLabel, createSecondWindow)
	})

	// 布局
	form := container.NewVBox(
		widget.NewLabel("Connect to server"),
		serverEntry,
		widget.NewLabel("Enter the server URL."),
	)

	settings := container.NewHBox(
		container.NewVBox(codecSelect),
		container.NewVBox(fpsSelect),
		container.NewVBox(
			container.NewPadded(initialRateSpinBox), // 添加 Padding 让控件更紧凑
		),
		container.NewVBox(
			container.NewPadded(maxRateSpinBox), // 添加 Padding 让控件更紧凑
		),
	)

	buttonContainer := container.NewHBox(
		layout.NewSpacer(),
		nextButton,
	)

	infoButtonContainer := container.NewHBox(
		infoButton,
		layout.NewSpacer(), // 使用 Spacer 将按钮推到左侧
	)

	content := container.NewVBox(
		form,
		widget.NewLabel("When connecting to local server, add http://<ip>:<port> and ws://<ip>:<port> to \"Insecure origins treated as secure\" flag in chrome://flags."),
		settings,
		container.NewHBox(
			infoButtonContainer, // 左下角的信息按钮
			layout.NewSpacer(),
			buttonContainer, // 右下角的 Next 按钮
		),
	)

	w.SetContent(content)
	w.ShowAndRun()
}
