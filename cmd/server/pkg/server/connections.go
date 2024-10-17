package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

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
			id, index, err := s.addClient(&conn, time.Duration(5)*time.Second, 1000)
			if err != nil {
				logger.Printf("Error Creating Client from Connection: %s", err)
				conn.Close()
				continue
			}
			// Getting a reference to the client object
			if c := client.GetClient(id); c != nil {
				go s.handle(index, c, logger, ctx, &conn)
				continue
			}
			panic(fmt.Errorf("Failed to get newly created client from active client list"))
		}
	}
}

func (s *Server) addClient(conn *net.Conn, tor time.Duration, qs int) (*client.Client, string, error) {
	c, id, err := s.clients.Connect(conn)

	return c, id, err
}

func (s *Server) handle(index int, c *client.Client, logger *log.Logger, ctx context.Context, conn *net.Conn) error {
	err := c.Serve(logger, ctx, conn)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			// Client has closed the connection properly
		}
		if errors.Is(err, syscall.ECONNRESET) {
			// connection has most likely dropped
		}
	}

	// Disconnect the client and remove the client object
	client.Disconect(s.clients[index])
	s.clientmutex.Lock()
	s.clients = append(s.clients[:index], s.clients[index:]...)
	s.clientmutex.Unlock()

	return nil
}
