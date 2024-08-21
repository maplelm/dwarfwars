package server

import (
	"context"
	"net"
	"sync"
)

var (
	clientIdMutex   sync.Mutex
	clientIdCounter int   = 0
	unusedClientIds []int = []int{}
)

type client struct {
	id   int             // Connection ID
	uid  *int            // User ID ( Null if not logged in )
	conn *net.Conn       // Connection State
	ctx  context.Context // Connection context
}

func NewClient(ctx context.Context, conn *net.Conn) *client {
	var (
		id int
	)
	clientIdMutex.Lock()
	defer clientIdMutex.Unlock()
	if len(unusedClientIds) > 0 {
		id = unusedClientIds[0]
		unusedClientIds = unusedClientIds[1:]
	} else {
		id = clientIdCounter
		clientIdCounter++
	}

	return &client{
		id:   id,
		uid:  nil,
		conn: conn,
		ctx:  ctx,
	}

}

func (c *client) ConnectionHandler(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}
