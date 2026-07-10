# 기여 가이드 (Contributing)

한슨어 번역기에 관심 가져주셔서 감사합니슨! 🐢
버그 제보, 기능 제안, 코드 기여 모두 환영합니다.

## 시작하기

1. 이 저장소를 **Fork** 하고 로컬에 클론합니다.
   ```bash
   git clone https://github.com/<your-username>/hanson-desktop-app.git
   cd hanson-desktop-app
   ```
2. Go(1.22 이상)를 설치합니다. → https://go.dev/dl
   - macOS: `brew install go`
   - macOS에서 처음 빌드 시 Xcode Command Line Tools 필요: `xcode-select --install`
3. 의존성을 받고 앱을 실행합니다.
   ```bash
   go mod tidy
   go run .
   ```

## 개발 워크플로우

1. 브랜치를 만듭니다: `git checkout -b feat/내-기능` (또는 `fix/...`).
2. 코드를 수정합니다.
3. **테스트와 정적 분석을 통과**시킵니다.
   ```bash
   go test ./...
   go vet ./...
   gofmt -l .        # 출력이 없어야 함 (있으면 `gofmt -w .` 로 정리)
   ```
4. 커밋하고 푸시한 뒤 **Pull Request**를 엽니다.

## 번역 규칙 바꾸기

한슨어 변환 로직은 [`internal/hanson/translate.go`](internal/hanson/translate.go)에 있습니다.
핵심은 `slap()` 함수로, 문장의 마지막 음절을 `슨`으로 교체합니다.

규칙을 바꾸거나 추가할 때는 **반드시 테스트를 함께 추가/수정**하세요:

```go
// internal/hanson/translate_test.go
{"새로운 입력", "기대하는 한슨어"},
```

그런 다음 `go test ./...` 로 검증합니다. 예시나 엣지 케이스가 많을수록 좋슨!

## 코드 스타일

- 표준 `gofmt` 포맷을 따릅니다.
- 주석과 UI 문구는 프로젝트 톤(친근한 한국어)에 맞춰 주세요.
- 외부 의존성 추가는 신중하게 — 이 앱은 **오프라인·경량**을 지향합니다.

## 커밋 메시지

- 제목은 간결하게, 본문에 "무엇을/왜" 바꿨는지 적어 주세요.
- 한국어/영어 모두 괜찮습니다.

## 아이콘을 바꿨다면

아이콘은 [`cmd/genicon`](cmd/genicon/main.go)에서 순수 Go로 생성합니다.
수정 후 아래로 `icon.png`를 다시 만들고 커밋하세요(앱에 임베드됩니다).

```bash
go run ./cmd/genicon
```

## 질문이 있으면

[Issues](https://github.com/yunhwane/hanson-desktop-app/issues)에 편하게 남겨 주세요.

행동 강령은 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)를 참고해 주세요.
