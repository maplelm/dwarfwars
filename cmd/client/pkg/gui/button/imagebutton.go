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

	// State
	hovered bool
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
	b.hovered = b.Bounds.X >= b.Bounds.X && mp.X <= b.Bounds.X+b.Bounds.Width && mp.Y >= b.Bounds.Y && mp.Y <= b.Bounds.Y+b.Bounds.Height
	return b.hovered
}

func (b *ImageButton) UpdateVisualState() {

	if !b.hovered && b.IsHovered() {
		b.Sprite.SetFrames(1)
	} else if b.hovered && !b.IsHovered() {
		b.Sprite.SetFrames(0)
	}

	if b.hovered && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		b.Sprite.SetFrames(2)
	}
	if b.Sprite.Frames() == 2 && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		if b.hovered {
			b.Sprite.SetFrames(1)
		} else {
			b.Sprite.SetFrames(0)
		}
	}
}

func (b *ImageButton) Update() {
	b.UpdateVisualState()
	if b.IsHovered() && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		b.Action()
	}
}

func (b *ImageButton) Draw() (err error) {
	b.Sprite.DrawAnimationFrame()
	return
}
