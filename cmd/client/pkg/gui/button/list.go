package button

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type List struct {
	Position        rl.Vector2
	ButtonSize      rl.Vector2
	BorderWidth     int32
	BorderColor     rl.Color
	FontSize        int32
	FontColor       rl.Color
	ForegroundColor rl.Color
	Buttons         []Button
	Columns         int32
	Font            rl.Font
}

func NewList(p, buttonsize rl.Vector2, cols, borderwidth, fontsize int32, bordercolor, foreground, fontcolor rl.Color, f rl.Font) *List {
	return &List{
		Position:        p,
		ButtonSize:      buttonsize,
		Font:            f,
		Columns:         cols,
		Buttons:         make([]Button, 0),
		BorderWidth:     borderwidth,
		BorderColor:     bordercolor,
		ForegroundColor: foreground,
		FontSize:        fontsize,
		FontColor:       fontcolor,
	}
}

func (l *List) Add(lab string, a func()) int {
	l.Buttons = append(l.Buttons, Button{
		Label:   lab,
		Clicked: false,
		Action:  a,
	})
	return len(l.Buttons)
}

func (l *List) Move(v rl.Vector2) {
	l.Position = rl.Vector2Add(l.Position, v)
}

func (l *List) Update() {
	for i, b := range l.Buttons {
		doubleBW := float32(l.BorderWidth) * 2
		localX := (float32(int32(i) % l.Columns)) * (float32(l.ButtonSize.X + doubleBW))
		localY := float32(math.Floor(float64(float32(i)/float32(l.Columns)))) * (l.ButtonSize.Y + doubleBW)

		b.UpdateWithBounds(rl.NewRectangle(
			l.Position.X+localX,
			l.Position.Y+localY,
			l.ButtonSize.X,
			l.ButtonSize.Y,
		))
	}
}

func (l *List) Size() rl.Vector2 {
	var (
		w          float32
		h          float32
		trueWidth  float32 = (l.ButtonSize.X + (float32(l.BorderWidth) * 2))
		trueHeight float32 = (l.ButtonSize.Y + (float32(l.BorderWidth) * 2))
		rows       float32 = float32(math.Floor(float64(len(l.Buttons)) / float64(l.Columns)))
	)

	if len(l.Buttons) >= int(l.Columns) {
		w = float32(l.Columns) * trueWidth
	} else {
		w = float32(len(l.Buttons)) * trueWidth
	}

	if len(l.Buttons) == 0 {
		h = 0
	} else {
		h = trueHeight + trueHeight*rows
	}
	return rl.Vector2{
		X: w,
		Y: h,
	}
}

func (l *List) Draw() {
	for i, v := range l.Buttons {
		doubleBW := float32(l.BorderWidth) * 2
		localX := (float32(int32(i) % l.Columns)) * (float32(l.ButtonSize.X + doubleBW))
		localY := float32(math.Floor(float64(float32(i)/float32(l.Columns)))) * (l.ButtonSize.Y + doubleBW)

		v.DrawWithGraphics(*NewbGraphics(rl.NewRectangle(l.Position.X+localX, l.Position.Y+localY, l.ButtonSize.X, l.ButtonSize.Y), l.BorderWidth, l.BorderColor, l.ForegroundColor, l.Font, l.FontSize, l.FontColor))

	}
}
