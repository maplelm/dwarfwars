package button

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type List struct {
	Position        rl.Vector2
	BorderWidth     int32
	BorderColor     rl.Color
	FontSize        int32
	FontColor       rl.Color
	ForegroundColor rl.Color
	Buttons         []Button
	Columns         int32
	Font            rl.Font

	size       rl.Vector2 // Cache this value as it is complex to calculate
	buttonsize rl.Vector2
}

func NewList(p, buttonsize rl.Vector2, cols, borderwidth, fontsize int32, bordercolor, foreground, fontcolor rl.Color, f rl.Font) *List {
	return &List{
		Position:        p,
		buttonsize:      buttonsize,
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
	l.size = l.Size()
	return len(l.Buttons)
}

func (l *List) Move(v rl.Vector2) {
	l.Position = rl.Vector2Add(l.Position, v)
}

func (l *List) SetButtonSize(s rl.Vector2) {
	l.buttonsize = s
	l.size = l.Size()
}

func (l *List) SetButtonSizeX(x float32) {
	l.buttonsize.X = x
	l.size = l.Size()
}

func (l *List) SetButtonSizeY(y float32) {
	l.buttonsize.Y = y
	l.size = l.Size()
}

func (l *List) Update(mp rl.Vector2) {
	if !l.IsHovered(mp) {
		return
	}
	for i, b := range l.Buttons {
		doubleBW := float32(l.BorderWidth) * 2
		localX := (float32(int32(i) % l.Columns)) * (float32(l.buttonsize.X + doubleBW))
		localY := float32(math.Floor(float64(float32(i)/float32(l.Columns)))) * (l.buttonsize.Y + doubleBW)

		b.UpdateWithBounds(mp, rl.NewRectangle(
			l.Position.X+localX,
			l.Position.Y+localY,
			l.buttonsize.X,
			l.buttonsize.Y,
		))
	}
}

func (l *List) IsHovered(mp rl.Vector2) bool {
	return mp.X > l.Position.X && mp.X < l.Position.X+l.size.X && mp.Y > l.Position.Y && mp.Y < l.Position.Y+l.size.Y
}

func (l *List) Size() rl.Vector2 {
	if len(l.Buttons) == 0 {
		return rl.Vector2{X: 0, Y: 0}
	}

	var (
		s         rl.Vector2
		btnWidth  float32 = (l.buttonsize.X + (float32(l.BorderWidth) * 2)) + 1
		btnHeight float32 = (l.buttonsize.Y + (float32(l.BorderWidth) * 2)) + 1
		rows      float32 = float32(math.Floor(float64(len(l.Buttons)) / float64(l.Columns)))
	)

	if len(l.Buttons) >= int(l.Columns) {
		s.X = float32(l.Columns) * btnWidth
	} else {
		s.X = float32(len(l.Buttons)) * btnWidth
	}

	s.Y = btnHeight * rows

	return s
}

func (l *List) Draw() {
	for i, v := range l.Buttons {
		doubleBW := float32(l.BorderWidth) * 2
		localX := (float32(int32(i) % l.Columns)) * (float32(l.buttonsize.X + doubleBW))
		localY := float32(math.Floor(float64(float32(i)/float32(l.Columns)))) * (l.buttonsize.Y + doubleBW)

		v.DrawWithGraphics(*NewbGraphics(rl.NewRectangle(l.Position.X+localX, l.Position.Y+localY, l.buttonsize.X, l.buttonsize.Y), l.BorderWidth, l.BorderColor, l.ForegroundColor, l.Font, l.FontSize, l.FontColor))

	}
}
