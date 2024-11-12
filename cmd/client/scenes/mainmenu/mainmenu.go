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

	ScreenSize rl.Vector2

	Font     rl.Font
	FontSize int32

	bgcolor rl.Color

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
	// Setting screen size so function does not need to be called every time the value is needed
	s.ScreenSize = rl.Vector2{
		X: float32(rl.GetScreenWidth()),
		Y: float32(rl.GetScreenHeight()),
	}

	s.FontSize = g.Opts.MustGet().General.FontSize

	s.bgcolor = rl.DarkBrown

	s.Menu = gui.NewButtonList(rl.NewRectangle(s.ScreenSize.X/2, s.ScreenSize.Y/2, 100, 40), MenuWidth, &g.Scale)
	s.Menu.AddMulti([]gui.Button{
		gui.InitButton("Connect", func() { Connect(g) }),
		gui.InitButton("Quit", func() { rl.CloseWindow() }),
		gui.InitButton("Sign Up", func() { g.PushScene(signup.New()) }),
		gui.InitButton("Login", func() { g.PushScene(login.New()) }),
	})
	s.Menu.Buttonsize = rl.Vector2{
		X: s.ScreenSize.X / 8,
		Y: s.ScreenSize.Y / 6,
	}
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
	s.Menu.Execute()
	return nil
}

func (s *Scene) Draw() error {
	rl.DrawRectangle(0, 0, int32(s.ScreenSize.X), int32(s.ScreenSize.Y), s.bgcolor)

	//rl.DrawTextEx()
	//rl.DrawText("Dwarf Wars", 100, 100, 32, rl.Black)
	rl.DrawText("Dwarf Wars", int32(s.ScreenSize.X/2.0-float32(rl.MeasureText("Dwarf Wars", 200))/2), int32(s.ScreenSize.Y/6), 200, rl.Black)
	s.Menu.Draw()
	return nil

}

func (s *Scene) Deconstruct() error {
	return nil
}

func (s *Scene) OnResize() error {
	s.ScreenSize = rl.Vector2{
		X: float32(rl.GetScreenWidth()),
		Y: float32(rl.GetScreenHeight()),
	}

	s.Menu.Buttonsize = rl.Vector2{
		X: s.ScreenSize.X / 8,
		Y: s.ScreenSize.Y / 6,
	}

	s.Menu.Position = rl.Vector2{
		X: s.ScreenSize.X / 2,
		Y: s.ScreenSize.Y / 2,
	}
	s.Menu.Center()
	return nil
}
