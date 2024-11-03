package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"encoding/base64"
)

var (
	userinput string
	content   []string
)

func init_UI() {
	sj_app := app.New()
	window := sj_app.NewWindow("SerializeJava")

	//初始化contentList
	contentList := widget.NewList(
		func() int {
			return len(content) // 返回内容数量
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Resize(fyne.NewSize(400, 300))

			label.MinSize()
			return label
		},
		func(i int, o fyne.CanvasObject) {
			label := o.(*widget.Label)                   // 将 CanvasObject 转换为 Label
			label.SetText(content[i])                    // 设置文本
			label.TextStyle = fyne.TextStyle{Bold: true} // 设置文本样式
			label.Resize(fyne.NewSize(300, 300))         // 设置文本大小
		},
	)

	//初始化entry
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter your content")

	//初始化submitButton
	submitButton := widget.NewButton("Submit", func() {
		userinput = entry.Text
		if userinput != "" {
			domain()
			contentList.Refresh()
		}
	})

	//设置尺寸和位置
	entry.Resize(fyne.NewSize(1400, 50))
	submitButton.Move(fyne.NewPos(0, 50))
	submitButton.Resize(fyne.NewSize(1400, 25))
	contentList.Move(fyne.NewPos(0, 75))
	contentList.Resize(fyne.NewSize(1400, 800))

	window.SetContent(container.NewWithoutLayout(entry, submitButton, contentList))
	window.Resize(fyne.NewSize(900, 900))
	window.ShowAndRun()
}

func domain() {

	decodedBytes, err := base64.StdEncoding.DecodeString(userinput)
	if err != nil {
		content = append(content, "请输入正确的base64编码内容！")
	}
	content = append(content, string(decodedBytes))

}

func main() {
	init_UI() // 初始化 UI
}
