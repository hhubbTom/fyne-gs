package main

import (
	"fmt"
	"fyne-gs/client"
	"fyne-gs/render"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// 从FPS字符串中提取数字值
func getFPSValue(fpsStr string) int { //暂未使用
	switch fpsStr {
	case "30FPS":
		return 30
	case "60FPS":
		return 60
	case "90FPS":
		return 90
	case "120FPS":
		return 120
	default:
		return 60
	}
}

// 将字符串转换为数字
func getNumericValue(text string) int { //暂未使用
	var value int
	fmt.Sscanf(text, "%d", &value)
	return value
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
	//待替换成下面的解析url
	// connectButton := widget.NewButton("连接", func() {
	//     // 获取并处理URL
	//     serverURL := serverURLEntry.Text

	//     // 创建配置
	//     gameConfig := &config.GameConfig{
	//         // 设置游戏配置
	//     }

	//     // 获取编解码器值
	//     codecMap := map[string]string{
	//         "H.264": "h264_nvenc",
	//         "H.265": "hevc_nvenc",
	//         "AV1":   "av1_nvenc",
	//     }

	//     codecConfig := &config.CodecConfig{
	//         Codec:          codecMap[codecSelect.Selected],
	//         FrameRate:      getFPSValue(fpsSelect.Selected),  // 从 FPS 字符串中提取数字值
	//         InitialBitrate: getNumericValue(initialBitrateEntry.Text) * 1_000_000,
	//         MaxBitrate:     getNumericValue(maxBitrateEntry.Text) * 1_000_000,
	//     }

	//     // 创建信令客户端并连接
	//     client := signaling.NewSignalClient(gameConfig, codecConfig)
	//     setupWebRTCConnection(client, serverURL)
	// })

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
		//handleConnect(serverURL.Text)   一个解析url的函数,待办
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
