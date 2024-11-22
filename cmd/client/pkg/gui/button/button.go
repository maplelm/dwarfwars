package button

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	// Properties
	Label   string
	Clicked bool
	Action  func()

	// Graphics
	Bounds      rl.Rectangle
	BorderWidth int32
	BorderColor rl.Color
	Color       rl.Color
	Font        rl.Font
	FontSize    int32
	FontColor   rl.Color
}

func New(l string, a func(), b rl.Rectangle, bw int32, bc rl.Color, fc rl.Color, f rl.Font, fs int32, fontc rl.Color) *Button {
	return &Button{
		Label:       l,
		Clicked:     false,
		Action:      a,
		Bounds:      b,
		BorderWidth: bw,
		BorderColor: bc,
		Color:       fc,
		Font:        f,
		FontSize:    fs,
		FontColor:   fontc,
	}
}

func (b *Button) Update() {
	mp := rl.GetMousePosition()

	if mp.X >= b.Bounds.X && mp.X <= b.Bounds.X+b.Bounds.Width && mp.Y >= b.Bounds.Y && mp.Y <= b.Bounds.Y+b.Bounds.Height && rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.Clicked == false {
		b.Action()
		b.Clicked = true
	} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) && b.Clicked == true {
		b.Clicked = false
	}

}

func (b *Button) IsHovered() bool {
	mp := rl.GetMousePosition()
	return mp.X >= b.Bounds.X && mp.X <= b.Bounds.X+b.Bounds.Width && mp.Y >= b.Bounds.Y && mp.Y <= b.Bounds.Y+b.Bounds.Height
}

func (b *Button) Draw() error {
	if b.IsHovered() && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		return b.DrawClick()
	}
	if b.IsHovered() {
		return b.DrawHover()
	}
	if b.BorderWidth > 0 {
		rl.DrawRectangle(int32(b.Bounds.X)-b.BorderWidth, int32(b.Bounds.Y)-b.BorderWidth, int32(b.Bounds.Width)+(b.BorderWidth*2), int32(b.Bounds.Y)+(b.BorderWidth*2), b.BorderColor)
	}

	lm := rl.MeasureTextEx(b.Font, b.Label, float32(b.FontSize), 0)

	rl.DrawText(b.Label, int32(b.Bounds.X)+(int32(b.Bounds.Width)/2)-int32(lm.X)/2, int32(b.Bounds.Y)+(int32(b.Bounds.Height/2))-int32(lm.Y)/2, 20, b.Color)
	return nil
}

func (b *Button) DrawHover() error {
	return nil
}

func (b *Button) DrawClick() error {
	return nil
}
