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

type Client struct {
	TimeoutRate time.Duration
	Connection  net.Conn

	readmut      sync.RWMutex
	readtimeouts int

	BufferSize int

	id            uint32
	dispatchQueue chan<- command.Command
	outbound      chan command.Command
}

func New(c net.Conn, tor time.Duration, qs int, id uint32, bs int, in chan<- command.Command) *Client {
	return &Client{
		id:            id,
		TimeoutRate:   tor,
		BufferSize:    bs,
		Connection:    c,
		dispatchQueue: in,
		outbound:      make(chan command.Command, qs),
		readtimeouts:  0,
	}
}

func (c *Client) GetID() uint32 {
	return c.id
}

func (c *Client) Send(cmd *command.Command) int {
	c.outbound <- *cmd
	return len(c.outbound)
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

	for {
		select {
		case <-clientctx.Done():
			return clientctx.Err()
		case e := <-readerr:
			return e
		case e := <-writeerr:
			return e
		}
	}
}

func (c *Client) Close() error {
	return c.Connection.Close()
}

func (c *Client) TimedOut() bool {
	c.readmut.RLock()
	defer c.readmut.RUnlock()
	return c.readtimeouts > 15
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
			c.dispatchQueue <- *cmd
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
		case cmd := <-c.outbound:
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

func (c *Client) Disconect() error {
	close(c.outbound)
	return c.Connection.Close()
}
