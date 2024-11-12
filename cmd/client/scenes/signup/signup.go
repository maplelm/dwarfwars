package signup

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type Scene struct {
	form          *gui.TextboxGroup
	Menu          gui.ButtonList
	minpasslength int
	back          bool
	initialized   bool
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Init(g *game.Game) error {
	defer func() { s.initialized = true }()
	fmt.Println("signin init")
	s.form = gui.NewTextboxGroup(rl.NewRectangle(100, 100, 200, 20), 255, 10, rl.Black)
	s.form.AddMulti([]gui.Textbox{
		gui.InitTextbox(s.form.Size, "Username", false, false, true, 0.9, 2),
		gui.InitTextbox(s.form.Size, "Password", true, false, false, 0.9, 2),
		gui.InitTextbox(s.form.Size, "Email", false, true, false, 0.9, 2),
	})
	s.form.Center()

	s.Menu = *gui.NewButtonList(rl.NewRectangle(400, 100, 100, 50), 2, nil)
	s.Menu.AddMulti([]gui.Button{
		gui.InitButton("Back", func() { g.PopScene() }),
		gui.InitButton("Password", func() {
			val, _ := s.form.ValueByLabel("Password")
			fmt.Println("password data: ", val)
		}),
		gui.InitButton("Username", func() {
			val, _ := s.form.ValueByLabel("Username")
			fmt.Println("Username data: ", val)
		}),
		gui.InitButton("Email", func() {
			val, _ := s.form.ValueByLabel("Email")
			fmt.Println("Email data: ", val)
		}),
	})

	return nil
}
func (s *Scene) IsInitialized() bool          { return s.initialized }
func (s *Scene) UserInput(g *game.Game) error { return nil }
func (s *Scene) Update(g *game.Game, cmds []*command.Command) error {
	s.Menu.Execute()
	return nil
}
func (s *Scene) Draw() error {

	s.form.Draw()
	s.Menu.Draw()
	return nil
}

func (s *Scene) Deconstruct() error {
	return nil
}

func (s *Scene) OnResize() error {
	return nil
}
