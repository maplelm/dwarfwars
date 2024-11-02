package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/pkg/command"
)

func TargetSprite(x, y float32) rl.Rectangle {
	return rl.Rectangle{
		X:      x * 15,
		Y:      y * 15,
		Width:  15,
		Height: 15,
	}
}

type Entity struct {
	Sprite   rl.Rectangle
	Size     rl.Vector2
	Pos      rl.Vector2
	Rotation float32
	Tint     rl.Color
}

type testlevel struct {
	SpriteSheet rl.Texture2D
	Entities    []Entity
	init        bool
	Camera      rl.Camera2D
}

func (tl *testlevel) Init(g *game.Game) error {

	tl.SpriteSheet = rl.LoadTexture("./assets/spritesheet.png")

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			tl.Entities = append(tl.Entities, Entity{Sprite: TargetSprite(float32(y%2), 0), Size: rl.Vector2{X: 32, Y: 32}, Pos: rl.Vector2{X: 32 * float32(x), Y: 32 * float32(y)}, Rotation: 0, Tint: rl.White})
		}
	}

	tl.Camera = rl.Camera2D{
		Offset:   rl.Vector2{X: 0, Y: 0},
		Target:   rl.Vector2{X: 0, Y: 0},
		Zoom:     1,
		Rotation: 0,
	}
	tl.init = true
	return nil
}

func (tl *testlevel) IsInitialized() bool {
	return tl.init
}

func (tl *testlevel) UserInput(g *game.Game) error {
	if !g.Paused && rl.IsKeyPressed(rl.KeyEnter) {
		return g.SetScene(1)
	}
	if !g.Paused && rl.IsKeyDown(rl.KeyUp) {
		tl.Camera.Target = rl.Vector2Add(tl.Camera.Target, rl.Vector2{X: 0, Y: -10})
	}
	if !g.Paused && rl.IsKeyDown(rl.KeyDown) {
		tl.Camera.Target = rl.Vector2Add(tl.Camera.Target, rl.Vector2{X: 0, Y: 10})
	}
	if !g.Paused && rl.IsKeyDown(rl.KeyRight) {
		tl.Camera.Target = rl.Vector2Add(tl.Camera.Target, rl.Vector2{X: 10, Y: 0})
	}
	if !g.Paused && rl.IsKeyDown(rl.KeyLeft) {
		tl.Camera.Target = rl.Vector2Add(tl.Camera.Target, rl.Vector2{X: -10, Y: 0})
	}
	if !g.Paused && rl.IsKeyPressed(rl.KeyS) {
		id := uint32(0)
		cmd, _ := command.New(id, 0, command.CommandType(1), []byte("Sent from client"))
		g.WriteQueue <- cmd
		fmt.Println("Command Queued")
	}
	return nil
}

func (tl *testlevel) Update(g *game.Game, cmds []*command.Command) error {
	return nil
}

func (tl *testlevel) Draw() error {
	rl.BeginMode2D(tl.Camera)
	for _, v := range tl.Entities {
		rl.DrawTexturePro(tl.SpriteSheet, v.Sprite, rl.Rectangle{X: v.Pos.X, Y: v.Pos.Y, Width: v.Size.X, Height: v.Size.Y}, rl.Vector2{X: 0, Y: 0}, v.Rotation, v.Tint)
	}
	rl.DrawText("Test Level", 190, 200, 20, rl.LightGray)
	rl.EndMode2D()
	rl.DrawText(fmt.Sprintf("Camera Pos: (%f,%f)", tl.Camera.Target.X, tl.Camera.Target.Y), 100, 10, 20, rl.Black)
	return nil
}

type testlevel2 struct {
	Sprites []rl.Texture2D
	init    bool
}

func (tl *testlevel2) Init(g *game.Game) error {
	tl.init = true
	return nil
}

func (tl *testlevel2) IsInitialized() bool {
	return tl.init
}

func (tl *testlevel2) UserInput(g *game.Game) error {
	if rl.IsKeyPressed(rl.KeyEnter) && !rl.IsKeyPressedRepeat(rl.KeyEnter) {
		return g.SetScene(2)
	}
	return nil
}

func (tl *testlevel2) Update(g *game.Game, cmds []*command.Command) error {
	return nil
}

func (tl *testlevel2) Draw() error {
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText("Test Level 2", 190, 200, 20, rl.LightGray)
	return nil
}
