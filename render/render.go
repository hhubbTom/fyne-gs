package render

//视频渲染组件
import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // 注册JPEG解码器
	_ "image/png"  // 注册PNG解码器
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// 加载rgb-pic目录下的所有图片
func loadImages(dirPath string) ([]image.Image, error) {
	var images []image.Image
	var imagePaths []string
	imgDirPath := "fyne-gs/rgb-pic"
	// 读取目录中的所有文件
	// 如果使用相对路径不成功，尝试使用绝对路径
	_, err := ioutil.ReadDir(imgDirPath)
	if err != nil {
		// 获取当前工作目录
		wd, _ := os.Getwd()
		log.Printf("当前工作目录: %s", wd)

		// 尝试从工作目录查找
		imgDirPath = filepath.Join(wd, "rgb-pic")

		// 如果还是找不到，尝试上级目录
		if _, err := ioutil.ReadDir(imgDirPath); err != nil {
			imgDirPath = filepath.Join(wd, "../rgb-pic")
		}
	}

	log.Printf("使用图片目录: %s", imgDirPath)

	files, err := ioutil.ReadDir(imgDirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v (路径: %s)", err, imgDirPath)
	}

	// 收集所有图片文件路径
	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				imagePaths = append(imagePaths, filepath.Join(dirPath, file.Name()))
			}
		}
	}

	// 按文件名排序
	sort.Strings(imagePaths)

	// 加载所有图片
	for _, path := range imagePaths {
		file, err := os.Open(path)
		if err != nil {
			log.Printf("无法打开图片 %s: %v", path, err)
			continue
		}

		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			log.Printf("无法解码图片 %s: %v", path, err)
			continue
		}

		images = append(images, img)
	}

	return images, nil
}

func CreateVideoWindow(a fyne.App, parent fyne.Window) {
	//要做的,,,,,需要用到FFmpeg吗?
	// 创建游戏和编解码器配置
	// 创建信令客户端
	// 设置SDP offer处理函数
	// 设置ICE candidate处理函数
	// 设置对等连接
	// 连接到服务器
	//保活

	videoWindow := a.NewWindow("Video Renderer")
	videoWindow.Resize(fyne.NewSize(800, 600))

	// 固定渲染画布的大小
	canvasWidth := 640
	canvasHeight := 480

	// 加载rgb-pic目录下的图片
	images, err := loadImages("rgb-pic")
	if err != nil {
		log.Printf("加载图片错误: %v", err)
	}

	// 当前显示的图片索引
	currentImageIndex := 0
	var currentImage image.Image

	// 如果有图片，设置第一张为当前图片
	if len(images) > 0 {
		currentImage = images[0]
	}

	renderer := canvas.NewRasterWithPixels(
		func(x, y, w, h int) color.Color {
			// 如果没有加载到图片，显示随机颜色
			if currentImage == nil {
				t := time.Now().UnixNano()
				return color.RGBA{
					R: uint8(t % 255),
					G: uint8(t % 255),
					B: uint8(t % 255),
					A: 255,
				}
			}

			// 将图片内容显示到渲染器中
			bounds := currentImage.Bounds()
			if x < bounds.Max.X && y < bounds.Max.Y {
				return currentImage.At(x, y)
			}

			// 超出图片范围的部分显示为黑色
			return color.RGBA{0, 0, 0, 255}
		},
	)

	// 设置渲染器大小
	renderer.Resize(fyne.NewSize(float32(canvasWidth), float32(canvasHeight)))
	// 包裹渲染器并设置固定大小
	rendererContainer := container.NewWithoutLayout(renderer)
	rendererContainer.Resize(fyne.NewSize(float32(canvasWidth), float32(canvasHeight)))

	// 使用 VBox 和 Spacer 将渲染器左上显示,暂未居中
	centeredContainer := container.NewVBox(
		container.NewHBox(
			rendererContainer,  // 渲染器容器
			layout.NewSpacer(), // 右侧 Spacer
		),
		layout.NewSpacer(), // 下方 Spacer
	)

	// 显示图片数量和当前帧率
	statusInfo := "未找到图片"
	if len(images) > 0 {
		statusInfo = fmt.Sprintf("已加载 %d 张图片, 渲染速率: 30 FPS", len(images))
	}
	statusLabel := widget.NewLabel(statusInfo)

	stopButton := widget.NewButton("Stop", func() {
		videoWindow.Close()
	})

	// 30FPS的渲染循环
	go func() {
		ticker := time.NewTicker(time.Second / 30)
		defer ticker.Stop()

		for range ticker.C {
			if len(images) > 0 {
				// 更新当前图片
				currentImage = images[currentImageIndex]
				currentImageIndex = (currentImageIndex + 1) % len(images)
			}
			renderer.Refresh()
		}
	}()

	videoWindow.SetContent(container.NewBorder(
		nil,
		container.NewHBox(statusLabel, layout.NewSpacer(), stopButton),
		nil,
		nil,
		centeredContainer,
	))

	videoWindow.SetOnClosed(func() {
		parent.Show()
	})

	videoWindow.Show()
	parent.Hide()
}
