// Command genicon 은 한슨어 번역기의 앱 아이콘(icon.png)을 순수 Go로 그립니다.
//
// 실행:  go run ./cmd/genicon
// 결과:  프로젝트 루트에 icon.png (512x512) 생성
//
// 외부 이미지 파일 없이 거북이 🐢 마스코트를 도형으로 그리고,
// 부드러운 가장자리를 위해 3배 해상도로 렌더링한 뒤 축소합니다.
package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	xdraw "golang.org/x/image/draw"
)

const (
	size = 512 // 최종 아이콘 한 변 크기(px)
	ss   = 3   // 슈퍼샘플링 배율 (안티에일리어싱용)
)

func main() {
	big := image.NewNRGBA(image.Rect(0, 0, size*ss, size*ss))
	drawTurtle(big)

	// 3배 렌더 -> 512로 부드럽게 축소
	out := image.NewNRGBA(image.Rect(0, 0, size, size))
	xdraw.CatmullRom.Scale(out, out.Bounds(), big, big.Bounds(), xdraw.Over, nil)

	f, err := os.Create("icon.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, out); err != nil {
		panic(err)
	}
	println("wrote icon.png (512x512)")
}

// 색상 팔레트
var (
	colBG    = color.NRGBA{0xFF, 0xF3, 0xD6, 0xFF} // 따뜻한 크림색 배경
	colShell = color.NRGBA{0x4C, 0xAF, 0x7D, 0xFF} // 등껍질(초록)
	colRim   = color.NRGBA{0x36, 0x8A, 0x60, 0xFF} // 등껍질 테두리/무늬(진초록)
	colLimb  = color.NRGBA{0x74, 0xCF, 0x9A, 0xFF} // 머리/다리(연초록)
	colWhite = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF} // 눈 흰자
	colBlack = color.NRGBA{0x2B, 0x2B, 0x2B, 0xFF} // 눈동자
	colBlush = color.NRGBA{0xF6, 0xA6, 0x8C, 0x99} // 볼터치(반투명)
)

// drawTurtle 은 img(슈퍼샘플 해상도)에 배경 라운드 사각형과 거북이를 그립니다.
// 좌표는 512 기준으로 적고 내부에서 ss 배율을 곱합니다.
func drawTurtle(img *image.NRGBA) {
	// 배경: 모서리가 둥근 사각형
	fillRoundRect(img, 24, 24, 488, 488, 100, colBG)

	// 다리 4개 (등껍질 뒤로 먼저 그림)
	fillEllipse(img, 150, 410, 50, 44, colLimb)
	fillEllipse(img, 362, 410, 50, 44, colLimb)
	fillEllipse(img, 206, 470, 46, 40, colLimb)
	fillEllipse(img, 306, 470, 46, 40, colLimb)
	// 꼬리
	fillEllipse(img, 256, 478, 24, 18, colLimb)

	// 머리 (등껍질 위로 확실히 나오도록 위쪽에 크게)
	fillEllipse(img, 256, 196, 90, 92, colLimb)

	// 등껍질: 테두리 -> 본체
	fillEllipse(img, 256, 352, 176, 140, colRim)
	fillEllipse(img, 256, 352, 156, 122, colShell)

	// 등껍질 무늬 (가운데 + 주변 4개)
	fillEllipse(img, 256, 352, 44, 40, colRim)
	fillEllipse(img, 180, 326, 32, 29, colRim)
	fillEllipse(img, 332, 326, 32, 29, colRim)
	fillEllipse(img, 210, 404, 30, 27, colRim)
	fillEllipse(img, 302, 404, 30, 27, colRim)

	// 눈 (머리 위)
	fillEllipse(img, 228, 182, 19, 21, colWhite)
	fillEllipse(img, 284, 182, 19, 21, colWhite)
	fillEllipse(img, 233, 186, 10, 11, colBlack)
	fillEllipse(img, 289, 186, 10, 11, colBlack)

	// 볼터치
	fillEllipse(img, 202, 222, 15, 9, colBlush)
	fillEllipse(img, 310, 222, 15, 9, colBlush)
}

// --- 그리기 헬퍼 (모두 512 좌표를 받아 내부에서 ss 배율 적용, 알파 블렌딩) ---

func blend(img *image.NRGBA, x, y int, c color.NRGBA) {
	b := img.Bounds()
	if x < 0 || y < 0 || x >= b.Dx() || y >= b.Dy() {
		return
	}
	if c.A == 0xFF {
		img.SetNRGBA(x, y, c)
		return
	}
	dst := img.NRGBAAt(x, y)
	a := float64(c.A) / 255
	img.SetNRGBA(x, y, color.NRGBA{
		R: uint8(float64(c.R)*a + float64(dst.R)*(1-a)),
		G: uint8(float64(c.G)*a + float64(dst.G)*(1-a)),
		B: uint8(float64(c.B)*a + float64(dst.B)*(1-a)),
		A: 0xFF,
	})
}

func fillEllipse(img *image.NRGBA, cx, cy, rx, ry float64, c color.NRGBA) {
	cx, cy, rx, ry = cx*ss, cy*ss, rx*ss, ry*ss
	for y := int(cy - ry); y <= int(cy+ry); y++ {
		for x := int(cx - rx); x <= int(cx+rx); x++ {
			dx := (float64(x) - cx) / rx
			dy := (float64(y) - cy) / ry
			if dx*dx+dy*dy <= 1 {
				blend(img, x, y, c)
			}
		}
	}
}

func fillRoundRect(img *image.NRGBA, x0, y0, x1, y1, r float64, c color.NRGBA) {
	x0, y0, x1, y1, r = x0*ss, y0*ss, x1*ss, y1*ss, r*ss
	cxL, cxR := x0+r, x1-r
	cyT, cyB := y0+r, y1-r
	for y := int(y0); y <= int(y1); y++ {
		for x := int(x0); x <= int(x1); x++ {
			fx, fy := float64(x), float64(y)
			inside := true
			switch {
			case fx < cxL && fy < cyT:
				inside = math.Hypot(fx-cxL, fy-cyT) <= r
			case fx > cxR && fy < cyT:
				inside = math.Hypot(fx-cxR, fy-cyT) <= r
			case fx < cxL && fy > cyB:
				inside = math.Hypot(fx-cxL, fy-cyB) <= r
			case fx > cxR && fy > cyB:
				inside = math.Hypot(fx-cxR, fy-cyB) <= r
			}
			if inside {
				blend(img, x, y, c)
			}
		}
	}
}
