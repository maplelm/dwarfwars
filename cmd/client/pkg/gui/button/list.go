package button

import rl "github.com/gen2brain/raylib-go/raylib"

type List struct {
	Bounds  rl.Rectangle
	Buttons []Button
	Columns int32
	Scale   float32
	Font    rl.Font
}

func NewList(bounds rl.Rectangle, f rl.Font, cols int32, s float32) *List {
	return &List{
		Bounds:  bounds,
		Font:    f,
		Columns: cols, 
		Scale:   s,
		Buttons: make([]Button, 1),
	}
}

func (l *List) Add(lab string, a func()) int {
	l.Buttons = append(l.Buttons, *New(lab, a

	}
}
