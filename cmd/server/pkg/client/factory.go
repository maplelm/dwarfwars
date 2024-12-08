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
	IdSize          uint32
	TimeoutRate     time.Duration
	TimeoutPollRate time.Duration

	IncomingCommands chan *command.Command

	mutex         sync.RWMutex
	unusedIds     []uint32
	nextId        uint32
	activeClients map[uint32]*Client
}

func NewFactory(idsize uint32, tor, tpr time.Duration) *Factory {
	return &Factory{
		IdSize:           idsize,
		TimeoutRate:      tor,
		TimeoutPollRate:  tpr,
		IncomingCommands: make(chan *command.Command, 100),
		unusedIds:        make([]uint32, 0),
		activeClients:    make(map[uint32]*Client),
	}
}

/*
Monitoring Per Client. This function will end when the client it is monitoring is finished
*/
func (f *Factory) Monitor(ctx context.Context, c *Client, logger *log.Logger, wgrp *sync.WaitGroup) error {
	wgrp.Add(1)
	defer wgrp.Done()

	cc := make(chan error)
	defer close(cc)

	go func() {
		cc <- c.Serve(logger, ctx)
	}()

eventloop:
	for {
		select {
		case e := <-cc:
			switch e {
			case net.ErrClosed:
				logger.Printf("Client (%d) closed connection", c.uid)
				break eventloop
			case syscall.ECONNRESET:
				logger.Printf("Client (%d) likely dropped connection", c.uid)
				break eventloop
			default:
				logger.Printf("Client (%d) Error: %s", c.uid, e)
			}
		case <-time.After(f.TimeoutPollRate):
			if time.Since(c.checkin) >= f.TimeoutRate {
				logger.Printf("Client (%d) has timed out", c.uid)
				break eventloop
			}
		}
	}

	return f.Disconnect(c.uid)
}

/*
Default Message dispatcher for clients
clients are designed to switch dispatchers as they swap contexts (what they are doing)
*/
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
				continue
			}

			switch cmd.Type {
			case command.TypeLobbyJoinRequest:
			case command.TypeLobbyLeaveRequest:
			case command.TypeRegister:

			case command.TypeLogin:
			case command.TypeWorldData, command.TypeWorldUpdate:
				// Error State, Client should be passed to a Game Instance Rather then the Dispatcher by the time these messages are being Send
			case command.TypeEcho:
				client.sendToClient <- cmd
			case command.TypeWelcome:
				logger.Printf("Warning: client %d sending Welcome command!", client.uid)
			default:
				response, err := command.New(client.uid, command.FormatText, command.TypeError, []byte("not a supported command type at this time"))
				if err == nil {
					client.sendToClient <- response
				} else {
					logger.Printf("Error: Failed to send error response: %s", err)
				}
			}
		}
	}
}

func (f *Factory) Connect(c net.Conn, logger *log.Logger) (*Client, error) {
	f.mutex.Lock()
	defer func() {
		f.mutex.TryLock()
		f.mutex.Unlock()
	}()

	var id uint32

	if len(f.unusedIds) > 0 {
		logger.Printf("reusing id %d", f.unusedIds[0])
		id = f.unusedIds[0]
		f.unusedIds = f.unusedIds[1:]
	} else {
		var bytes []byte = make([]byte, 4)
		for binary.LittleEndian.Uint32(bytes) == 0 {
			rand.Read(bytes)
		}
		id = binary.LittleEndian.Uint32(bytes)
		logger.Printf("new id %d", id)
	}

	cli := New(c, uint32(100), id, f.IncomingCommands) // New Client
	f.activeClients[id] = cli                          // Adding Client to Active Client List

	// Send client their ID
	cmd, err := command.New(id, command.FormatText, command.TypeWelcome, []byte(fmt.Sprint(id)))
	if err != nil {
		logger.Printf("Error: failed to send client their id, closing connection")
		f.mutex.Unlock()
		f.Disconnect(id)
		return nil, err
	}

	_, err = cmd.Send(cli.Connection)
	if err != nil {
		logger.Printf("Error: failed to send client their id, closing connection")
		f.mutex.Unlock()
		f.Disconnect(id)
		return nil, err
	}

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
	for k := range f.activeClients {
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
