package mainmenu

import (
	"context"
	"fmt"
	"time"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/types"
)

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
