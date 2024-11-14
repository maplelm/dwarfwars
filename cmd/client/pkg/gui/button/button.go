package button

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	Label       string
	BorderWidth int32
	BorderColor rl.Color
	Color       rl.Color
	Clicked     bool
	Action      func()
}

func New(l string, a func()) *Button {
	return &Button{
		Clicked: false,
		Label:   l,
		Action:  a,
	}
}

/*
func (b *Button) Draw(bounds rl.Rectangle, f *rl.Font ) {
	rl
}
*/
