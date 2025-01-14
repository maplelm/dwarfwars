package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Entity struct {

	// Render components
	CurrentSprite   rl.Rectangle
	CurrentSize     rl.Vector2
	CurrentPos      rl.Vector2
	CurrentRotation float32
	CurrentTint     rl.Color

	AnimationSprite   rl.Rectangle
	AnimationSize     rl.Vector2
	AnimationPos      rl.Vector2
	AnimationRotation float32
	AnimationTint     rl.Color
}
