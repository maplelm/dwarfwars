package mainmenu

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui"
	"github.com/maplelm/dwarfwars/pkg/command"

	"github.com/maplelm/dwarfwars/cmd/client/scenes/login"
	"github.com/maplelm/dwarfwars/cmd/client/scenes/signup"
)

type Scene struct {
	connect  bool
	quit     bool
	sendecho bool

	init bool
	Menu *gui.ButtonList
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Init(g *game.Game) error {
	var (
		MenuWidth int = 2
		err       error
	)

	s.Menu = gui.NewButtonList(rl.NewRectangle(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2, 100, 40), MenuWidth, &g.Scale)
	s.Menu.AddMulti([]gui.Button{
		gui.InitButton("Connect", func() { Connect(g) }),
		gui.InitButton("Quit", func() { rl.CloseWindow() }),
		gui.InitButton("Sign Up", func() { g.PushScene(signup.New()) }),
		gui.InitButton("Login", func() { g.PushScene(login.New()) }),
	})
	s.Menu.Center()

	// Connect to the Network
	if err = Connect(g); err != nil {
		fmt.Printf("Warning: Failed to connect to server, %s\n", err)
	}

	s.init = true
	return err
}

func (s *Scene) IsInitialized() bool { return s.init }

func (s *Scene) UserInput(g *game.Game) error { return nil }

func (s *Scene) Update(g *game.Game, cmds []*command.Command) error {

	if s.Menu != nil {
		s.Menu.Execute()
	} else {
		fmt.Println("s.Menu not initilized!")
	}
	return nil

}

func (s *Scene) Draw() error {
	s.Menu.Draw()
	return nil

}
