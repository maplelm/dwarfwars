package button

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/maplelm/dwarfwars/pkg/engine"
)

type ImageButton struct {
	// Properties
	Label   string
	Clicked bool
	Action  func()

	// Graphics
	Bounds      rl.Rectangle
	Sprite      engine.AnimationMatrix
	BorderWidth int32
	BorderColor rl.Color
	Font        rl.Font
	FontSize    int32
	FontColor   rl.Color
}

func NewImageButton(l string, action func(), bounds rl.Rectangle, s engine.AnimationMatrix, bw int32, bc rl.Color, f rl.Font, fs int32, fc rl.Color) *ImageButton {
	return &ImageButton{
		Label:       l,
		Clicked:     false,
		Action:      action,
		Bounds:      bounds,
		Sprite:      s,
		BorderWidth: bw,
		BorderColor: bc,
		Font:        f,
		FontSize:    fs,
		FontColor:   fc,
	}
}

/*
I want to update this so that I can calculate if the button is being hovered even if there is a roation in the button. that is too much math for basic implementaiton though.
*/
func (b *ImageButton) IsHovered() bool {
	mp := rl.GetMousePosition()
	return mp.X >= b.Bounds.X && mp.X <= b.Bounds.X+b.Bounds.Width && mp.Y >= b.Bounds.Y && mp.Y <= b.Bounds.Y+b.Bounds.Height
}

func (b *ImageButton) Execute() error {
	return nil
}

func (b *ImageButton) Draw() error {
	return nil
}

func (b *ImageButton) DrawHover() error {
	return nil
}

func (b *ImageButton) DrawClicked() error {
	return nil
}
