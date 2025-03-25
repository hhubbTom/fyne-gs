package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createSpinBox(initialValue, minValue, maxValue int) (*fyne.Container, *widget.Label) {
	value := initialValue
	valueLabel := widget.NewLabel(strconv.Itoa(value))

	// 增加按钮
	incrementButton := widget.NewButton("+", func() {
		if value < maxValue {
			value++
			valueLabel.SetText(strconv.Itoa(value))
		}
	})

	// 减少按钮
	decrementButton := widget.NewButton("-", func() {
		if value > minValue {
			value--
			valueLabel.SetText(strconv.Itoa(value))
		}
	})

	// 布局
	spinBox := container.NewVBox(
		valueLabel,
		container.NewHBox(decrementButton, incrementButton),
	)

	return spinBox, valueLabel
}

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
	initialRateSpinBox, initialRateLabel := createSpinBox(5, 1, 100) // 初始值为 5，范围为 1 到 100
	maxRateSpinBox, maxRateLabel := createSpinBox(30, 1, 100)        // 初始值为 30，范围为 1 到 100

	// 创建第二个窗口的函数
	createSecondWindow := func() {
		secondWindow := a.NewWindow("Settings Summary")
		secondWindow.Resize(fyne.NewSize(800, 600))

		// 获取 SpinBox 的当前值
		initialRateValue := initialRateLabel.Text
		maxRateValue := maxRateLabel.Text

		// 创建显示内容
		summary := widget.NewLabel(
			"Settings Summary:\n\n" +
				"Server URL: " + serverEntry.Text + "\n" +
				"Codec: " + codecSelect.Selected + "\n" +
				"FPS: " + fpsSelect.Selected + "\n" +
				"Initial Rate: " + initialRateValue + " Mbps\n" +
				"Max Rate: " + maxRateValue + " Mbps",
		)

		// 返回按钮
		backButton := widget.NewButton("Back", func() {
			secondWindow.Close()
			w.Show()
		})

		secondWindow.SetContent(container.NewVBox(
			summary,
			backButton,
		))

		// 处理窗口关闭
		secondWindow.SetOnClosed(func() {
			w.Show()
		})

		secondWindow.Show()
		w.Hide()
	}

	// 按钮
	nextButton := widget.NewButton("Next", func() {

		println("Next button clicked (no functionality).")
	})

	// 按钮
	infoButton := widget.NewButton("Info", func() {
		// 获取输入框的值并验证
		// 获取 SpinBox 的当前值
		initialRateValue := initialRateLabel.Text
		maxRateValue := maxRateLabel.Text

		println("Server URL:", serverEntry.Text)
		println("Codec:", codecSelect.Selected)
		println("FPS:", fpsSelect.Selected)
		println("Initial Rate:", initialRateValue)
		println("Max Rate:", maxRateValue)
		createSecondWindow()

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
		container.NewVBox(initialRateSpinBox),
		container.NewVBox(maxRateSpinBox),
	)

	buttonContainer := container.NewHBox(
		layout.NewSpacer(),
		nextButton,
	)

	// 信息按钮容器（左下角）
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
