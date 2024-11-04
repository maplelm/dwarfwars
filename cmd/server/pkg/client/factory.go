package client

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/maplelm/dwarfwars/pkg/command"
)

type Factory struct {
	IdSize      int
	BufferSize  int
	TimeoutRate time.Duration

	IncomingCommands chan command.Command

	mutex         sync.RWMutex
	unusedIds     []uint32
	nextId        uint32
	activeClients map[uint32]*Client
}

func NewFactory(idsize, buffsize int, tor time.Duration) *Factory {
	return &Factory{
		IdSize:           idsize,
		BufferSize:       buffsize,
		TimeoutRate:      tor,
		IncomingCommands: make(chan command.Command, 100),
		unusedIds:        make([]uint32, 0),
		activeClients:    make(map[uint32]*Client),
	}
}

func (f *Factory) Monitor(ctx context.Context, c *Client, logger *log.Logger, wgrp *sync.WaitGroup) error {
	wgrp.Add(1)
	defer wgrp.Done()

	err := c.Serve(logger, ctx)
	switch err {
	case net.ErrClosed:
		logger.Printf("Client (%d) closed connection", c.id)
	case syscall.ECONNRESET:
		logger.Printf("Client (%d) likely dropped connection", c.id)
	default:
		logger.Printf("Client (%d) Error: %s", c.id, err)
	}
	return f.Disconnect(c.id)
}

func (f *Factory) DispatchIncomingCommands(ctx context.Context, logger *log.Logger) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case cmd := <-f.IncomingCommands:
			logger.Printf("Command Received from: %s, dispatching", strconv.Itoa(int(cmd.ClientID)))
			client, ok := f.activeClients[cmd.ClientID]
			if !ok {
				logger.Printf("Error: Failed to get client with id: %s", strconv.Itoa(int(cmd.ClientID)))
				for k, _ := range f.activeClients {
					logger.Printf("\t client key: %d", k)
				}
				continue
			}
			response, err := command.New(client.id, 0, command.CommandType(0), []byte("command received"))
			if err != nil {
				logger.Printf("Error (%d): %s", client.id, err)
				continue
			}
			client.outbound <- *response
		}
	}
}

func (f *Factory) Connect(c net.Conn, logger *log.Logger) (*Client, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	var id uint32

	if len(f.unusedIds) > 0 {
		logger.Printf("reusing id %d", f.unusedIds[0])
		id = f.unusedIds[0]
		f.unusedIds = f.unusedIds[1:]
	} else {
		var bytes []byte = make([]byte, 4)
		rand.Read(bytes)
		id = binary.LittleEndian.Uint32(bytes)
		logger.Printf("new id %d", id)
	}

	cli := New(c, f.TimeoutRate, 100, id, f.BufferSize, f.IncomingCommands)
	f.activeClients[id] = cli

	return f.activeClients[id], nil
}

func (f *Factory) Disconnect(id uint32) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if c, ok := f.activeClients[id]; !ok {
		return fmt.Errorf("factory does not contain that client")
	} else {
		c.Disconect()
		delete(f.activeClients, id)
		f.unusedIds = append(f.unusedIds, id)
		return nil
	}
}

func (f *Factory) Keys() []uint32 {
	keys := make([]uint32, len(f.activeClients))
	index := 0
	for k, _ := range f.activeClients {
		keys[index] = k
		index++
	}
	return keys
}

func (f *Factory) Get(id uint32) (*Client, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if c, ok := f.activeClients[id]; !ok {
		return nil, fmt.Errorf("factory does not contain that client")
	} else {
		return c, nil
	}
}
