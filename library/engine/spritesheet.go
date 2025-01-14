packagesrc/server/pkg engine

import rl "github.com/gen2brain/raylib-go/raylib"

type SpriteSheet struct {
	OffsetX float32
	OffsetY float32

	Width  float32
	Height float32

	MarginTop    float32
	MarginBottom float32
	MarginLeft   float32
	MarginRight  float32
}

func (ss *SpriteSheet) GetFrame(x, y int) rl.Rectangle {
	return rl.Rectangle{
		X:      ss.OffsetX + ((float32(x) + 1) * ss.MarginLeft) + (float32(x) * ss.MarginRight) + (float32(x) * ss.Width),
		Y:      ss.OffsetY + ((float32(y) + 1) * ss.MarginTop) + (float32(y) * ss.MarginBottom) + (float32(y) * ss.Height),
		Width:  ss.Width,
		Height: ss.Height,
	}
}
