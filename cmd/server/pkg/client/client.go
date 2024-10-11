package client

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/maplelm/dwarfwars/pkg/command"
)

const (
	idsize     = 255
	buffersize = 1024
)

var matrixMutex sync.RWMutex
var unusedIdMatrix []string = []string{}

var activemutex sync.RWMutex
var active map[string]*Client = make(map[string]*Client)

type Client struct {
	ID          string
	TimeoutRate time.Duration
	readQueue   chan command.Command
	writeQueue  chan command.Command
	connection  *net.Conn

	rtMutex      sync.RWMutex
	readtimeouts int
}

func New(c *net.Conn, tor time.Duration, qs int) (string, error) {
	var id string
	if len(unusedIdMatrix) > 0 {
		matrixMutex.Lock()
		id = unusedIdMatrix[0]
		unusedIdMatrix = unusedIdMatrix[1:]
		matrixMutex.Unlock()
	} else {
		var bytes []byte = make([]byte, idsize)
		for {
			_, err := rand.Read(bytes)
			if err != nil {
				return "", fmt.Errorf("Failed to Generate Client ID: %s", err)
			}
			id = string(bytes)
			if _, ok := active[id]; ok {
				continue
			}
			break
		}

	}

	activemutex.Lock()
	defer activemutex.Unlock()
	active[id] = &Client{
		ID:           id,
		TimeoutRate:  tor,
		connection:   c,
		readQueue:    make(chan command.Command, qs),
		writeQueue:   make(chan command.Command, qs),
		readtimeouts: 0,
	}
	return id, nil
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
		err <- read(logger, ctx, c.readQueue, conn, c.TimeoutRate, &(c.readtimeouts), &(c.rtMutex))
	}()
	// Write Outbound Data
	go func() {
		err <- write(logger, ctx, c.writeQueue, conn)
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
	c.rtMutex.RLock()
	defer c.rtMutex.RUnlock()
	return c.readtimeouts > 15
}

func read(logger *log.Logger, ctx context.Context, c chan<- command.Command, conn *net.Conn, tor time.Duration, timeouts *int, toMutex *sync.RWMutex) error {
	var (
		buf    []byte = make([]byte, buffersize)
		msg    []byte = make([]byte, buffersize*3)
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
				(*conn).SetReadDeadline(time.Now().Add(tor)) // tor: time out rate
				n, err := (*conn).Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					var opErr net.Error
					if errors.As(err, &opErr) && opErr.Timeout() {
						toMutex.Lock()
						*timeouts += 1
						toMutex.Unlock()
					}
					if errors.Is(err, net.ErrClosed) {
						return err
					}
					if errors.Is(err, syscall.ECONNRESET) {
						return err
					}
				}
				toMutex.Lock()
				*timeouts = 0
				toMutex.Unlock()
				buflen += n
				msg = append(msg, buf...)
			}
			msg = msg[:buflen]
			cmd, err := command.Unmarshal(msg)
			if err != nil {
				logger.Printf("Failed to Unmarshal Message into a command: %s", err)
				continue
			}
			c <- *cmd
		}
	}
}

func write(logger *log.Logger, ctx context.Context, c <-chan command.Command, conn *net.Conn) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case cmd := <-c:
			for i := 0; i < 3; i++ {
				if _, err := (*conn).Write(cmd.Marshal()); err != nil {
					if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
						logger.Printf("Connection Closed: %s", (*conn).RemoteAddr())
						return ctx.Err()
					}
					logger.Printf("Connection Failed to Write Data (%d) %s", i+1, (*conn).RemoteAddr())
					if i == 2 {
						logger.Printf("Connection Failed to Write Data 3 times, giving up (%s)", (*conn).RemoteAddr())
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

func GetClient(id string) *Client {
	activemutex.RLock()
	defer activemutex.RUnlock()
	if c, ok := active[id]; ok {
		return c
	}
	return nil
}

func Disconect(id string) error {
	if _, ok := active[id]; ok {
		activemutex.Lock()
		defer activemutex.Unlock()
		delete(active, id)
		return nil
	}
	return fmt.Errorf("Invalid id: %s", id)
}
