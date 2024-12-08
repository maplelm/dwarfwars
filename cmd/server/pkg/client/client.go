package client

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/maplelm/dwarfwars/pkg/command"
)

/*
This class is the base unit for the server. clients will be based around and
they will be the onlyes activating work. The only exception would be the game
world goruitines.
*/
type Client struct {

	// Client State
	Account        *Account // Links client with a game account
	GameInstanceID uint32   // What game the client is currently engaged with

	// Network Connection State
	sendToInternal chan<- *command.Command // Channel the client will use to send message received from connection to an internal system within the server
	sendToClient   chan *command.Command   // commands that will be sent to the connection
	Connection     net.Conn
	checkin        time.Time // last time a message was recieved over connection

	// Network Security State
	Secret []byte
	uid    uint32

	// Concurrency
	readlock *sync.Mutex
}

func New(c net.Conn, qs, uid uint32, InternalSystemQueue chan<- *command.Command) *Client {
	return &Client{
		uid:            uid,
		Connection:     c,
		sendToInternal: InternalSystemQueue,
		sendToClient:   make(chan *command.Command, qs),
		readlock:       new(sync.Mutex),
	}
}

func (c *Client) Uid() uint32 {
	return c.uid
}

func (c *Client) Send(cmd *command.Command) int {
	c.sendToClient <- cmd
	return len(c.sendToClient)
}

func (c *Client) Serve(logger *log.Logger, ctx context.Context) error {

	readerr := make(chan error)
	writeerr := make(chan error)

	clientctx, cancelConn := context.WithCancel(ctx)
	defer cancelConn()

	// Read Inbound Data
	go c.read(readerr, logger, clientctx)
	// Write Outbound Data
	go c.write(writeerr, logger, clientctx)

	select {
	case <-clientctx.Done():
		return clientctx.Err()
	case e := <-readerr:
		return e
	case e := <-writeerr:
		return e
	}
}

func (c *Client) Close() error {
	return c.Connection.Close()
}

func (c *Client) read(e chan<- error, logger *log.Logger, ctx context.Context) error {
	var timeoutCount int = 0
	defer close(e)
	for {
		select {
		case <-ctx.Done():
			e <- ctx.Err()
			return ctx.Err()
		default:

			cmd, err := command.Recieve(c.Connection)
			if err != nil {
				if errors.Is(err, io.EOF) {
					e <- err
					return err
				} else if opErr, ok := err.(net.Error); ok && opErr.Timeout() {
					timeoutCount++
					continue
				} else if errors.Is(err, syscall.ECONNRESET) {
					e <- err
					return err
				} else {
					logger.Printf("Network Read Error: %s", err)
				}
			}
			c.sendToInternal <- cmd
			c.checkin = time.Now()
		}
	}
}

func (c *Client) write(e chan<- error, logger *log.Logger, ctx context.Context) error {
	defer close(e)
	for {
		select {
		case <-ctx.Done():
			e <- ctx.Err()
			return ctx.Err()
		case cmd := <-c.sendToClient:
			for i := 0; i < 3; i++ {
				if n, err := cmd.Send(c.Connection); err != nil {
					if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
						logger.Printf("Connection Closed: %s", c.Connection.RemoteAddr())
						e <- ctx.Err()
						return ctx.Err()
					}
					logger.Printf("Connection Failed to Write Data (%d) %s", i+1, c.Connection.RemoteAddr())
					if i == 2 {
						logger.Printf("Connection Failed to Write Data 3 times, giving up (%s)", c.Connection.RemoteAddr())
					} else {
						time.Sleep(100 * time.Nanosecond)
					}
				} else {
					logger.Printf("Network Write: %d bytes", n)
					break
				}
			}
		}
	}
}

// This is where the client will join an active game that has started
func (c *Client) JoinGame(id uint32) error {
	return nil
}

func (c *Client) Disconect() error {
	close(c.sendToClient)
	return c.Connection.Close()
}
