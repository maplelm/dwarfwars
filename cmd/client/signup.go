package main

import (
	rlgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type Signup struct {
	Username string
	Password string
	Email    string
	SSO      string
	Menu     gui.ButtonList

	back bool

	initialized bool
}

func (s *Signup) Init(g *game.Game) error {
	defer func() { s.initialized = true }()
	return nil
}
func (s *Signup) IsInitialized() bool          { return s.initialized }
func (s *Signup) UserInput(g *game.Game) error { return nil }
func (s *Signup) Update(g *game.Game, cmds []*command.Command) error {
	if s.back {
		g.PopScene()
	}
	return nil
}
func (s *Signup) Draw() error {

	rlgui.TextBox(rl.NewRectangle(0, 40, 255, 40), &s.Username, 12, true)
	rl.DrawText("Username", 0, 0, 12, rl.Black)
	s.back = rlgui.Button(rl.NewRectangle(100, 100, 100, 40), "Back")
	return nil
}
