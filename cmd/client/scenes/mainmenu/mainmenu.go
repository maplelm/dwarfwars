package mainmenu

import (
	"fmt"
	"path/filepath"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui/button"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/pkg/command"
	"github.com/maplelm/dwarfwars/pkg/engine"

	"github.com/maplelm/dwarfwars/cmd/client/scenes/login"
	"github.com/maplelm/dwarfwars/cmd/client/scenes/signup"
)

type Scene struct {
	// Background
	bgcolor rl.Color

	// Interface
	Font     rl.Font
	FontSize int32

	// Interface Elements
	Menu button.List
	//Menu       *gui.ButtonList
	testbutton button.Button

	// State Tracking
	ScreenSize rl.Vector2
	init       bool

	// Testing
	testimagebutton button.ImageButton
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Init(g *game.Game) error {

	// Background Setup
	s.bgcolor = rl.DarkBrown

	// Interface Setup
	opts, err := g.Opts.Get()
	if err != nil {
		fmt.Println("Failed to Get Options to setup font: ", err)
		s.Font = rl.LoadFontEx(filepath.Join("./assets/fonts/", "Arial.ttf"), 400, nil, 0)
		s.FontSize = 32
	} else {
		s.Font = rl.LoadFontEx(filepath.Join("./assets/fonts/", opts.General.Font), opts.General.FontRes, nil, 0)
		s.FontSize = opts.General.FontSize
	}

	// Interface Element Setups
	a, err := engine.NewAnimationMatrix(3, 1, 3, 0, rl.LoadTexture("./assets/"), rl.Vector2{X: 32, Y: 32}, rl.White, nil)

	s.testimagebutton = *button.NewImageButton("Sign in", func() { g.PushScene(signup.New()) }, rl.NewRectangle(100, 100, 300, 100), *a, 0, rl.Black, rl.Font{}, 32, rl.Black)

	s.Menu = *button.NewList(rl.Vector2{X: float32(rl.GetScreenWidth()) / 2, Y: float32(rl.GetScreenHeight()) / 2}, s.Font, 2, 1, 2, rl.Black, rl.Green, s.FontSize, rl.Black, rl.Vector2{X: float32(rl.GetScreenWidth()) / 8, Y: float32(rl.GetScreenHeight()) / 6})
	s.Menu.Add("Connect", func() { Connect(g) })
	s.Menu.Add("Quit", func() { rl.CloseWindow() })
	s.Menu.Add("Sign Up", func() { g.PushScene(signup.New()) })
	s.Menu.Add("Login", func() { g.PushScene(login.New()) })

	// State Tracking Setup
	s.ScreenSize = rl.Vector2{
		X: float32(rl.GetScreenWidth()),
		Y: float32(rl.GetScreenHeight()),
	}

	// Connect to the Network
	if err = Connect(g); err != nil {
		fmt.Printf("Warning: Failed to connect to server, %s\n", err)
	}

	// Init Finished
	s.init = true
	return err
}

func (s *Scene) IsInitialized() bool { return s.init }

func (s *Scene) UserInput(g *game.Game) error { return nil }

func (s *Scene) Update(g *game.Game, cmds []*command.Command) error {
	s.testimagebutton.Update()

	s.Menu.Update()

	return nil
}

func (s *Scene) Draw() error {
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Brown)
	s.testimagebutton.Draw()
	s.Menu.Draw()
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
	s.Menu.ButtonSize.X = s.ScreenSize.X / 6
	s.Menu.ButtonSize.Y = s.ScreenSize.Y / 8
	s.Menu.Move(rl.Vector2{
		X: s.ScreenSize.X/2 - s.Menu.Position.X,
		Y: s.ScreenSize.Y/2 - s.Menu.Position.Y,
	})

	return nil
}

func (s *Scene) PausedUpdate(g *game.Game, cmds []*command.Command) error { return nil }
