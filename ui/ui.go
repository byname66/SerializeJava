package ui

import (
	"fmt"
	"image/color"
	"main/structures"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	userinput string
	result    string
)

type AppState struct {
	InputEntry    *widget.Entry
	MainContainer *fyne.Container
	MyWindow      fyne.Window
}

func (state *AppState) UpdateMainContainer(newContent *fyne.Container) {
	state.MainContainer.Objects = newContent.Objects
	state.MainContainer.Refresh()
}

func NewContainer(state *AppState) *fyne.Container {
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Enter your content/filename (base64)")
	maxCharsPerLine := 200
	inputEntry.OnChanged = func(content string) {
		lines := strings.Split(content, "\n")
		var wrappedLines []string
		for _, line := range lines {
			for len(line) > maxCharsPerLine {
				wrappedLines = append(wrappedLines, line[:maxCharsPerLine])
				line = line[maxCharsPerLine:]
			}
			wrappedLines = append(wrappedLines, line)
		}
		newContent := strings.Join(wrappedLines, "\n")
		if newContent != content {
			inputEntry.SetText(newContent)
		}
	}
	state.InputEntry = inputEntry
	scrollContainer := container.NewScroll(inputEntry)
	scrollContainer.SetMinSize(fyne.NewSize(800, 10))

	cleanButton := widget.NewButton("Clean", func() {
		inputEntry.SetText("")
		inputEntry.Refresh()
	})

	pasteButton := widget.NewButton("Paste", func() {
		clipboard := state.MyWindow.Clipboard()
		inputEntry.SetText("")
		content := clipboard.Content()
		inputEntry.SetText(content)
	})

	c1 := container.NewVBox()
	c1.Objects = NewStruContainer(state).Objects
	c2 := container.NewVBox()
	c2.Objects = NewFuncContainer(state).Objects
	tab1 := container.NewTabItem("Show STREAM Structure", c1)
	tab2 := container.NewTabItem("Modify STREAM Data", c2)
	tabs := container.NewAppTabs(tab1, tab2)

	state.MainContainer = container.NewVBox()

	return container.NewVBox(
		container.NewHBox(scrollContainer, cleanButton, pasteButton),
		tabs,
		state.MainContainer,
	)
}

func NewStruContainer(state *AppState) *fyne.Container {

	richText := widget.NewRichText()
	label := widget.NewLabel("")
	scrollContainer := container.NewScroll(richText)
	//scrollContainer := container.NewScroll(label)

	scrollContainer.SetMinSize(fyne.NewSize(600, 700))

	showButton := widget.NewButton("Show the Stream Structure", func() {
		userinput = strings.TrimSpace(state.InputEntry.Text)
		userinput = strings.ReplaceAll(userinput, "\n", "")
		if userinput == "" {
			dialog.ShowInformation("Error", "Input is empty!", state.MyWindow)
			return
		}
		domain() //func domain() use global variable "userinput" and change global variable "result".
		// label.Text = result
		// label.Refresh()
		richText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{Text: result},
		}
		richText.Refresh()
	})

	copyButton := widget.NewButton("Copy", func() {
		state.MyWindow.Clipboard().SetContent(label.Text)
		dialog.ShowInformation("Success", "Copied to clipboard!", state.MyWindow)
	})

	return container.NewVBox(
		showButton,
		scrollContainer,
		copyButton,
	)
}

func NewFuncContainer(state *AppState) *fyne.Container {
	var stream *structures.Stream
	empty := widget.NewLabel("")
	line := canvas.NewLine(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	line1 := canvas.NewLine(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	line2 := canvas.NewLine(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	line3 := canvas.NewLine(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	//func InsertDirtyData:
	funcDirtyTheme := canvas.NewText("Insert dirty data", color.RGBA{R: 0, G: 0, B: 255, A: 255}) // 蓝色字体
	funcDirtyTheme.TextStyle.Bold = true
	funcDirtyTheme.TextSize = 16
	funcDirtyIntroduce1 := canvas.NewText("You can input a number to specify how many bytes (JAVA_TC_RESET) you want to insert into the serialized stream.", color.Black)
	funcDirtyIntroduce2 := canvas.NewText("This can help bypass WAFs that impose a length limit.", color.Black)
	funcDirtyInputEntry := widget.NewEntry()
	funcDirtyCheck := widget.NewCheck("check", func(checked bool) {
		if checked {
			funcDirtyInputEntry.Enable()
		} else {
			funcDirtyInputEntry.Disable()
		}
		funcDirtyInputEntry.Refresh()
	})
	funcDirtyInputLabel := widget.NewLabel("The number of byte you want to insert")
	//func Utf OverLoad Encoding.
	funcUtfTheme := canvas.NewText("UTF OverLong Encoding", color.RGBA{R: 0, G: 0, B: 255, A: 255}) // blue
	funcUtfTheme.TextStyle.Bold = true
	funcUtfTheme.TextSize = 16
	funcUtfIntroduce1 := canvas.NewText("You can choose model 2 or 3 (Overlong encoding of 1 UTF character using 2 or 3 bytes).", color.Black)
	funcUtfIntroduce2 := canvas.NewText("This approach can bypass some WAFs that check UTF characters (such as class names or parameters).", color.Black)
	funcUtfCheck := widget.NewCheck("check", func(checked bool) {
	})
	funcUtfLabel := widget.NewLabel("The model of overlong encode UTF.")
	radio := widget.NewRadioGroup([]string{"2", "3"}, func(selected string) {
	})
	funcUtfCheck.OnChanged = func(checked bool) {
		if checked {
			radio.Enable()
		} else {
			radio.Disable()
		}
		radio.Refresh()
	}
	//func change serialVersionUID
	funcUidTheme := canvas.NewText("Change Class SerialVersionUID", color.RGBA{R: 0, G: 0, B: 255, A: 255})
	funcUidTheme.TextStyle.Bold = true
	funcUidTheme.TextSize = 16
	funcUidIntroduce1 := canvas.NewText("You can modify the SerialVersionUID of the class you want to change in the Java serialization stream.", color.Black)
	funcUidIntroduce2 := canvas.NewText("This helps to avoid errors caused by class version mismatches between the serialization (client) and deserialization (server) processes.", color.Black)
	scrollContainer := container.NewScroll(empty)
	scrollContainer.SetMinSize(fyne.NewSize(200, 200))
	funcUidCheck := widget.NewCheck("check", func(checked bool) {
	})
	funcUidCheck.OnChanged = func(checked bool) {
		if checked {
			userinput = strings.TrimSpace(state.InputEntry.Text)
			if userinput == "" {
				dialog.ShowInformation("Error", "Input is empty!", state.MyWindow)
				return
			}
			stream, err := ConvertInputToStream(1)
			if err != nil {
				res := err.Error()
				scrollContainer.Content = widget.NewLabel(res)
				scrollContainer.Refresh()
				return
			}
			vContent := container.NewVBox()
			for i := 0; i < len(stream.SerVersionUIDs); i++ {
				entry := widget.NewEntry()
				entry.Text = fmt.Sprintf("%d", stream.SerVersionUIDs[i].SerialVersionUID)
				c := container.NewScroll(entry)
				c.SetMinSize(fyne.NewSize(200, 10))
				hContent := container.NewHBox(widget.NewLabel(stream.SerVersionUIDs[i].ClassName), c)
				vContent.Add(hContent)
			}
			scrollContainer.Content = vContent
			scrollContainer.Refresh()
		} else {
			scrollContainer.Content = empty
			scrollContainer.Refresh()
		}
	}
	//sumbitButton and outputEntry
	outputEntry := widget.NewMultiLineEntry()
	maxCharsPerLine := 200
	outputEntry.OnChanged = func(content string) {
		lines := strings.Split(content, "\n")
		var wrappedLines []string
		for _, line := range lines {
			for len(line) > maxCharsPerLine {
				wrappedLines = append(wrappedLines, line[:maxCharsPerLine])
				line = line[maxCharsPerLine:]
			}
			wrappedLines = append(wrappedLines, line)
		}
		newContent := strings.Join(wrappedLines, "\n")
		if newContent != content {
			outputEntry.SetText(newContent)
		}
	}
	outputContainer := container.NewScroll(outputEntry)
	outputContainer.SetMinSize(fyne.NewSize(800, 10))
	changeButton := widget.NewButton("change", func() {
		model := 1
		userinput = strings.TrimSpace(state.InputEntry.Text)
		if userinput == "" {
			dialog.ShowInformation("Error", "Input is empty!", state.MyWindow)
			return
		}
		stream, err = ConvertInputToStream(1)

		if err != nil {
			res := err.Error()
			outputEntry.SetText(res)
			return
		}
		if funcDirtyCheck.Checked {
			stream, err = insertDirtyData(funcDirtyInputEntry.Text, stream)
		}
		if radio.Selected == "2" {
			model = 2
		} else if radio.Selected == "3" {
			model = 3
		}
		if funcUidCheck.Checked {
			var UIDs []int64
			content := scrollContainer.Content
			if vbox, ok := content.(*fyne.Container); ok {
				for _, obj := range vbox.Objects {
					if hbox, ok := obj.(*fyne.Container); ok {
						for _, objInHBox := range hbox.Objects {
							if obj, ok := objInHBox.(*container.Scroll); ok {
								entry := obj.Content
								if entry, ok := entry.(*widget.Entry); ok {
									value, _ := strconv.ParseInt(entry.Text, 10, 64)
									UIDs = append(UIDs, value)
								}
							}
						}
					}
				}
			}
			for i := 0; i < len(UIDs); i++ {
				if UIDs[i] != stream.SerVersionUIDs[i].SerialVersionUID {
					stream.SerVersionUIDs[i].StructPtr.SerialVersionUID = UIDs[i]
				}
			}
		}
		res, err := ConvertStreamToBase64(stream, model)
		if err != nil {
			res := err.Error()
			outputEntry.SetText(res)
			return
		}
		outputEntry.SetText(res)
	})
	copyButton := widget.NewButton("Copy", func() {
		state.MyWindow.Clipboard().SetContent(strings.ReplaceAll(outputEntry.Text, "\n", ""))
	})
	return container.NewVBox(line, funcDirtyTheme, funcDirtyIntroduce1, funcDirtyIntroduce2, funcDirtyCheck, container.NewHBox(funcDirtyInputLabel, empty, funcDirtyInputEntry), line1, funcUtfTheme, funcUtfIntroduce1, funcUtfIntroduce2, funcUtfCheck, funcUtfLabel, container.NewHBox(empty, radio), line2, funcUidTheme, funcUidIntroduce1, funcUidIntroduce2, funcUidCheck, scrollContainer, line3, changeButton, outputContainer, copyButton)
}

func wrapText(input string, maxLineLength int) string {
	var result string
	for i, r := range input {
		if i > 0 && i%maxLineLength == 0 {
			result += "\n"
		}
		result += string(r)
	}
	return result
}

func InitUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SerializeJava")
	state := &AppState{MyWindow: myWindow}

	mainContainer := NewContainer(state)
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, mainContainer))
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
