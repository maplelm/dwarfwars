package engine

import rl "github.com/gen2brain/raylib-go/raylib"

type Entity struct {
	Sprite   rl.Rectangle
	Size     rl.Vector2
	Pos      rl.Vector2
	Rotation float32
	Tint     rl.Color
}
