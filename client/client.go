package client

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateSpinBox(initialValue, minValue, maxValue int) (*fyne.Container, *widget.Label) {
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

func CreateSecondWindow(a fyne.App, parent fyne.Window, serverEntry *widget.Entry, codecSelect, fpsSelect *widget.Select, initialRateLabel, maxRateLabel *widget.Label) {
	secondWindow := a.NewWindow("Settings Summary")
	secondWindow.Resize(fyne.NewSize(800, 600))

	initialRateValue := initialRateLabel.Text
	maxRateValue := maxRateLabel.Text

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
		parent.Show()
	})

	secondWindow.SetContent(container.NewVBox(
		summary,
		backButton,
	))

	secondWindow.SetOnClosed(func() {
		parent.Show()
	})

	secondWindow.Show()
	parent.Hide()
}

func HandleInfoButton(serverEntry *widget.Entry, codecSelect, fpsSelect *widget.Select, initialRateLabel, maxRateLabel *widget.Label, createSecondWindow func()) {
	println("Server URL:", serverEntry.Text)
	println("Codec:", codecSelect.Selected)
	println("FPS:", fpsSelect.Selected)
	println("Initial Rate:", initialRateLabel.Text)
	println("Max Rate:", maxRateLabel.Text)
	createSecondWindow()
}
