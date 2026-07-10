<div align="center">

<img src="icon.png" width="128" alt="한슨어 번역기" />

# 한슨어 번역기 🐢

**한국어를 한슨어로 바꿔주는 데스크톱 앱**

인터넷 · AI · 설치 과정 없이, 오프라인으로 즉시 번역합니다.

[![Download](https://img.shields.io/badge/⬇️_다운로드-Releases-368A60?style=for-the-badge)](https://github.com/yunhwane/hanson-desktop-app/releases/latest)

</div>

---

## 한슨어가 뭔가슨?

문장의 **마지막 음절을 `슨`으로 교체**하는 말투입니다. 문장부호는 그대로 유지돼요.

| 한국어 | 한슨어 |
|--------|--------|
| 안녕하세요 | 안녕하세**슨** |
| 존나 웃기네 | 존나 웃기**슨** |
| 과자 | 과**슨** |
| 아웃겨! | 아웃**슨**! |
| 밥 먹었어요? | 밥 먹었어**슨**? |

## 다운로드 (비개발자용)

[**최신 릴리스**](https://github.com/yunhwane/hanson-desktop-app/releases/latest)에서 운영체제에 맞는 파일을 받으세요.

| OS | 파일 | 실행 방법 |
|----|------|-----------|
| 🪟 Windows | **`HansonTranslator.exe`** (압축 없음, 추천) | 다운로드 후 **바로 더블클릭** |
| 🪟 Windows | `HansonTranslator-windows-amd64.zip` | 압축 풀고 `.exe` 더블클릭 |
| 🍎 macOS (Apple Silicon) | `HansonTranslator-macOS-arm64.zip` | 압축 풀고 `.app` 더블클릭 |

> **⚠️ 처음 실행 시 경고가 뜰 수 있어요** (코드 서명이 없는 앱이라 그렇습니다)
> - **Windows**: "Windows의 PC 보호" 창 → **추가 정보 → 실행**
> - **macOS**: 앱 **우클릭 → 열기** → **열기** (최초 1회만)

### 압축이 안 풀리거나 "파일이 올바르지 않다"고 할 때

Windows 탐색기 기본 압축 풀기는 가끔 오작동합니다. zip을 **우클릭 → 압축 풀기**로 풀거나,
그래도 안 되면 압축 파일 안의 `.exe`를 바탕화면 등으로 **드래그해서 꺼내** 실행하세요.

## 사용법

1. 왼쪽 `한국어` 칸에 문장을 입력합니다.
2. 오른쪽 `한슨어` 칸에 **실시간으로** 번역 결과가 나옵니다.
3. `결과 복사` 버튼으로 복사해서 어디든 붙여넣으세요.

---

## 개발자용 (직접 빌드)

Go로 작성됐고 GUI는 [Fyne](https://fyne.io)을 사용합니다.

### 실행

```bash
brew install go            # macOS (Windows/Linux는 https://go.dev/dl)
git clone https://github.com/yunhwane/hanson-desktop-app.git
cd hanson-desktop-app
go mod tidy
go run .
```

### 테스트

```bash
go test ./...
```

### 배포용 패키징

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
go run ./cmd/genicon          # 아이콘 icon.png 생성(변경 시에만)

# 현재 OS용 패키지 (FyneApp.toml 의 이름/아이콘/버전 사용)
fyne package -os darwin       # macOS .app
fyne package -os windows      # Windows .exe (Windows에서 실행 시)
```

Mac에서 **Windows용 `.exe`를 크로스 컴파일**하려면 (Docker 필요):

```bash
go install github.com/fyne-io/fyne-cross@latest
fyne-cross windows -arch=amd64
# 결과: fyne-cross/dist/windows-amd64/한슨어 번역기.exe.zip
```

## 번역 규칙을 바꾸고 싶다면

핵심 로직은 [`internal/hanson/translate.go`](internal/hanson/translate.go)에 있습니다.
`slap()` 함수가 문장의 마지막 음절을 `슨`으로 교체합니다. 규칙을 바꾼 뒤 `go test ./...`로 검증하세요.

## 프로젝트 구조

```
hanson-desktop-app/
├── main.go                     # Fyne GUI (테마 · 마스코트 · 실시간 번역 · 복사)
├── internal/hanson/
│   ├── translate.go            # 한슨어 변환 규칙 (핵심 로직)
│   └── translate_test.go       # 테스트
├── cmd/genicon/main.go         # 거북이 아이콘 생성기 (순수 Go)
├── icon.png                    # 앱 아이콘 (앱에 임베드됨)
├── FyneApp.toml                # 패키징 메타데이터
└── go.mod
```

## 라이선스

[MIT](LICENSE)
