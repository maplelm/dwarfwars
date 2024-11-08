package main

import (
	"context"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type MainMenu struct {
	connect  bool
	quit     bool
	sendecho bool

	init bool
	Menu *gui.ButtonList
}

func (mm *MainMenu) Init(g *game.Game) error {
	var (
		MenuWidth int = 2
		err       error
	)

	mm.Menu = gui.NewButtonList(rl.NewRectangle(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2, 100, 40), MenuWidth, &g.Scale)
	mm.Menu.AddMulti([]gui.Button{
		gui.InitButton("Connect", func() { Connect(g) }),
		gui.InitButton("Quit", func() { rl.CloseWindow() }),
		gui.InitButton("Sign Up", func() { g.PushScene(&Signup{}) }),
		gui.InitButton("Login", func() { g.PushScene(&Login{}) }),
	})
	mm.Menu.Center()

	// Connect to the Network
	if err = Connect(g); err != nil {
		fmt.Printf("Warning: Failed to connect to server, %s\n", err)
	}

	mm.init = true
	return err
}

func (mm *MainMenu) IsInitialized() bool { return mm.init }

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
	var (
		done          chan struct{} = make(chan struct{})
		opts          *types.Options
		err           error
		dl            context.Context
		dlc           func()
		attempt       context.Context
		cancelattempt func(error)
	)
	fmt.Println("Connecting to Game Server")

	if opts, err = g.Opts.ForceRefresh(); err != nil {
		fmt.Printf("Warning, failed to refresh settings\n")
	}

	// Setting up contexts
	dl, dlc = context.WithDeadline(g.NetworkCtx, time.Now().Add(opts.Network.ConnTimeout*time.Millisecond))
	defer dlc()
	attempt, cancelattempt = context.WithCancelCause(dl)

	// Running the networking code
	go func() {
		err = g.Network(g.NetworkCtx, opts.Network.Addr, opts.Network.Port)
		if err != nil {
			cancelattempt(err)
			return
		}
		close(done)
	}()

	select {
	case <-attempt.Done():
		if attempt.Err() != nil {
			fmt.Printf("Failed to connect to server, %s\n", attempt.Err())
		} else {
			fmt.Printf("Connection Attempt Timed out\n")
		}
		return attempt.Err()
	case _, ok := <-done:
		fmt.Printf("Successfully connected to server!\n")
		if ok {
			close(done)
		}
		return nil
	}
}
