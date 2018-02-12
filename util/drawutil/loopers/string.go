package loopers

import (
	"fmt"
	"image"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// glyph metrics
// https://developer.apple.com/library/content/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyph_metrics_2x.png

type String struct {
	EmbedLooper

	Face    font.Face
	Str     string
	Ri      int // will be equal to len(Str) at end-of-string rune
	Ru      rune
	PrevRu  rune
	Pen     fixed.Point26_6 // upper left corner
	Kern    fixed.Int26_6
	Advance fixed.Int26_6
	Metrics font.Metrics

	// use externally to help detect extra runes (ex: wraplinerune, annotations)
	riClone int
}

func MakeString(face font.Face, str string) String {
	return String{
		Face:    face,
		Str:     str,
		Metrics: face.Metrics(),
	}
}

func (lpr *String) Loop(fn func() bool) {
	lpr.OuterLooper().Loop(func() bool {
		if lpr.Ri > len(lpr.Str) {
			panic(fmt.Sprintf("ri>len: %v %v", lpr.Ri, len(lpr.Str)))
		}
		if lpr.Ri == len(lpr.Str) {
			// end of string
			lpr.Ru = 0
			lpr.riClone = 0
			_ = lpr.Iterate(fn)
			return false
		}
		ru, w := utf8.DecodeRuneInString(lpr.Str[lpr.Ri:])
		lpr.Ru = ru
		lpr.riClone = 0
		if !lpr.Iterate(fn) {
			return false
		}
		lpr.Ri += w
		return true
	})
}
func (lpr *String) Iterate(fn func() bool) bool {
	lpr.AddKern()
	lpr.CalcAdvance()
	if ok := fn(); !ok {
		return false
	}
	lpr.PrevRu = lpr.Ru
	lpr.Pen.X = lpr.PenXAdvance()
	return true
}
func (lpr *String) AddKern() {
	lpr.Kern = lpr.Face.Kern(lpr.PrevRu, lpr.Ru)
	lpr.Pen.X += lpr.Kern
}
func (lpr *String) CalcAdvance() bool {
	adv, ok := lpr.Face.GlyphAdvance(lpr.Ru)
	if !ok {
		lpr.Advance = 0
		return false
	}
	lpr.Advance = adv
	return true
}
func (lpr *String) PenXAdvance() fixed.Int26_6 {
	return lpr.Pen.X + lpr.Advance
}
func (lpr *String) PenBounds() *fixed.Rectangle26_6 {
	var r fixed.Rectangle26_6
	r.Min.X = lpr.Pen.X
	r.Max.X = lpr.PenXAdvance()
	r.Min.Y = lpr.LineY0()
	r.Max.Y = lpr.LineY1()
	return &r
}
func (lpr *String) PenBoundsForImage() *image.Rectangle {
	pb := lpr.PenBounds()

	// both min and max should use the same function (floor/ceil/round) because while the first rune uses ceil on max, if the next rune uses floor on min it will overwrite the previous rune on one pixel. This is noticeable in painting backgrounds.
	min := image.Point{pb.Min.X.Round(), pb.Min.Y.Round()}
	max := image.Point{pb.Max.X.Round(), pb.Max.Y.Round()}

	r := image.Rect(min.X, min.Y, max.X, max.Y)
	return &r
}

func (lpr *String) Baseline() fixed.Int26_6 {
	return lpr.Metrics.Ascent
}
func (lpr *String) LineHeight() fixed.Int26_6 {
	lh := lpr.Baseline() + lpr.Metrics.Descent
	// line height needs to be aligned with an int to have predictable line positions to be used in calculations.
	return fixed.I(lh.Ceil())
}
func (lpr *String) LineY0() fixed.Int26_6 {
	return lpr.Pen.Y
}
func (lpr *String) LineY1() fixed.Int26_6 {
	return lpr.LineY0() + lpr.LineHeight()
}

func (lpr *String) PushRiClone() {
	lpr.riClone++
}
func (lpr *String) PopRiClone() {
	lpr.riClone--
}
func (lpr *String) IsRiClone() bool {
	return lpr.riClone > 0
}

// Implements PosDataKeeper
func (lpr *String) KeepPosData() interface{} {
	d := &StringData{
		Ri:     lpr.Ri,
		PrevRu: lpr.PrevRu,
		Pen:    lpr.Pen,
	}
	return d
}

// Implements PosDataKeeper
func (lpr *String) RestorePosData(data interface{}) {
	d := data.(*StringData)
	lpr.Ri = d.Ri
	lpr.PrevRu = d.PrevRu
	lpr.Pen = d.Pen
}

type StringData struct {
	Ri     int
	PrevRu rune
	Pen    fixed.Point26_6
}
