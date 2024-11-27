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

func New(l string, a func(), b rl.Rectangle, bw int32, bc rl.Color, c rl.Color, f rl.Font, fs int32, fc rl.Color) *Button {
	return &Button{
		Label:       l,
		Clicked:     false,
		Action:      a,
		Bounds:      b,
		BorderWidth: bw,
		BorderColor: bc,
		Color:       c,
		Font:        f,
		FontSize:    fs,
		FontColor:   fc,
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

func (b *Button) UpdateWithBounds(bounds rl.Rectangle) {
	if b.IsHoveredWithBounds(bounds) && rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.Clicked == false {
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

func (b *Button) IsHoveredWithBounds(bounds rl.Rectangle) bool {
	mp := rl.GetMousePosition()
	return mp.X >= bounds.X && mp.X <= bounds.X+bounds.Width && mp.Y >= bounds.Y && mp.Y <= bounds.Y+bounds.Height
}

func (b *Button) Draw() error {
	if b.IsHovered() {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			return b.DrawClick()
		} else {
			return b.DrawHover()
		}
	}
	if b.BorderWidth > 0 {
		rl.DrawRectangle(int32(b.Bounds.X)-b.BorderWidth, int32(b.Bounds.Y)-b.BorderWidth, int32(b.Bounds.Width)+(b.BorderWidth*2), int32(b.Bounds.Y)+(b.BorderWidth*2), b.BorderColor)
	}
	rl.DrawRectangle(int32(b.Bounds.X), int32(b.Bounds.Y), int32(b.Bounds.Width), int32(b.Bounds.Height), b.Color)

	lm := rl.MeasureTextEx(b.Font, b.Label, float32(b.FontSize), 0)

	rl.DrawTextEx(b.Font, b.Label, rl.Vector2{X: b.Bounds.X + b.Bounds.Width/2 - lm.X/2, Y: b.Bounds.Y + b.Bounds.Height/2 - lm.Y/2}, float32(b.FontSize), 0, b.FontColor)
	return nil
}

func (b *Button) DrawHover() error {
	hcolor := b.Color
	hcolor.R += 15
	hcolor.G += 15
	hcolor.B += 15
	if b.BorderWidth > 0 {
		rl.DrawRectangle(int32(b.Bounds.X)-b.BorderWidth, int32(b.Bounds.Y)-b.BorderWidth, int32(b.Bounds.Width)+(b.BorderWidth*2), int32(b.Bounds.Y)+(b.BorderWidth*2), b.BorderColor)
	}
	rl.DrawRectangle(int32(b.Bounds.X), int32(b.Bounds.Y), int32(b.Bounds.Width), int32(b.Bounds.Height), hcolor)

	lm := rl.MeasureTextEx(b.Font, b.Label, float32(b.FontSize), 0)

	rl.DrawTextEx(b.Font, b.Label, rl.Vector2{X: b.Bounds.X + b.Bounds.Width/2 - lm.X/2, Y: b.Bounds.Y + b.Bounds.Height/2 - lm.Y/2}, float32(b.FontSize), 0, b.FontColor)
	return nil
}

func (b *Button) DrawClick() error {
	if b.BorderWidth > 0 {
		rl.DrawRectangle(int32(b.Bounds.X)-b.BorderWidth, int32(b.Bounds.Y)-b.BorderWidth, int32(b.Bounds.Width)+(b.BorderWidth*2), int32(b.Bounds.Y)+(b.BorderWidth*2), b.BorderColor)
	}
	rl.DrawRectangle(int32(b.Bounds.X), int32(b.Bounds.Y), int32(b.Bounds.Width), int32(b.Bounds.Height), b.Color)

	lm := rl.MeasureTextEx(b.Font, b.Label, float32(b.FontSize), 0)

	rl.DrawTextEx(b.Font, b.Label, rl.Vector2{X: b.Bounds.X + b.Bounds.Width/2 - lm.X/2, Y: b.Bounds.Y + b.Bounds.Height/2 - lm.Y/2}, float32(b.FontSize), 0, b.FontColor)
	return nil
}

func (b *Button) DrawWithGraphics(
	bounds rl.Rectangle,
	border int32,
	bordercolor, color rl.Color,
	font rl.Font,
	fontsize int32,
	fontcolor rl.Color,
) error {

	if border > 0 {
		rl.DrawRectangle(int32(bounds.X)-border, int32(bounds.Y)-border, int32(bounds.Width)+(border*2), int32(bounds.Y)+(border*2), bordercolor)
	}
	rl.DrawRectangle(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height), color)

	lm := rl.MeasureTextEx(font, b.Label, float32(fontsize), 0)

	rl.DrawTextEx(font, b.Label, rl.Vector2{X: bounds.X + (bounds.Width/2 - lm.X/2), Y: bounds.Y + (bounds.Height/2 - lm.Y/2)}, float32(fontsize), 0, fontcolor)
	return nil
}
