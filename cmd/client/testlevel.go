package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type testlevel struct {
	Sprites []rl.Texture2D
}

func (tl *testlevel) Init() error {
	return nil
}

func (tl *testlevel) UserInput() error {
	return nil
}

func (tl *testlevel) Update(nd [][]byte) error {
	return nil
}

func (tl *testlevel) Draw() error {
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText("Test Level", 190, 200, 20, rl.LightGray)
	return nil
}
