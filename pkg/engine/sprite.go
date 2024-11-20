package engine

import rl "github.com/gen2brain/raylib-go/raylib"

type Sprite struct {
	Texture     rl.Texture2D
	Spritesheet *SpriteSheet
	Animation   *AnimationMatrix
}

func NewSprite(t rl.Texture2D, ss *SpriteSheet, am *AnimationMatrix) *Sprite {
	return &Sprite{
		Texture:     t,
		Spritesheet: ss,
		Animation:   am,
	}
}

func (s *Sprite) PushFrame() rl.Rectangle {
	x, y := 1, 1 // s.Animation.AnimationFrame()
	return s.Spritesheet.GetFrame(x, y)
}
