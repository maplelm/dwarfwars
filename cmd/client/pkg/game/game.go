package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Handler interface {
	Init() error
	UserInput() error
	Update([][]byte) error // Might need to add some arguements here so that networked data can be passed here.
	Draw() error
}

type Game struct {
	Scenes      []Handler
	ActiveScene int

	Camera     rl.Camera2D
	ScreenSize rl.Vector2
}

func New(screenx, screeny float32, title string) *Game {
	rl.InitWindow(int32(screenx), int32(screeny), title)
	return &Game{
		Scenes:      []Handler{},
		ActiveScene: 0,
		ScreenSize: rl.Vector2{
			X: screenx,
			Y: screeny,
		},
	}
}

func (g *Game) Run() {
	defer rl.CloseWindow()
	for !rl.WindowShouldClose() {
		g.Network()
		g.UserInput()
		g.Update()
		g.Draw()
	}
}
func (g *Game) UserInput() {
	err := g.Scenes[g.ActiveScene].UserInput()
	if err != nil {
		fmt.Printf("Warning Game Scene (%d) User Input Error: %s\n", g.ActiveScene, err)
	}
}

func (g *Game) Network() {}

func (g *Game) Update() {
	err := g.Scenes[g.ActiveScene].Update([][]byte{}) // will need to repace [][]byte{} with network traffic
	if err != nil {
		fmt.Printf("Warning Game Scene (%d) Draw Error: %s\n", g.ActiveScene, err)
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	err := g.Scenes[g.ActiveScene].Draw()
	if err != nil {
		fmt.Printf("Warning Game Scene (%d) Draw Error: %s\n", g.ActiveScene, err)
	}
	rl.EndDrawing()
}
