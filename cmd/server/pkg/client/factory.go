package client

import (
	"crypto/rand"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/maplelm/dwarfwars/pkg/command"
)

type Factory struct {
	IdSize      int
	BufferSize  int
	TimeoutRate time.Duration

	mutex         sync.RWMutex
	unusedIds     []string
	activeClients map[string]*Client
}

func NewFactory(idsize, buffsize int, tor time.Duration) *Factory {
	return &Factory{
		IdSize:        idsize,
		BufferSize:    buffsize,
		TimeoutRate:   tor,
		unusedIds:     make([]string, 10),
		activeClients: make(map[string]*Client),
	}
}

func (f *Factory) Connect(c *net.Conn) (*Client, string, error) {
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
				return nil, "", err
			}
			id = string(bytes)
			if _, ok := f.activeClients[id]; !ok {
				break
			}
		}
	}

	f.activeClients[id] = &Client{
		id:           id,
		TimeoutRate:  f.TimeoutRate,
		connection:   c,
		readQueue:    make(chan command.Command, f.BufferSize),
		writeQueue:   make(chan command.Command, f.BufferSize),
		readtimeouts: 0,
	}

	return f.activeClients[id], id, nil
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

func (f *Factory) Get(id string) (*Client, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if c, ok := f.activeClients[id]; !ok {
		return nil, fmt.Errorf("factory does not contain that client")
	} else {
		return c, nil
	}
}
