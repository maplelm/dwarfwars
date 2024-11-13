package login

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type Scene struct {
	init bool
	menu *gui.ButtonList
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Init(g *game.Game) error {
	defer func() { s.init = true }()

	s.menu = gui.NewButtonList(raygui.GetFont(), rl.NewRectangle(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2, 100, 40), 3, &g.Scale)
	s.menu.Add("Back", func() { g.PopScene() })
	s.menu.Add("Login", func() { fmt.Println("This is where you would login") })
	s.menu.Center()

	return nil
}

func (s *Scene) IsInitialized() bool { return s.init }

func (s *Scene) UserInput(g *game.Game) error {
	return nil
}

func (s *Scene) Update(g *game.Game, cmds []*command.Command) error {
	s.menu.Execute()
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
