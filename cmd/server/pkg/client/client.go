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
	Connection  *net.Conn

	readmut      sync.RWMutex
	readtimeouts int

	BufferSize int

	id         string
	readQueue  chan command.Command
	writeQueue chan command.Command
}

func New(c *net.Conn, tor time.Duration, qs int, id string, bs int) (*Client, string) {
	return &Client{
		id:           id,
		TimeoutRate:  tor,
		BufferSize:   bs,
		Connection:   c,
		readQueue:    make(chan command.Command, qs),
		writeQueue:   make(chan command.Command, qs),
		readtimeouts: 0,
	}, id
}

func (c *Client) GetID() string {
	return c.id
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

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-err:
			return e
		case msg := <-c.readQueue:
			logger.Printf(string(msg.Marshal()))
			c.writeQueue <- msg
		}
	}
}

func (c *Client) Close() error {
	return (*c.Connection).Close()
}

func (c *Client) TimedOut() bool {
	c.readmut.RLock()
	defer c.readmut.RUnlock()
	return c.readtimeouts > 15
}

func (c *Client) read(logger *log.Logger, ctx context.Context) error {
	var (
		header       []byte = make([]byte, command.HeaderSize)
		buffer       []byte
		msg          []byte
		timeoutCount int = 0
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

			_, err := (*c.Connection).Read(header)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return err
				}
				var opErr net.Error
				if errors.As(err, &opErr) && opErr.Timeout() {
					timeoutCount++
					continue
				}
				if errors.Is(err, syscall.ECONNRESET) {
					return err
				}
			}
			l, _, err := command.ValidateHeader(header)
			if err != nil {
				logger.Printf("failed to validate header: %s", err)
				continue
			}

			buffer = make([]byte, l)
			n, err := (*c.Connection).Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return err
				}
				var opErr net.Error
				if errors.As(err, &opErr) && opErr.Timeout() {
					timeoutCount++
					continue
				}
				if errors.Is(err, syscall.ECONNRESET) {
					return err
				}
			}
			if n != int(l) {
				logger.Printf("Warning, did not get expected command length from client: %d, %d", n, l)
			}

			msg = make([]byte, int(l)+(int(command.HeaderSize)))
			for i, v := range header {
				msg[i] = v
			}
			for i, v := range buffer {
				msg[i+int(int(command.HeaderSize))] = v
			}

			cmd, err := command.Unmarshal(msg)
			if err != nil {
				logger.Printf("Error Unmarshalling command: %s", err)
				continue
			}
			c.readQueue <- *cmd

			/*
				msg = []byte{}
				buflen = 0
				for {
					//(*c.Connection).SetReadDeadline(time.Now().Add(c.TimeoutRate)) // tor: time out rate
					n, err := (*c.Connection).Read(buf)
					if err != nil {
						logger.Printf("Read Error (%s): %s", (*c.Connection).RemoteAddr(), err)
						if errors.Is(err, io.EOF) {
							break
						}
						var opErr net.Error
						if errors.As(err, &opErr) && opErr.Timeout() {
							/*
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
					/*
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
			*/
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
				if _, err := (*c.Connection).Write(cmd.Marshal()); err != nil {
					if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
						logger.Printf("Connection Closed: %s", (*c.Connection).RemoteAddr())
						return ctx.Err()
					}
					logger.Printf("Connection Failed to Write Data (%d) %s", i+1, (*c.Connection).RemoteAddr())
					if i == 2 {
						logger.Printf("Connection Failed to Write Data 3 times, giving up (%s)", (*c.Connection).RemoteAddr())
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
	return (*c.Connection).Close()
}
