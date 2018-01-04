package widget

import (
	"image/color"

	"github.com/jmigpin/editor/util/imageutil"
)

type Label struct {
	EmbedNode
	Text   *BasicText
	Border *Pad
	Pad    *Pad
	Color  *color.Color
	ctx    Context
}

func NewLabel(ctx Context) *Label {
	l := &Label{ctx: ctx}
	l.Text = NewBasicText(ctx)
	l.Pad = NewPad(ctx, l.Text)
	l.Border = NewPad(ctx, l.Pad)
	l.Append(l.Border)
	return l
}
func (l *Label) Paint() {
	if l.Color == nil {
		return
	}
	imageutil.FillRectangle(l.ctx.Image(), &l.Bounds, *l.Color)
}
