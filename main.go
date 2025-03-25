package main

import (
	"fyne.io/fyne/layout"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Piongs Client")
	w.Resize(fyne.NewSize(800, 600)) // 设置窗口初始大小

	// 输入框
	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("Enter the server URL")

	// 下拉框
	codecSelect := widget.NewSelect([]string{"H.264", "VP8", "VP9"}, nil)
	codecSelect.SetSelected("H.264")

	fpsSelect := widget.NewSelect([]string{"30FPS", "60FPS", "120FPS"}, nil)
	fpsSelect.SetSelected("60FPS")

	// 输入框（初始速率和最大速率）
	initialRateEntry := widget.NewEntry()
	initialRateEntry.SetPlaceHolder("Initial Rate Mbps")

	maxRateEntry := widget.NewEntry()
	maxRateEntry.SetPlaceHolder("Max Rate Mbps")

	// 按钮
	nextButton := widget.NewButton("Next", func() {
		// 在这里处理按钮点击事件
		println("Server URL:", serverEntry.Text)
		println("Codec:", codecSelect.Selected)
		println("FPS:", fpsSelect.Selected)
		println("Initial Rate:", initialRateEntry.Text)
		println("Max Rate:", maxRateEntry.Text)
	})

	// 布局
	form := container.NewVBox(
		widget.NewLabel("Connect to server"),
		serverEntry,
		widget.NewLabel("Enter the server URL."),
		// codecSelect,
		// fpsSelect,
		// initialRateEntry,
		// maxRateEntry,
		// nextButton,
	)
	settings := container.NewGridWithColumns(4,
		container.NewVBox(
			codecSelect,
		),
		container.NewVBox(
			fpsSelect,
		),
		container.NewVBox(
			initialRateEntry,
			widget.NewLabel("Initial Rate Mbps."),
		),
		container.NewVBox(
			maxRateEntry,
			widget.NewLabel("Max Rate Mbps."),
		),
	)

	buttonContainer := container.NewHBox(
		layout.NewSpacer(), // 使用 Spacer 将按钮推到右侧
		nextButton,
	)
	content := container.NewVBox(
		form,
		widget.NewLabel("When connecting to local server, add http://<ip>:<port> and ws://<ip>:<port> to \"Insecure origins treated as secure\" flag in chrome://flags."),
		settings,
		buttonContainer,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
