package main

import (
	"context"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui"
	"github.com/maplelm/dwarfwars/pkg/command"
)

func TargetSprite(x, y float32) rl.Rectangle {
	return rl.Rectangle{
		X:      x * 15,
		Y:      y * 15,
		Width:  15,
		Height: 15,
	}
}

type MainMenu struct {
	connect  bool
	quit     bool
	sendecho bool

	init bool
	Menu *gui.ButtonList
}

func (mm *MainMenu) Init(g *game.Game) error {
	defer func() { mm.init = true }()
	fmt.Println("\n\nMainMenu Init\n\n")
	mm.Menu = gui.NewButtonList(rl.Vector2{X: 100, Y: 40}, rl.Vector2{X: 200, Y: 100}, 5, &g.Scale)
	mm.Menu.Add("Connect", func() {
		if g.IsConnected() {
			fmt.Println("Already connected to server")
		} else if g.IsConnecting() {
			fmt.Println("Still Establishing Connection to server")
		} else {
			fmt.Println("Connecting to server")
			connectionAttempt, AttemptCancel := context.WithDeadline(g.NetworkCtx, time.Now().Add(time.Duration(5)*time.Second))
			success := make(chan struct{})
			defer AttemptCancel()
			go func(ctx context.Context, success chan struct{}) {
				opts := g.Opts.MustGet()
				if err := g.Network(g.NetworkCtx, success, opts.Network.Addr, opts.Network.Port); err != nil {
					fmt.Printf("Network failure, %s", err)
				}
			}(connectionAttempt, success)
			select {
			case <-connectionAttempt.Done():
				if connectionAttempt.Err() != nil {
					fmt.Printf("Failed to connect to server, %s\n", connectionAttempt.Err())
				} else {
					fmt.Printf("Connection Attempt Timed out")
				}
			case <-success:
				fmt.Println("Successfully connected to server!")
			}
		}
	})
	mm.Menu.Add("Quit", func() {
		rl.CloseWindow()
	})
	mm.Menu.Add("Echo", func() {
		if g.IsConnected() {
			cmd, _ := command.New(g.ServerID, command.FormatText, command.TypeEcho, []byte("Ping!"))
			g.WriteQueue <- cmd
		} else if !g.IsConnected() {
			fmt.Println("Can't send message, not connected to server")
		}
	})
	mm.Menu.Add("Sign Up", func() {
		g.PushScene(&Signup{})
	})
	mm.Menu.Add("Login", func() {
		fmt.Println("Login button pressed")
	})

	mm.Menu.Position = rl.Vector2{X: float32(rl.GetScreenWidth()) / 2, Y: float32(rl.GetScreenHeight()) / 2}
	mm.Menu.Position = mm.Menu.Centered()

	return nil
}
func (mm *MainMenu) IsInitialized() bool          { return mm.init }
func (mm *MainMenu) UserInput(g *game.Game) error { return nil }
func (mm *MainMenu) Update(g *game.Game, cmds []*command.Command) error {

	if mm.Menu != nil {
		mm.Menu.Execute()
	} else {
		fmt.Println("mm.Menu not initilized!")
	}
	return nil

}
func (mm *MainMenu) Draw() error {
	mm.Menu.Draw()
	return nil

}
