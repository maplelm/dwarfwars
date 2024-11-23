package button

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	// Properties
	Label   string
	Clicked bool
	Action  func()

	// Graphics
	/*
		Bounds      rl.Rectangle
		BorderWidth int32
		BorderColor rl.Color
		Color       rl.Color
		Font        rl.Font
		FontSize    int32
		FontColor   rl.Color
	*/
}

func New(l string, a func()) *Button {
	return &Button{
		Label:   l,
		Clicked: false,
		Action:  a,
		/*
			Bounds:      b,
			BorderWidth: bw,
			BorderColor: bc,
			Color:       fc,
			Font:        f,
			FontSize:    fs,
			FontColor:   fontc,
		*/
	}
}

func (b *Button) Update(bounds rl.Rectangle) {
	mp := rl.GetMousePosition()

	if mp.X >= bounds.X && mp.X <= bounds.X+bounds.Width && mp.Y >= bounds.Y && mp.Y <= bounds.Y+bounds.Height && rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.Clicked == false {
		b.Action()
		b.Clicked = true
	} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) && b.Clicked == true {
		b.Clicked = false
	}

}

func (b *Button) IsHovered(bounds rl.Rectangle) bool {
	mp := rl.GetMousePosition()
	return mp.X >= bounds.X && mp.X <= bounds.X+bounds.Width && mp.Y >= bounds.Y && mp.Y <= bounds.Y+bounds.Height
}

func (b *Button) Draw(bounds rl.Rectangle, border int32, bordercolor, color rl.Color, font rl.Font, fontsize int32, fontcolor rl.Color) error {
	if b.IsHovered(bounds) {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			return b.DrawClick(bounds, border, bordercolor, color, font, fontsize, fontcolor)
		} else {
			hcolor := color
			hcolor.R += 15
			hcolor.G += 15
			hcolor.B += 15
			return b.DrawHover(bounds, border, bordercolor, hcolor, font, fontsize, fontcolor)
		}
	}
	if border > 0 {
		rl.DrawRectangle(int32(bounds.X)-border, int32(bounds.Y)-border, int32(bounds.Width)+(border*2), int32(bounds.Y)+(border*2), bordercolor)
	}
	rl.DrawRectangle(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height), color)

	lm := rl.MeasureTextEx(font, b.Label, float32(fontsize), 0)

	rl.DrawTextEx(font, b.Label, rl.Vector2{X: bounds.X + bounds.Width/2 - lm.X/2, Y: bounds.Y + bounds.Height/2 - lm.Y/2}, float32(fontsize), 0, fontcolor)
	return nil
}

func (b *Button) DrawHover(bounds rl.Rectangle, border int32, bordercolor, color rl.Color, font rl.Font, fontsize int32, fontcolor rl.Color) error {
	if border > 0 {
		rl.DrawRectangle(int32(bounds.X)-border, int32(bounds.Y)-border, int32(bounds.Width)+(border*2), int32(bounds.Y)+(border*2), bordercolor)
	}
	rl.DrawRectangle(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height), color)

	lm := rl.MeasureTextEx(font, b.Label, float32(fontsize), 0)

	rl.DrawTextEx(font, b.Label, rl.Vector2{X: bounds.X + bounds.Width/2 - lm.X/2, Y: bounds.Y + bounds.Height/2 - lm.Y/2}, float32(fontsize), 0, fontcolor)
	return nil
}

func (b *Button) DrawClick(bounds rl.Rectangle, border int32, bordercolor, color rl.Color, font rl.Font, fontsize int32, fontcolor rl.Color) error {
	if border > 0 {
		rl.DrawRectangle(int32(bounds.X)-border, int32(bounds.Y)-border, int32(bounds.Width)+(border*2), int32(bounds.Y)+(border*2), bordercolor)
	}
	rl.DrawRectangle(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height), color)

	lm := rl.MeasureTextEx(font, b.Label, float32(fontsize), 0)

	rl.DrawTextEx(font, b.Label, rl.Vector2{X: bounds.X + bounds.Width/2 - lm.X/2, Y: bounds.Y + bounds.Height/2 - lm.Y/2}, float32(fontsize), 0, fontcolor)
	return nil
}
