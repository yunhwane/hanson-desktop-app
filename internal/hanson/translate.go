// Package hanson 은 한국어를 "한슨어"로 변환하는 규칙 기반 번역기입니다.
//
// 말투(Style)는 두 가지입니다.
//
//	슨체:  문장 마지막 "음절"을 무조건 "슨"으로 교체 (한슨어 고유의 개그)
//	누체:  마지막 어절의 "어미"를 떼고 어간에 "누~"를 붙임 (디시 "~노/누" 말투)
//
// 두 말투 모두 문장 끝 문장부호(!, ?, . 등)는 그대로 뒤에 유지합니다.
//
//	안녕하세요   ->  안녕하세슨   (슨체)  /  안녕하누~    (누체)
//	감사합니다   ->  감사합니슨   (슨체)  /  감사하누~    (누체)
//	밥 먹었어요? ->  밥 먹었어슨? (슨체)  /  밥 먹었누~?  (누체)
//	과자         ->  과슨         (슨체)  /  과누~        (누체)
package hanson

import "strings"

// 문장 끝에 붙을 수 있는 문장부호 집합입니다.
const punctuation = ".!?~…。！？"

// Style 은 한슨어 변환 말투(종류)를 나타냅니다.
// 내부 규칙(tail/transform)은 감추고, 화면 표시용 Label 만 노출합니다.
type Style struct {
	Label string // 화면 표시용 이름 (예: "슨체", "누체")
	tail  string // 마지막 어절 뒤·문장부호 앞에 붙일 꼬리 (예: "", "~")
	// transform 은 마지막 어절(공백으로 구분된 마지막 단어) 하나를 이 말투로 바꿉니다.
	transform func(word string) string
}

// 기본 제공하는 말투들.
var (
	// StyleSeun 은 기본 "슨체" 입니다. (마지막 음절 -> 슨)
	StyleSeun = Style{Label: "슨체", tail: "", transform: seunWord}
	// StyleNu 는 "누체" 입니다. (어미 분석 후 어간 + 누, 끝에 ~)
	StyleNu = Style{Label: "누체", tail: "~", transform: nuWord}
)

// Styles 는 UI 선택기 등에서 쓸 수 있는 전체 말투 목록입니다. (표시 순서 유지)
var Styles = []Style{StyleSeun, StyleNu}

// Translate 는 여러 줄/여러 문장으로 이루어진 한국어 텍스트를 주어진 말투(style)로 변환합니다.
func Translate(text string, style Style) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = translateLine(line, style)
	}
	return strings.Join(lines, "\n")
}

// translateLine 은 한 줄을 문장 단위로 나눠 각 문장을 변환합니다.
// (예: "아웃겨! 대박" 처럼 문장부호로 끝나는 문장이 여러 개일 때 각각 처리)
func translateLine(line string, style Style) string {
	if strings.TrimSpace(line) == "" {
		return line // 빈 줄/공백 줄은 그대로 둡니다.
	}

	var b, sentence strings.Builder
	flush := func() {
		if sentence.Len() > 0 {
			b.WriteString(translateSentence(sentence.String(), style))
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

// translateSentence 은 한 문장의 마지막 어절을 말투에 맞게 바꿉니다.
// 앞 공백과 끝 문장부호는 원래 위치에 그대로 보존합니다.
//
//	"아웃겨!" -> 앞:"" / 알맹이:"아웃겨" / 뒤:"!"  ->  변환("아웃겨") + "!"
func translateSentence(sentence string, style Style) string {
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

	// 3) 마지막 어절만 변환 (앞 어절들은 그대로 유지)
	head, word := splitLastWord(core)
	converted := style.transform(word)

	// 4) 꼬리(tail)는 변환된 어절 뒤·문장부호 앞에 붙입니다.
	//    다만 문장부호가 이미 같은 꼬리(예: "~")로 시작하면 중복을 피합니다.
	tail := style.tail
	if tail != "" && strings.HasPrefix(trailing, tail) {
		tail = ""
	}

	return leading + head + converted + tail + trailing
}

// splitLastWord 는 core 를 "마지막 공백까지의 앞부분(head)"과 "마지막 어절(word)"로 나눕니다.
// 공백이 없으면 head 는 빈 문자열입니다.
func splitLastWord(core string) (head, word string) {
	if i := strings.LastIndexAny(core, " \t"); i >= 0 {
		return core[:i+1], core[i+1:]
	}
	return "", core
}

// ── 슨체 ────────────────────────────────────────────────────────────────────

// seunWord 는 어절의 마지막 음절을 "슨"으로 교체합니다. (한슨어 고유 규칙)
func seunWord(word string) string {
	return slap(word, '슨')
}

// slap 은 어절의 마지막 음절을 교체 글자(replace)로 바꿉니다.
// 마지막 글자가 한글 음절이면 교체하고, 아니면(영어/숫자/기호 등) 뒤에 붙입니다.
//
//	웃기네 -> 웃기슨,  과자 -> 과슨,  아웃겨 -> 아웃슨,  안녕하세슨 -> 안녕하세슨(그대로)
func slap(word string, replace rune) string {
	wr := []rune(word)
	last := wr[len(wr)-1]
	if isHangulSyllable(last) {
		wr[len(wr)-1] = replace
		return string(wr)
	}
	return word + string(replace)
}

// ── 누체 (어미 분석) ─────────────────────────────────────────────────────────

// nuEndings 는 어간에서 통째로 떼어낼 "2음절 이상" 종결/공손 어미 목록입니다.
// 반드시 긴 것부터(longest-first) 두어 그리디 매칭이 올바르게 동작하게 합니다.
//
// 1음절 어미(다/네/지/요/어 등)는 넣지 않습니다. 그런 경우 어간+누 == "마지막 음절
// 교체"라 fallback(slap)이 같은 결과를 내기 때문입니다. (예: 웃기네 -> 웃기누)
var nuEndings = []string{
	// 하십시오체(격식체)
	"습니다", "습니까", "읍니다", "읍니까",
	// 해요체
	"이에요", "으세요", "으셔요",
	"세요", "셔요", "예요", "에요",
	"어요", "아요", "여요",
	"네요", "군요", "지요", "고요", "구요", "나요",
	// 해체/한다체
	"는다", "구나",
}

// nuWord 는 어절의 어미를 분석해 어간에 "누"를 붙입니다.
// 어미를 못 찾으면 마지막 음절을 "누"로 교체합니다(fallback).
func nuWord(word string) string {
	if stem, ok := stripNuEnding(word); ok {
		return stem + "누"
	}
	return slap(word, '누')
}

// stripNuEnding 은 어절에서 종결/공손 어미를 떼어 어간을 돌려줍니다.
// 규칙 기반 휴리스틱(사전 없음)이라 완벽하진 않지만 흔한 말끝은 잘 잡습니다.
func stripNuEnding(word string) (string, bool) {
	// 1) 2음절+ 어미를 긴 것부터 잘라내기
	for _, suf := range nuEndings {
		if stem, ok := strings.CutSuffix(word, suf); ok && endsWithHangul(stem) {
			return stem, true
		}
	}
	// 2) 받침에 융합된 어미 처리 (합니다->하, 갑니까->가, 간다->가, 한다->하)
	if stem, ok := stripFusedEnding(word); ok {
		return stem, true
	}
	return word, false
}

// stripFusedEnding 은 받침으로 축약된 어미(ㅂ니다/ㅂ니까/ㄴ다)를 처리합니다.
// 예) 합니다 = 하 + ㅂ니다  ->  ㅂ 받침을 떼어 "하"
//
//	간다  = 가 + ㄴ다   ->  ㄴ 받침을 떼어 "가"
func stripFusedEnding(word string) (string, bool) {
	rs := []rune(word)
	n := len(rs)

	// ...ㅂ니다 / ...ㅂ니까  (앞 음절 받침이 ㅂ 또는 ㄴ)
	if n >= 3 {
		if tail2 := string(rs[n-2:]); tail2 == "니다" || tail2 == "니까" {
			if cho, jung, jong, ok := decompose(rs[n-3]); ok && (jong == jongB || jong == jongN) {
				rs[n-3] = compose(cho, jung, jongNone)
				return string(rs[:n-2]), true
			}
		}
	}

	// ...ㄴ다  (앞 음절 받침이 ㄴ:  간다->가, 한다->하)
	if n >= 2 && rs[n-1] == '다' {
		if cho, jung, jong, ok := decompose(rs[n-2]); ok && jong == jongN {
			rs[n-2] = compose(cho, jung, jongNone)
			return string(rs[:n-1]), true
		}
	}

	return word, false
}

// endsWithHangul 은 문자열이 비어있지 않고 한글 음절로 끝나는지 봅니다.
// (어간이 남아있고, 그 끝이 한글이어야 "누"를 자연스럽게 붙일 수 있습니다.)
func endsWithHangul(s string) bool {
	rs := []rune(s)
	return len(rs) > 0 && isHangulSyllable(rs[len(rs)-1])
}

// ── 한글 자모 유틸 ───────────────────────────────────────────────────────────

const (
	hangulBase = 0xAC00
	hangulLast = 0xD7A3
	jungCount  = 21 // 중성 개수
	jongCount  = 28 // 종성 개수(받침 없음 포함)

	jongNone = 0  // 받침 없음
	jongN    = 4  // 받침 ㄴ
	jongB    = 17 // 받침 ㅂ
)

// decompose 는 완성형 한글 음절을 (초성, 중성, 종성) 인덱스로 분해합니다.
func decompose(r rune) (cho, jung, jong int, ok bool) {
	if r < hangulBase || r > hangulLast {
		return 0, 0, 0, false
	}
	s := int(r) - hangulBase
	cho = s / (jungCount * jongCount)
	jung = (s % (jungCount * jongCount)) / jongCount
	jong = s % jongCount
	return cho, jung, jong, true
}

// compose 는 (초성, 중성, 종성) 인덱스로 완성형 한글 음절을 만듭니다.
func compose(cho, jung, jong int) rune {
	return rune(hangulBase + (cho*jungCount+jung)*jongCount + jong)
}

// isHangulSyllable 은 룬이 완성형 한글 음절(가~힣)인지 판별합니다.
func isHangulSyllable(r rune) bool {
	return r >= hangulBase && r <= hangulLast
}
