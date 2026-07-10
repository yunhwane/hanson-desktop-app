// Package hanson 은 한국어를 "한슨어"로 변환하는 규칙 기반 번역기입니다.
//
// 한슨어 규칙: 문장의 "마지막 음절"을 무조건 "슨"으로 교체합니다.
// (뒤에 덧붙이는 게 아니라, 마지막 글자 자체를 슨으로 바꿉니다.)
// 문장 끝 문장부호(!, ?, . 등)는 그대로 뒤에 유지합니다.
//
//	안녕하세요  ->  안녕하세슨
//	존나 웃기네 ->  존나 웃기슨
//	과자        ->  과슨
//	아웃겨!     ->  아웃슨!
package hanson

import "strings"

// 문장 끝에 붙을 수 있는 문장부호 집합입니다.
const punctuation = ".!?~…。！？"

// Translate 는 여러 줄/여러 문장으로 이루어진 한국어 텍스트를 한슨어로 변환합니다.
func Translate(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = translateLine(line)
	}
	return strings.Join(lines, "\n")
}

// translateLine 은 한 줄을 문장 단위로 나눠 각 문장을 변환합니다.
// (예: "아웃겨! 대박" 처럼 문장부호로 끝나는 문장이 여러 개일 때 각각 처리)
func translateLine(line string) string {
	if strings.TrimSpace(line) == "" {
		return line // 빈 줄/공백 줄은 그대로 둡니다.
	}

	var b, sentence strings.Builder
	flush := func() {
		if sentence.Len() > 0 {
			b.WriteString(translateSentence(sentence.String()))
			sentence.Reset()
		}
	}

	for _, r := range line {
		sentence.WriteRune(r)
		if strings.ContainsRune(punctuation, r) {
			flush()
		}
	}
	flush()
	return b.String()
}

// translateSentence 은 한 문장에서 마지막 음절을 "슨"으로 바꿉니다.
// 앞 공백과 끝 문장부호는 원래 위치에 그대로 보존합니다.
//
//	"아웃겨!" -> 앞:"" / 알맹이:"아웃겨" / 뒤:"!"  ->  "아웃슨" + "!"  ->  "아웃슨!"
func translateSentence(sentence string) string {
	// 1) 앞쪽 공백 분리
	trimmed := strings.TrimLeft(sentence, " \t")
	leading := sentence[:len(sentence)-len(trimmed)]

	// 2) 끝쪽 공백 + 문장부호 분리 (나중에 그대로 복원)
	runes := []rune(trimmed)
	end := len(runes)
	for end > 0 {
		r := runes[end-1]
		if r == ' ' || r == '\t' || strings.ContainsRune(punctuation, r) {
			end--
			continue
		}
		break
	}
	core := string(runes[:end])
	trailing := string(runes[end:])

	if core == "" {
		return sentence // 알맹이가 없으면(공백/부호만) 그대로.
	}

	return leading + slap(core) + trailing
}

// slap 은 어절(문장의 핵심 부분)의 마지막 음절을 "슨"으로 교체합니다.
// 마지막 글자가 한글 음절이면 교체하고, 아니면(영어/숫자/기호 등) 뒤에 붙입니다.
//
//	웃기네 -> 웃기슨,  과자 -> 과슨,  아웃겨 -> 아웃슨,  안녕하세슨 -> 안녕하세슨(그대로)
func slap(core string) string {
	cr := []rune(core)
	last := cr[len(cr)-1]
	if isHangulSyllable(last) {
		cr[len(cr)-1] = '슨'
		return string(cr)
	}
	return core + "슨"
}

// isHangulSyllable 은 룬이 완성형 한글 음절(가~힣)인지 판별합니다.
func isHangulSyllable(r rune) bool {
	return r >= 0xAC00 && r <= 0xD7A3
}
