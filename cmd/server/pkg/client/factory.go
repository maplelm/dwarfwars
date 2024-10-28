package client

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/maplelm/dwarfwars/pkg/command"
)

type Factory struct {
	IdSize      int
	BufferSize  int
	TimeoutRate time.Duration

	IncomingConnections chan net.Conn
	IncomingCommands    chan command.Command

	mutex         sync.RWMutex
	unusedIds     []string
	activeClients map[string]*Client
}

func NewFactory(idsize, buffsize int, tor time.Duration) *Factory {
	return &Factory{
		IdSize:              idsize,
		BufferSize:          buffsize,
		TimeoutRate:         tor,
		IncomingConnections: make(chan net.Conn, 100),
		IncomingCommands:    make(chan command.Command, 100),
		unusedIds:           make([]string, 10),
		activeClients:       make(map[string]*Client),
	}
}

func (f *Factory) MonitorIncomingConnections(ctx context.Context, logger *log.Logger) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c := <-f.IncomingConnections:
			client, err := f.Connect(c)
			if err != nil {
				logger.Printf("Failed to Connect with Client: %s", err)
			}
			// kicking off service to client
			go func() {
				err := client.Serve(logger, ctx)
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						// Client has closed the connection properly
					}
					if errors.Is(err, syscall.ECONNRESET) {
						// connection has most likely dropped
					}
				}
				f.Disconnect(client.GetID())
			}()
		}
	}
}

func (f *Factory) DispatchIncomingCommands(ctx context.Context, logger *log.Logger) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case cmd := <-f.IncomingCommands:
			logger.Printf("Command Received from: %s, dispatching", cmd.ClientID)
		}
	}
}

func (f *Factory) Connect(c net.Conn) (*Client, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	var id string

	if len(f.unusedIds) > 0 {
		id = f.unusedIds[0]
		f.unusedIds = f.unusedIds[1:]
	} else {
		var bytes []byte = make([]byte, f.IdSize)
		for {
			_, err := rand.Read(bytes)
			if err != nil {
				return nil, err
			}
			id = string(bytes)
			if _, ok := f.activeClients[id]; !ok {
				break
			}
		}
	}

	cli := New(c, f.TimeoutRate, 100, id, f.BufferSize, f.IncomingCommands)
	f.activeClients[id] = cli

	return f.activeClients[id], nil
}

func (f *Factory) Disconnect(id string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if c, ok := f.activeClients[id]; !ok {
		return fmt.Errorf("factory does not contain that client")
	} else {
		c.Disconect()
		delete(f.activeClients, id)
		return nil
	}
}

func (f *Factory) Keys() []string {
	keys := make([]string, len(f.activeClients))
	index := 0
	for k, _ := range f.activeClients {
		keys[index] = k
		index++
	}
	return keys
}

func (f *Factory) Get(id string) (*Client, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if c, ok := f.activeClients[id]; !ok {
		return nil, fmt.Errorf("factory does not contain that client")
	} else {
		return c, nil
	}
}
