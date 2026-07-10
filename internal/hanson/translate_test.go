package hanson

import "testing"

func TestTranslate(t *testing.T) {
	cases := []struct{ in, want string }{
		// 마지막 음절을 무조건 '슨'으로 교체
		{"안녕하세요", "안녕하세슨"},
		{"감사합니다", "감사합니슨"},
		{"존나 웃기네", "존나 웃기슨"},
		{"과자", "과슨"},
		{"모자", "모슨"},
		{"어머니", "어머슨"},

		// 문장부호는 끝에 그대로 유지 (음절 교체는 부호 앞에서)
		{"아웃겨!", "아웃슨!"},
		{"밥 먹었어요?", "밥 먹었어슨?"},
		{"반갑습니다.", "반갑습니슨."},
		{"대박~", "대슨~"},

		// 이미 '슨'으로 끝나면 그대로 (멱등)
		{"안녕하세슨", "안녕하세슨"},

		// 여러 문장 / 빈 입력
		{"아웃겨! 대박이다", "아웃슨! 대박이슨"},
		{"안녕하세요. 반가워요!", "안녕하세슨. 반가워슨!"},
		{"", ""},
	}
	for _, c := range cases {
		if got := Translate(c.in); got != c.want {
			t.Errorf("Translate(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTranslateMultiline(t *testing.T) {
	in := "안녕하세요\n밥 먹었어요"
	want := "안녕하세슨\n밥 먹었어슨"
	if got := Translate(in); got != want {
		t.Errorf("Translate(%q) = %q, want %q", in, got, want)
	}
}
