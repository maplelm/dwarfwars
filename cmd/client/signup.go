package main

import (
	rlgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type Signup struct {
	form          *gui.TextBoxGroup
	Menu          gui.ButtonList
	minpasslength int
	back          bool
	initialized   bool
}

func (s *Signup) Init(g *game.Game) error {
	defer func() { s.initialized = true }()
	s.form = gui.NewTextBoxGroup(rl.NewRectangle(100, 100, 200, 50), 255, 10, rl.Black)
	s.form.AddMulti([]gui.Textbox{
		gui.InitTextbox(s.form.Size, "Username", false, false, false, 0.9, 2),
		gui.InitTextbox(s.form.Size, "Password", true, false, false, 0.9, 2),
		gui.InitTextbox(s.form.Size, "Email", false, true, false, 0.9, 2),
	})
	s.form.Center()

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

	s.form.Draw()
	/*
		if rlgui.TextBox(rl.NewRectangle(float32(rl.GetScreenWidth())/2.0 - (255.0/2.0), float32(rl.GetScreenHeight()/2 - (40.0/2.0) + , 255, 40), &s.Username, 255, s.usernameactivestate) {
			s.usernameactivestate = !s.usernameactivestate
		}
		if rlgui.TextBox(rl.NewRectangle(0, 80, 255, 40), &s.Password, 255, s.passwordactivestate) {
			s.passwordactivestate = !s.passwordactivestate
		}
		if rlgui.TextBox(rl.NewRectangle(0, 120, 255, 40), &s.Email, 255, s.emailactivestate) {
			s.emailactivestate = !s.emailactivestate
		}
		rl.DrawText("Username", 0, 0, 12, rl.Black)
	*/
	s.back = rlgui.Button(rl.NewRectangle(400, 100, 100, 40), "Back")
	return nil
}
