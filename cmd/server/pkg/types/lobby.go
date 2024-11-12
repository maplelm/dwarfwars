package types

import (
	"context"
	"golang.org/x/time/rate"
	"time"

	"github.com/maplelm/dwarfwars/cmd/server/pkg/client"
	//"github.com/maplelm/dwarfwars/cmd/server/pkg/player"
	"github.com/maplelm/dwarfwars/pkg/command"
	"github.com/maplelm/dwarfwars/pkg/engine"
	"github.com/maplelm/dwarfwars/pkg/types"
)

type Lobby struct {
	players      map[uint32]<-chan *command.Command // have a channel of inputs per player connected to lobby?
	LobbyLeader  uint32
	max          int
	mode         types.LobbyMode
	password     *string
	world        *engine.World
	cmdQueue     chan *command.Command
	PauseOnEmpty bool
	tickRate     int
}

func NewLobby(mode types.LobbyMode, tickrate, mp int, pass *string, poe bool) *Lobby {
	return &Lobby{
		players:      make(map[uint32]<-chan *command.Command),
		max:          mp,
		mode:         mode,
		password:     pass,
		world:        nil,
		cmdQueue:     make(chan *command.Command, 100),
		tickRate:     tickrate,
		PauseOnEmpty: poe,
	}

}

func (l *Lobby) ConnectPlayer(c *client.Client) error {
	return nil
}

func (l *Lobby) PreGame(ctx context.Context, x, y, z, seed int) error {
	rlim := rate.NewLimiter(rate.Every(time.Second/time.Duration(l.tickRate)), 1)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(rlim.Reserve().Delay()):
			if len(l.cmdQueue) == 0 {
				continue
			}
			cmds := make([]*command.Command, len(l.cmdQueue))
			length := len(l.cmdQueue)
			for i := 0; i < length; i++ {
				cmds[i] = <-l.cmdQueue
			}
			for _, v := range cmds[:length] {
				switch v.Type {
				case command.TypeInput:
				case command.TypeStartGame:
				default:
				}
			}
		default:
		}
	}
}

func (l *Lobby) Start(ctx context.Context) error {
	rlim := rate.NewLimiter(rate.Every(time.Second/time.Duration(l.tickRate)), 1)
	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(rlim.Reserve().Delay()):
			// pause game if lobby is empty
			if len(l.players) != 0 || !l.PauseOnEmpty {
				l.world.Update()
			}
		}

	}
}
