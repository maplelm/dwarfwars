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
	fmt.Println("MainMenu Init")

	mm.Menu = gui.NewButtonList(rl.Vector2{X: 100, Y: 40}, rl.Vector2{X: 200, Y: 100}, 3, &g.Scale)
	mm.Menu.Add("Connect", func() {
		Connect(g)
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
	mm.Menu.Add("NULL", func() {})

	mm.Menu.Position = rl.Vector2{X: float32(rl.GetScreenWidth()) / 2, Y: float32(rl.GetScreenHeight()) / 2}
	mm.Menu.Center()

	// Connect to the Network
	err := Connect(g)
	if err != nil {
		fmt.Printf("Warning: Failed to connect to server, %s\n", err)
	}

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

///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////

func Connect(g *game.Game) error {
	if g.IsConnected() {
		return fmt.Errorf("already connected to server")
	}
	if g.IsConnecting() {
		return fmt.Errorf("already attempting to connected to server")
	}

	fmt.Println("Connecting to Game Server")

	opts, err := g.Opts.ForceRefresh()
	if err != nil {
		fmt.Printf("Warning, failed to refresh settings\n")
	}

	dl, _ := context.WithDeadline(g.NetworkCtx, time.Now().Add(opts.Network.ConnTimeout*time.Millisecond))
	attempt, cancelattempt := context.WithCancelCause(dl)

	done := make(chan struct{})
	defer close(done)

	go func(ctx context.Context, c chan struct{}) {
		opts, err := g.Opts.ForceRefresh()
		if err != nil {
			fmt.Printf("Failed to get settings, %s\n", err)
			cancelattempt(err)
			return
		}
		if err = g.Network(g.NetworkCtx, c, opts.Network.Addr, opts.Network.Port); err != nil {
			fmt.Printf("Network Failure, %s\n", err)
			cancelattempt(err)
		}
	}(attempt, done)

	select {
	case <-attempt.Done():
		if attempt.Err() != nil {
			fmt.Printf("Failed to connect to server, %s\n", attempt.Err())
			return attempt.Err()
		}
		fmt.Printf("Connection Attempt Timed out\n")
		return attempt.Err()
	case <-done:
		fmt.Printf("Successfully connected to server!\n")
		return nil
	}
}
