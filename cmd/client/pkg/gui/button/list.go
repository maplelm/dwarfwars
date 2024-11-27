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
		BorderWidth:     borderwidth,
		BorderColor:     bordercolor,
		ForegroundColor: foreground,
		FontSize:        fontsize,
		FontColor:       fontcolor,
		ButtonSize:      buttonsize,
	}
}

func (l *List) Add(lab string, a func()) int {
	b := Button{
		Label:   lab,
		Clicked: false,
		Action:  a,
	}
	x := float32(int32(len(l.Buttons)) % l.Columns)
	y := float32(math.Floor(float64(len(l.Buttons)) / float64(l.Columns)))
	b.Bounds = rl.NewRectangle(
		l.Position.X+x*(l.ButtonSize.X+float32(l.BorderWidth)*2),
		l.Position.Y+y*(l.ButtonSize.Y+float32(l.BorderWidth)*2),
		l.ButtonSize.X,
		l.ButtonSize.Y,
	)
	l.Buttons = append(l.Buttons, b)
	return len(l.Buttons)
}

func (l *List) Move(v rl.Vector2) {
	l.Position = rl.Vector2Add(l.Position, v)
}

func (l *List) Update() {
	for i, b := range l.Buttons {
		b.UpdateWithBounds(rl.NewRectangle(
			float32(int32(l.Position.X)+((int32(i)%l.Columns)*int32((l.ButtonSize.X+float32(l.BorderWidth)*2)))),
			float32(int32(l.Position.Y)+(int32(math.Floor(float64(i)/float64(l.Columns)))*int32(l.ButtonSize.Y+float32(l.BorderWidth)*2))),
			l.ButtonSize.X,
			l.ButtonSize.Y,
		))
	}
}

func (l *List) Draw() {
	for i, v := range l.Buttons {
		gridx := float32(i % int(l.Columns))
		gridy := float32(math.Floor(float64(float32(i) / float32(l.Columns))))
		x := l.Position.X
		y := l.Position.Y
		w := l.ButtonSize.X + float32(l.BorderWidth)*2
		h := l.ButtonSize.Y + float32(l.BorderWidth)*2
		v.DrawWithGraphics(rl.NewRectangle(x+(gridx*w), y+(gridy*h), l.ButtonSize.X, l.ButtonSize.Y), l.BorderWidth, l.BorderColor, l.ForegroundColor, l.Font, l.FontSize, l.FontColor)
	}
}
