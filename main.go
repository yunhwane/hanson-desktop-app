package main

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"hanson-desktop/internal/hanson"
)

//go:embed icon.png
var iconBytes []byte

// 브랜드 색상 (아이콘의 등껍질 초록과 맞춤)
var brandGreen = color.NRGBA{0x36, 0x8A, 0x60, 0xFF}

// hansonTheme 은 기본 테마에서 글씨를 키우고 강조색을 초록으로 바꿉니다.
type hansonTheme struct{ fyne.Theme }

func (h hansonTheme) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case theme.SizeNameText:
		return 16 // 본문 글씨 키우기 (비개발자 가독성)
	case theme.SizeNameInnerPadding:
		return 10
	}
	return h.Theme.Size(n)
}

func (h hansonTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if n == theme.ColorNamePrimary {
		return brandGreen // 커서/버튼 강조색
	}
	return h.Theme.Color(n, v)
}

func main() {
	a := app.New()
	icon := fyne.NewStaticResource("icon.png", iconBytes)
	a.SetIcon(icon)
	a.Settings().SetTheme(hansonTheme{theme.DefaultTheme()})

	w := a.NewWindow("한슨어 번역기")
	w.SetIcon(icon)
	w.Resize(fyne.NewSize(720, 480))

	// --- 헤더: 마스코트 거북이 + 제목 ---
	mascot := canvas.NewImageFromResource(icon)
	mascot.FillMode = canvas.ImageFillContain
	mascot.SetMinSize(fyne.NewSize(60, 60))

	title := canvas.NewText("한슨어 번역기", brandGreen)
	title.TextSize = 26
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("한국어를 입력하면 한슨어로 바꿔드려슨!", color.NRGBA{0x88, 0x88, 0x88, 0xFF})
	subtitle.TextSize = 13

	header := container.NewHBox(
		mascot,
		container.NewVBox(title, subtitle),
	)

	// --- 입력 / 출력 ---
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("여기에 한국어를 입력하세요...\n\n예) 안녕하세요")
	input.Wrapping = fyne.TextWrapWord

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.SetPlaceHolder("번역 결과가 여기에 나와슨!")

	// 입력할 때마다 실시간 번역
	input.OnChanged = func(text string) {
		output.SetText(hanson.Translate(text))
	}

	// --- 버튼 ---
	copyBtn := widget.NewButtonWithIcon("결과 복사", theme.ContentCopyIcon(), func() {
		if output.Text == "" {
			return
		}
		w.Clipboard().SetContent(output.Text)
		dialog.ShowInformation("복사 완료", "한슨어가 복사되었슨! 🐢", w)
	})
	copyBtn.Importance = widget.HighImportance

	clearBtn := widget.NewButtonWithIcon("지우기", theme.ContentClearIcon(), func() {
		input.SetText("")
		output.SetText("")
	})

	// --- 레이아웃 ---
	inputBox := container.NewBorder(
		widget.NewLabelWithStyle("한국어", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil, input,
	)
	outputBox := container.NewBorder(
		widget.NewLabelWithStyle("한슨어", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil, output,
	)

	content := container.NewBorder(
		container.NewVBox(header, widget.NewSeparator()), // 상단
		container.NewHBox(copyBtn, clearBtn),             // 하단
		nil, nil,
		container.NewGridWithColumns(2, inputBox, outputBox), // 가운데(확장)
	)

	w.SetContent(content)
	w.ShowAndRun()
}
