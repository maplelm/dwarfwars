package main

import (
	"context"
	"fmt"
	"time"

	rlgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
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
	connect bool
	quit    bool
}

func (mm *MainMenu) Init(g *game.Game) error      { return nil }
func (mm *MainMenu) IsInitialized() bool          { return true }
func (mm *MainMenu) UserInput(g *game.Game) error { return nil }
func (mm *MainMenu) Update(g *game.Game, cmds []*command.Command) error {

	if mm.connect {
		if g.IsConnected() {
			fmt.Println("Already connected to server")
		} else {
			fmt.Println("Connecting to server")
			connectctx, conncancel := context.WithDeadline(g.Ctx, time.Now().Add(time.Duration(5)*time.Second))
			go func(ctx context.Context) {
				defer conncancel()
				opts := g.Opts.MustGet()
				if err := g.Network(g.Ctx, opts.Network.Addr, opts.Network.Port); err != nil {
					fmt.Printf("Network failure, %s", err)
				}
			}(connectctx)
			select {
			case <-connectctx.Done():
				if connectctx.Err() != nil && !g.IsConnected() {
					fmt.Printf("Failed to connect to server, %s", connectctx.Err())
				} else {
					fmt.Printf("Successfully connected")
				}
			}
		}
		fmt.Println("going to connect")
	}
	if mm.quit {
		rl.CloseWindow()
	}

	return nil

}
func (mm *MainMenu) Draw() error {

	rl.BeginDrawing()
	mm.connect = rlgui.Button(rl.NewRectangle(0, 0, 100, 40), "Connect")
	mm.quit = rlgui.Button(rl.NewRectangle(0, 40, 100, 40), "Quit")

	return nil

}
