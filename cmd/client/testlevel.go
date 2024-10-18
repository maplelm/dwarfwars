package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
)

type testlevel struct {
	Sprites []rl.Texture2D
}

func (tl *testlevel) Init(g *game.Game) error {
	return nil
}

func (tl *testlevel) UserInput(g *game.Game) error {
	if !g.Paused && rl.IsKeyPressed(rl.KeyEnter) && !rl.IsKeyPressedRepeat(rl.KeyEnter) {
		g.ActiveScene = 1
	}
	return nil
}

func (tl *testlevel) Update(g *game.Game, nd [][]byte) error {
	return nil
}

func (tl *testlevel) Draw() error {
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText("Test Level", 190, 200, 20, rl.LightGray)
	return nil
}

type testlevel2 struct {
	Sprites []rl.Texture2D
}

func (tl *testlevel2) Init(g *game.Game) error {
	return nil
}

func (tl *testlevel2) UserInput(g *game.Game) error {
	if rl.IsKeyPressed(rl.KeyEnter) && !rl.IsKeyPressedRepeat(rl.KeyEnter) {
		g.ActiveScene = 2
	}
	return nil
}

func (tl *testlevel2) Update(g *game.Game, nd [][]byte) error {
	return nil
}

func (tl *testlevel2) Draw() error {
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText("Test Level 2", 190, 200, 20, rl.LightGray)
	return nil
}
