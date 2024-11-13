package gui

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	Bounds  rl.Rectangle
	Clicked bool
	Label   string
	Action  func()
}
