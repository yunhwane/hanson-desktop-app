# 변경 이력 (Changelog)

이 프로젝트의 주요 변경 사항을 기록합니다.
형식은 [Keep a Changelog](https://keepachangelog.com/ko/1.1.0/)를 따르며,
버전은 [유의적 버전(SemVer)](https://semver.org/lang/ko/)을 따릅니다.

## [Unreleased]

## [1.0.2] - 2026-07-10

### 추가
- **릴리스 자동화** — `v*` 태그를 push하면 GitHub Actions가 macOS `.app`와
  Windows `.exe`를 자동으로 빌드·패키징(영문 파일명)하고 릴리스에 업로드합니다.
- 오픈소스 문서: `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, 이슈/PR 템플릿.
- CI 워크플로우(`gofmt`/`vet`/`test`/`build`)와 README 배지.

## [1.0.1] - 2026-07-10

### 개선
- 한글 타이핑 시 입력 렉 제거 — 입력창의 워드랩을 꺼서 Fyne의
  Multiline + `TextWrapWord` 성능 버그([#4221](https://github.com/fyne-io/fyne/issues/4221),
  [#5297](https://github.com/fyne-io/fyne/issues/5297))를 우회.
- 번역/화면 갱신을 80ms 디바운스 처리하고, `fyne.Do`로 메인 스레드에서 안전하게 갱신.
- Fyne `v2.5.2` → `v2.7.4` 업그레이드 (한글 IME 입력 처리 개선 포함).

### 수정
- 배포용 압축 파일 안의 파일명을 영문(`HansonTranslator.exe`)으로 변경하여,
  Windows 탐색기에서 *"파일이 올바르지 않습니다 / 압축이 풀리지 않음"* 문제를 해결.
- 압축 해제가 필요 없는 단일 실행 파일 `HansonTranslator.exe`를 릴리스에 추가.

## [1.0.0] - 2026-07-10

### 추가
- 한국어를 **한슨어**로 바꾸는 데스크톱 앱 첫 릴리스.
- 번역 규칙: 문장의 마지막 음절을 `슨`으로 교체 (예: `아웃겨!` → `아웃슨!`).
- Fyne 기반 GUI — 실시간 번역, 결과 복사, 초록 테마, 거북이 마스코트.
- 순수 Go 아이콘 생성기(`cmd/genicon`).
- macOS(`.app`) / Windows(`.exe`) 빌드 및 GitHub 릴리스.
- 테스트, MIT 라이선스.

[Unreleased]: https://github.com/yunhwane/hanson-desktop-app/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/yunhwane/hanson-desktop-app/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/yunhwane/hanson-desktop-app/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/yunhwane/hanson-desktop-app/releases/tag/v1.0.0
