package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/gui"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type Login struct {
	init bool
	menu *gui.ButtonList
}

func (l *Login) Init(g *game.Game) error {
	defer func() { l.init = true }()

	l.menu = gui.NewButtonList(rl.NewRectangle(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2, 100, 40), 3, &g.Scale)
	l.menu.Add("Back", func() { g.PopScene() })
	l.menu.Add("Login", func() { fmt.Println("This is where you would login") })
	l.menu.Center()

	return nil
}

func (l *Login) IsInitialized() bool { return l.init }

func (l *Login) UserInput(g *game.Game) error {
	return nil
}

func (l *Login) Update(g *game.Game, cmds []*command.Command) error {
	l.menu.Execute()
	return nil
}

func (l *Login) Draw() error {
	rl.DrawText("Login Screen", 0, 0, 20, rl.Black)
	l.menu.Draw()
	return nil
}
