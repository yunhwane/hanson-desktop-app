package main

import (
	_ "embed"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"hanson-desktop/internal/hanson"
)

// 타이핑이 멈춘 뒤 번역/화면갱신까지 기다리는 시간.
// 키 입력마다 무거운 렌더링을 하지 않도록 이 시간만큼 디바운스합니다.
const debounceDelay = 80 * time.Millisecond

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
	// TextWrapWord + Multiline 은 키 입력마다 줄바꿈을 재계산해 타이핑이 느려지는
	// 알려진 Fyne 성능 버그(#4221, #5297)가 있습니다. 입력창은 워드랩을 꺼서
	// 타이핑을 매끄럽게 유지합니다(긴 줄은 가로 스크롤).
	input.Wrapping = fyne.TextWrapOff

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.SetPlaceHolder("번역 결과가 여기에 나와슨!")

	// --- 말투(스타일) 상태 ---
	// 라디오 콜백(메인 스레드)과 디바운스 타이머 콜백(별도 goroutine)에서
	// 함께 읽으므로 뮤텍스로 보호합니다.
	var (
		styleMu sync.Mutex
		style   = hanson.Styles[0]
	)
	translate := func(text string) string {
		styleMu.Lock()
		s := style
		styleMu.Unlock()
		return hanson.Translate(text, s)
	}

	// 입력 렌더링이 빠른 타이핑을 막지 않도록 번역/화면갱신을 디바운스합니다.
	// 키를 칠 때는 타이머만 리셋하고 즉시 반환 → 입력이 매끄럽습니다.
	// 타이핑이 잠깐 멈추면(debounceDelay) 그때 한 번만 번역해 결과를 갱신합니다.
	var (
		debounceMu    sync.Mutex
		debounceTimer *time.Timer
	)
	input.OnChanged = func(text string) {
		debounceMu.Lock()
		defer debounceMu.Unlock()
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceTimer = time.AfterFunc(debounceDelay, func() {
			result := translate(text)
			// 타이머 콜백은 별도 goroutine이므로 UI 갱신은 fyne.Do로 메인 스레드에서.
			fyne.Do(func() {
				if output.Text != result {
					output.SetText(result)
				}
			})
		})
	}

	// --- 말투 선택기 (슨체 / 누체) ---
	// 고르면 즉시 현재 입력을 다시 번역합니다.
	styleNames := make([]string, len(hanson.Styles))
	for i, s := range hanson.Styles {
		styleNames[i] = s.Label
	}
	styleSelector := widget.NewRadioGroup(styleNames, func(selected string) {
		styleMu.Lock()
		for _, s := range hanson.Styles {
			if s.Label == selected {
				style = s
				break
			}
		}
		styleMu.Unlock()
		// 라디오 콜백은 메인 스레드이므로 바로 갱신해도 안전합니다.
		output.SetText(translate(input.Text))
	})
	styleSelector.Horizontal = true
	styleSelector.Required = true
	styleSelector.SetSelected(hanson.Styles[0].Label)

	styleRow := container.NewHBox(
		widget.NewLabelWithStyle("말투", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		styleSelector,
	)

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
		container.NewVBox(header, styleRow, widget.NewSeparator()), // 상단
		container.NewHBox(copyBtn, clearBtn),                       // 하단
		nil, nil,
		container.NewGridWithColumns(2, inputBox, outputBox), // 가운데(확장)
	)

	w.SetContent(content)
	w.ShowAndRun()
}
