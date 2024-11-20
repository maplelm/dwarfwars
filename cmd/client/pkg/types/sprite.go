package types

import rl "github.com/gen2brain/raylib-go/raylib"

type Sprite struct {
	texture    rl.Texture2D
	gridoffset rl.Vector2
	boxsize    rl.Vector2
	margin     Margin
}

type Margin struct {
	Top    float32
	Bottom float32
	Left   float32
	Right  float32
}
