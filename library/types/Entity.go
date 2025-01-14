package types

import rl "github.com/gen2brain/raylib-go/raylib"

type Entity struct {
	Position          float64 // (tile Index).(Inner Tile Position) (000).(000)(000)(000) (XYZ).(X,Y,Z)
	AnimationPosition rl.Vector2
}
