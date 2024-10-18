package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Handler interface {
	Init(*Game) error
	UserInput(*Game) error
	Update(*Game, [][]byte) error // Might need to add some arguements here so that networked data can be passed here.
	Draw() error
}

type Game struct {
	Scenes      []Handler
	ActiveScene int

	Camera     rl.Camera2D
	ScreenSize rl.Vector2

	Paused bool
}

func New(screenx, screeny float32, title string, Scenes []Handler) *Game {
	rl.InitWindow(int32(screenx), int32(screeny), title)
	return &Game{
		Scenes:      Scenes,
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

	if rl.IsKeyPressed(rl.KeyP) && !rl.IsKeyPressedRepeat(rl.KeyP) {
		g.Paused = !g.Paused
	}

	if !g.Paused {
		err := g.Scenes[g.ActiveScene].UserInput(g)
		if err != nil {
			fmt.Printf("Warning Game Scene (%d) User Input Error: %s\n", g.ActiveScene, err)
		}
	}
}

func (g *Game) Network() {}

func (g *Game) Update() {

	if !g.Paused {
		err := g.Scenes[g.ActiveScene].Update(g, [][]byte{}) // will need to repace [][]byte{} with network traffic
		if err != nil {
			fmt.Printf("Warning Game Scene (%d) Draw Error: %s\n", g.ActiveScene, err)
		}
	}
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.RayWhite)
	rl.BeginDrawing()
	err := g.Scenes[g.ActiveScene].Draw()
	if err != nil {
		fmt.Printf("Warning Game Scene (%d) Draw Error: %s\n", g.ActiveScene, err)
	}

	if g.Paused {
		rl.DrawText("Paused", 0, 0, 12, rl.Black)
	} else {
		rl.DrawText("Unpaused", 0, 0, 12, rl.Black)
	}

	rl.EndDrawing()
}
