package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Scenes      []Scene
	ActiveScene int

	Camera rl.Camera2D

	Quiting bool
}

func (g *Game) Init() {}

func (g *Game) Run() {
	for !rl.WindowShouldClose() {
		g.Network()
		g.UserInput()
		g.Update()
		g.Draw()
	}
}
func (g *Game) UserInput() {}

func (g *Game) Network() {}

func (g *Game) Update() {}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.EndDrawing()
}
