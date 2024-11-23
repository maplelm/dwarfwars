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
	buttonBounds    []rl.Rectangle
	Columns         int32
	Scale           float32
	Font            rl.Font
}

func NewList(p rl.Vector2, f rl.Font, cols int32, s float32, borderwidth int32, bordercolor rl.Color, foreground rl.Color, fontsize int32, fontcolor rl.Color, buttonsize rl.Vector2) *List {
	return &List{
		Position:        p,
		Font:            f,
		Columns:         cols,
		Scale:           s,
		Buttons:         make([]Button, 0),
		buttonBounds:    make([]rl.Rectangle, 0),
		BorderWidth:     borderwidth,
		BorderColor:     bordercolor,
		ForegroundColor: foreground,
		FontSize:        fontsize,
		FontColor:       fontcolor,
		ButtonSize:      buttonsize,
	}
}

func (l *List) Add(lab string, a func()) int {

	l.Buttons = append(l.Buttons, *New(
		lab,
		a,
	))
	x := float32(int32(len(l.Buttons)) % l.Columns)
	y := float32(math.Floor(float64(len(l.Buttons)) / float64(l.Columns)))
	l.buttonBounds = append(l.buttonBounds, rl.NewRectangle(
		l.Position.X+x*(l.ButtonSize.X+float32(l.BorderWidth)*2),
		l.Position.Y+y*(l.ButtonSize.Y+float32(l.BorderWidth)*2),
		l.ButtonSize.X,
		l.ButtonSize.Y,
	))
	return len(l.Buttons)
}

func (l *List) Move(v rl.Vector2) {
	l.Position = rl.Vector2Add(l.Position, v)
	for i := range l.Buttons {
		bb := l.buttonBounds[i]
		bb.X += v.X
		bb.Y += v.Y
		l.buttonBounds[i] = bb
	}
}

func (l *List) Update() {
	for i, b := range l.Buttons {
		b.Update(l.buttonBounds[i])
	}
}

func (l *List) Draw() {
	for i, v := range l.Buttons {
		v.Draw(l.buttonBounds[i], l.BorderWidth, l.BorderColor, l.ForegroundColor, l.Font, l.FontSize, l.FontColor)
	}
}
