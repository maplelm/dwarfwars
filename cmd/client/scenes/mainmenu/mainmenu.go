package mainmenu

import (
	"fmt"
	"path/filepath"

	// rlgui "github.com/gen2brain/raylib-go/raygui"
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

	Font     rl.Font
	FontSize int32

	ScreenSize rl.Vector2

	bgcolor rl.Color

	init             bool
	Menu             *gui.ButtonList
	ButtonMenuOrigin rl.Rectangle
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Init(g *game.Game) error {
	var (
		MenuWidth int = 2
		err       error
		fontname  string
	)
	// Setting screen size so function does not need to be called every time the value is needed
	s.ScreenSize = rl.Vector2{
		X: float32(rl.GetScreenWidth()),
		Y: float32(rl.GetScreenHeight()),
	}

	// Loading Font into Memory
	opts, err := g.Opts.Get()
	if err != nil {
		fontname = "Arial.ttf"
		s.FontSize = 100
	} else {
		fontname = opts.General.Font
		s.FontSize = opts.General.FontSize
	}

	s.Font = rl.LoadFontEx(filepath.Join("./assets/fonts/", fontname), 400, nil, 0)

	// Setting Up Button Menu
	s.ButtonMenuOrigin = rl.NewRectangle(
		float32(rl.GetScreenWidth())/2.0,
		float32(rl.GetScreenHeight())/2.0,
		float32(rl.GetScreenWidth())/6,
		float32(rl.GetScreenHeight())/8,
	)
	s.Menu = gui.NewButtonList(s.Font, s.ButtonMenuOrigin, MenuWidth, &g.Scale)

	s.bgcolor = rl.DarkBrown

	s.Menu = gui.NewButtonList(s.Font, rl.NewRectangle(s.ScreenSize.X/2, s.ScreenSize.Y/2, 100, 40), MenuWidth, &g.Scale)
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
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Brown)

	sizing := rl.MeasureTextEx(s.Font, "Dwarf  Wars", 100, 0)
	rl.DrawTextEx(s.Font,
		"Dwarf  Wars",
		rl.Vector2{
			X: float32(rl.GetScreenWidth())/2.0 - sizing.X/2.0,
			Y: float32(int32(rl.GetScreenHeight())/10) - sizing.Y/2.0,
		},
		100,
		0,
		rl.Black,
	)
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

	// Re-Center button Cluster
	s.Menu.Buttonsize.X = s.ScreenSize.X / 6
	s.Menu.Buttonsize.Y = s.ScreenSize.Y / 8
	s.Menu.Position = rl.Vector2{
		X: s.ScreenSize.X / 2.0,
		Y: s.ScreenSize.Y / 2.0,
	}
	s.Menu.Center()

	return nil
}
