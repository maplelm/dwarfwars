package login

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui/button"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type Scene struct {
	init bool
	menu *button.List
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Init(g *game.Game) error {
	defer func() { s.init = true }()

	s.menu = button.NewList(rl.Vector2{X: 100, Y: 100}, rl.Vector2{X: 150, Y: 50}, 2, 1, 32, rl.Black, rl.Blue, rl.Black, raygui.GetFont())
	s.menu.Add("Back", func() { g.PopScene() })
	s.menu.Add("Login", func() { fmt.Println("This is where you would login") })

	return nil
}

func (s *Scene) IsInitialized() bool { return s.init }

func (s *Scene) UserInput(g *game.Game) error {
	return nil
}

func (s *Scene) Update(g *game.Game, cmds []*command.Command) error {
	s.menu.Update(g.MP)
	return nil
}

func (s *Scene) Draw() error {
	rl.DrawText("Login Screen", 0, 0, 20, rl.Black)
	s.menu.Draw()
	return nil
}

func (s *Scene) Deconstruct() error {
	return nil
}

func (s *Scene) OnResize() error {
	return nil
}

func (s *Scene) PausedUpdate(g *game.Game, cmds []*command.Command) error { return nil }
