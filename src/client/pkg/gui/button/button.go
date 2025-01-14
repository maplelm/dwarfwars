package button

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/src/client/pkg/gui"
)

type bGraphics struct {
	Bounds      rl.Rectangle
	BorderWidth int32
	BorderColor rl.Color
	Color       rl.Color
	Font        rl.Font
	FontSize    int32
	FontColor   rl.Color
}

func NewbGraphics(bounds rl.Rectangle, borderwidth int32, bordercolor, color rl.Color, font rl.Font, fontsize int32, fontcolor rl.Color) *bGraphics {
	return &bGraphics{
		Bounds:      bounds,
		BorderWidth: borderwidth,
		BorderColor: bordercolor,
		Color:       color,
		Font:        font,
		FontSize:    fontsize,
		FontColor:   fontcolor,
	}
}

type Button struct {
	Label   string
	Clicked bool
	Action  func()

	Graphics *bGraphics
}

func New(l string, a func(), graphics *bGraphics) *Button {
	return &Button{
		Label:    l,
		Clicked:  false,
		Action:   a,
		Graphics: graphics,
	}
}

func (b *Button) Update(mp rl.Vector2) error {
	if b.Graphics == nil {
		return gui.ErrNoGraphics("button update failed, must use UpdateWithBounds if Graphics is nil")
	}

	if b.MustIsHovered(mp) && rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.Clicked == false {
		b.Action()
		b.Clicked = true
	} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) && b.Clicked == true {
		b.Clicked = false
	}
	return nil
}

func (b *Button) UpdateWithBounds(mp rl.Vector2, bounds rl.Rectangle) {
	if b.IsHoveredWithBounds(mp, bounds) && rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.Clicked == false {
		b.Action()
		b.Clicked = true
	} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) && b.Clicked == true {
		b.Clicked = false
	}
}

func (b *Button) IsHovered(mp rl.Vector2) (bool, error) {
	if b.Graphics == nil {
		return false, gui.ErrNoGraphics("button IsHovered failed, must use IsHoveredWithBounds if Graphics is nil")
	}
	return b.IsHoveredWithBounds(mp, b.Graphics.Bounds), nil
}

func (b *Button) MustIsHovered(mp rl.Vector2) bool {
	if b.Graphics == nil {
		panic(gui.ErrNoGraphics("button IsHovered failed, must use IsHoveredWithBounds if Graphics is nil"))
	}
	return b.IsHoveredWithBounds(mp, b.Graphics.Bounds)
}

func (b *Button) IsHoveredWithBounds(mp rl.Vector2, bounds rl.Rectangle) bool {
	return mp.X > bounds.X && mp.X < bounds.X+bounds.Width && mp.Y > bounds.Y && mp.Y < bounds.Y+bounds.Height
}

func (b *Button) Draw(mp rl.Vector2) error {
	if b.Graphics == nil {
		return gui.ErrNoGraphics("button draw failed, must use DrawWithGraphics if Grphics is null")
	}
	if b.MustIsHovered(mp) {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			// Button Clicked on
			g := *b.Graphics
			g.Color.R += 15
			g.Color.G += 15
			g.Color.B += 15
			b.DrawWithGraphics(g)
		} else {
			// Button Hovered
			g := *b.Graphics
			g.Color.R -= 15
			g.Color.G -= 15
			g.Color.B -= 15
			b.DrawWithGraphics(g)
		}
		return nil
	}
	b.DrawWithGraphics(*b.Graphics)
	return nil
}

func (b *Button) DrawCustom(
	bounds rl.Rectangle,
	border int32,
	bordercolor, color rl.Color,
	font rl.Font,
	fontsize int32,
	fontcolor rl.Color,
) {
	b.DrawWithGraphics(bGraphics{
		Bounds:      bounds,
		BorderWidth: border,
		BorderColor: bordercolor,
		Color:       color,
		Font:        font,
		FontSize:    fontsize,
		FontColor:   fontcolor,
	})
}

func (b *Button) DrawWithGraphics(g bGraphics) {

	// Only Draw Border if it will showon the button
	if g.BorderWidth > 0 {
		x := int32(g.Bounds.X - float32(g.BorderWidth))
		y := int32(g.Bounds.Y - float32(g.BorderWidth))
		w := int32(g.Bounds.Width + float32(g.BorderWidth*2))
		h := int32(g.Bounds.Height + float32(g.BorderWidth*2))
		rl.DrawRectangle(x, y, w, h, g.BorderColor)
	}
	rl.DrawRectangle(int32(g.Bounds.X), int32(g.Bounds.Y), int32(g.Bounds.Width), int32(g.Bounds.Height), g.Color)

	lm := rl.MeasureTextEx(g.Font, b.Label, float32(g.FontSize), 0).X / 2

	pos := rl.Vector2{
		X: (g.Bounds.X + g.Bounds.Width/2) - lm,
		Y: (g.Bounds.Y + g.Bounds.Height/2) - float32(g.FontSize/2),
	}

	rl.DrawTextEx(g.Font, b.Label, pos, float32(g.FontSize), 0, g.FontColor)
}
