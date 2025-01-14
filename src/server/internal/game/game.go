package game

import (
	"context"
	"fmt"
	"log"
	"time"

	"library/command"
	"library/engine"
	"library/types"
	"server/internal/client"
)

type Game struct {
	players    []*client.Client
	MaxPlayers int
	Mode       types.LobbyMode
	password   string
	World      *engine.World
	TickRate   time.Duration
	stop       chan struct{}
}

func New(mode types.LobbyMode, pass string, maxp, x, y, z, seed int, tr time.Duration) *Game {
	return &Game{
		players:    make([]*client.Client, maxp),
		MaxPlayers: maxp,
		Mode:       mode,
		password:   pass,
		TickRate:   tr,
		World:      engine.WorldGen(x, y, z, seed),
	}
}

func (g *Game) Run(ctx context.Context, logger *log.Logger) error {
	timer := time.NewTimer(g.TickRate)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			timer.Reset(g.TickRate - g.loop(logger))
		}
	}
}

/*
This will be the main game loop but right now there is just filler code sending a command to the client to keep the connection from timing out.
*/
func (g *Game) loop(logger *log.Logger) (delta time.Duration) {

	s := time.Now()

	// Main Loop //

	for _, c := range g.players {
		cmd, err := command.New(c.Uid(), command.FormatText, command.TypeKeepAlive, []byte("ping from server"))
		if err != nil {
			logger.Printf("Error Sending Ping to Client (%d), %s", c.Uid(), err)
			continue
		}
		if n := c.Send(cmd); n == 0 {
			logger.Printf("Warning: No Bytes where sent to client (%d)", c.Uid())
		}
	}

	////////////

	delta = time.Since(s)
	return
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
