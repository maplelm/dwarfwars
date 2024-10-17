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

	readmut      sync.RWMutex
	readtimeouts int

	BufferSize int

	id         string
	readQueue  chan command.Command
	writeQueue chan command.Command
	connection *net.Conn
}

func New(c *net.Conn, tor time.Duration, qs int, id string, bs int) (*Client, string) {
	return &Client{
		id:           id,
		TimeoutRate:  tor,
		BufferSize:   bs,
		connection:   c,
		readQueue:    make(chan command.Command, qs),
		writeQueue:   make(chan command.Command, qs),
		readtimeouts: 0,
	}, id
}

func (c *Client) Read() (*command.Command, int) {
	if len(c.readQueue) == 0 {
		return nil, 0
	}
	if cmd, ok := <-c.readQueue; !ok {
		return nil, 0
	} else {
		return &cmd, len(c.readQueue)
	}
}

func (c *Client) Write(cmd *command.Command) int {
	c.writeQueue <- *cmd
	return len(c.writeQueue)
}

func (c *Client) Serve(logger *log.Logger, ctx context.Context, conn *net.Conn) error {

	var err chan error = make(chan error)

	// Read Inbound Data
	go func() {
		err <- c.read(logger, ctx)
	}()
	// Write Outbound Data
	go func() {
		err <- c.write(logger, ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-err:
		return e
	}
}

func (c *Client) Close() error {
	return (*c.connection).Close()
}

func (c *Client) TimedOut() bool {
	c.readmut.RLock()
	defer c.readmut.RUnlock()
	return c.readtimeouts > 15
}

func (c *Client) read(logger *log.Logger, ctx context.Context) error {
	var (
		buf    []byte = make([]byte, c.BufferSize)
		msg    []byte = make([]byte, c.BufferSize*3)
		buflen int
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg = []byte{}
			buflen = 0
			for {
				(*c.connection).SetReadDeadline(time.Now().Add(c.TimeoutRate)) // tor: time out rate
				n, err := (*c.connection).Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					var opErr net.Error
					if errors.As(err, &opErr) && opErr.Timeout() {
						c.readmut.Lock()
						c.readtimeouts += 1
						c.readmut.Unlock()
					}
					if errors.Is(err, net.ErrClosed) {
						return err
					}
					if errors.Is(err, syscall.ECONNRESET) {
						return err
					}
				}
				c.readmut.Lock()
				c.readtimeouts = 0
				c.readmut.Unlock()
				buflen += n
				msg = append(msg, buf...)
			}
			msg = msg[:buflen]
			cmd, err := command.Unmarshal(msg)
			if err != nil {
				logger.Printf("Failed to Unmarshal Message into a command: %s", err)
				continue
			}
			c.readQueue <- *cmd
		}
	}
}

func (c *Client) write(logger *log.Logger, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case cmd := <-c.writeQueue:
			for i := 0; i < 3; i++ {
				if _, err := (*c.connection).Write(cmd.Marshal()); err != nil {
					if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
						logger.Printf("Connection Closed: %s", (*c.connection).RemoteAddr())
						return ctx.Err()
					}
					logger.Printf("Connection Failed to Write Data (%d) %s", i+1, (*c.connection).RemoteAddr())
					if i == 2 {
						logger.Printf("Connection Failed to Write Data 3 times, giving up (%s)", (*c.connection).RemoteAddr())
					} else {
						time.Sleep(100 * time.Nanosecond)
					}
					continue
				}
				break
			}
		}
	}
}

func (c *Client) Disconect() error {
	close(c.readQueue)
	close(c.writeQueue)
	return (*c.connection).Close()
}
