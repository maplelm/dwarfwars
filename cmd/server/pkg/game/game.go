package game

import (
	"fmt"

	"github.com/maplelm/dwarfwars/cmd/server/pkg/client"
	"github.com/maplelm/dwarfwars/pkg/engine"
	"github.com/maplelm/dwarfwars/pkg/types"
)

type Game struct {
	players    []*client.Client
	MaxPlayers int
	Mode       types.LobbyMode
	password   string
	World      *engine.World
}

func New(mode types.LobbyMode, pass string, maxp, x, y, z, seed int) *Game {
	return &Game{
		players:    make([]*client.Client, maxp),
		MaxPlayers: maxp,
		Mode:       mode,
		password:   pass,
		World:      engine.WorldGen(x, y, z, seed),
	}
}

func (g *Game) IsFull() bool {
	return len(g.players) >= g.MaxPlayers
}

func (g *Game) Connect(c *client.Client) error {
	if g.IsFull() {
		return fmt.Errorf("game full")
	}
	g.players[len(g.players)] = c
	return nil
}
