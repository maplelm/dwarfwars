package server

import (
	"context"
	"errors"
	"log"
	"net"
	"syscall"

	"github.com/maplelm/dwarfwars/cmd/server/pkg/client"
)

func (s *Server) connMgr(logger *log.Logger, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		// Wait for connection to come through from the listener
		case conn := <-s.CC:
			// Create a Client object to use with the new connection
			newClient, _, err := s.clients.Connect(&conn)
			if err != nil {
				logger.Printf("Error Creating Client from Connection: %s", err)
				conn.Close()
				continue
			}
			go s.handleConnection(newClient, logger, ctx)
		}
	}
}

func (s *Server) handleConnection(c *client.Client, logger *log.Logger, ctx context.Context) error {
	err := c.Serve(logger, ctx, c.Connection)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			// Client has closed the connection properly
		}
		if errors.Is(err, syscall.ECONNRESET) {
			// connection has most likely dropped
		}
	}

	// Disconnect the client and remove the client object
	s.clients.Disconnect(c.GetID())

	return nil
}
