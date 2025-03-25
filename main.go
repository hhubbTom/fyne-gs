package main

import (
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createSpinBox(initialValue, minValue, maxValue int) (*fyne.Container, *widget.Label) {
	value := initialValue
	valueLabel := widget.NewLabel(strconv.Itoa(value))
	valueEntry := widget.NewEntry()
	valueEntry.SetText(strconv.Itoa(value))

	// 增加按钮
	incrementButton := widget.NewButton("^", func() {
		if value < maxValue {
			value++
			valueEntry.SetText(strconv.Itoa(value))
			valueLabel.SetText(strconv.Itoa(value)) // 同步更新 Label
		}
	})
	incrementButton.Resize(fyne.NewSize(3, 10)) // 设置按钮大小

	// 减少按钮
	decrementButton := widget.NewButton("v", func() {
		if value > minValue {
			value--
			valueEntry.SetText(strconv.Itoa(value))
			valueLabel.SetText(strconv.Itoa(value)) // 同步更新 Label
		}
	})
	decrementButton.Resize(fyne.NewSize(3, 10)) // 设置按钮大小

	// 手动编辑逻辑
	valueEntry.OnChanged = func(text string) {
		if text == "" {
			// 删不了，用这个试试
			return
		}
		if newValue, err := strconv.Atoi(text); err == nil {
			if newValue >= minValue && newValue <= maxValue {
				value = newValue
				valueLabel.SetText(strconv.Itoa(value)) // 同步更新 Label
			} else {
				valueEntry.SetText(strconv.Itoa(value))
			}
		} else {
			valueEntry.SetText(strconv.Itoa(value))
		}
	}

	// 布局：按钮垂直排列在右侧
	spinBox := container.NewBorder(
		nil, nil, nil,
		container.NewVBox(incrementButton, decrementButton),
		valueEntry,
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
	initialRateSpinBox, initialRateLabel := createSpinBox(5, 0, 100) // 初始值为 5，范围为 1 到 100
	maxRateSpinBox, maxRateLabel := createSpinBox(30, 0, 100)        // 初始值为 30，范围为 1 到 100

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

	// 创建视频渲染窗口
	createVideoWindow := func() {
		videoWindow := a.NewWindow("Video Renderer")
		videoWindow.Resize(fyne.NewSize(640, 480))

		// 创建画布用于渲染
		renderer := canvas.NewRasterWithPixels(
			func(_, _, w, h int) color.Color {
				// 这里生成随机颜色模拟视频帧
				return color.RGBA{
					R: uint8(time.Now().UnixNano() % 255),
					G: uint8(time.Now().UnixNano() % 255),
					B: uint8(time.Now().UnixNano() % 255),
					A: 255,
				}
			},
		)

		// 控制面板
		statusLabel := widget.NewLabel("Rendering at 60 FPS")
		stopButton := widget.NewButton("Stop", func() {
			videoWindow.Close()
		})

		// 模拟60FPS渲染
		go func() {
			ticker := time.NewTicker(time.Second / 60)
			defer ticker.Stop()

			for range ticker.C {
				// 刷新画布
				renderer.Refresh()
			}
		}()

		// 布局
		videoWindow.SetContent(container.NewBorder(
			nil,
			container.NewHBox(statusLabel, layout.NewSpacer(), stopButton),
			nil,
			nil,
			renderer,
		))

		videoWindow.SetOnClosed(func() {
			w.Show()
		})

		videoWindow.Show()
		w.Hide()
	}

	// 按钮
	nextButton := widget.NewButton("Next", createVideoWindow)

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
