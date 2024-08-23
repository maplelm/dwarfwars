package tcp

import (
	"net"
	"sync"
	"time"
)

var (
	connectionIdCounter uint64 = 0
	unusedConnectionIds []uint64
	connectionIdMutex   sync.Mutex
)

type connection struct {
	Tcp          net.TCPConn // TCP Connection object
	Udp          net.UDPConn
	Id           uint64 // Connection ID
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewConnection(c net.Conn) *connection {
	var id uint64
	connectionIdMutex.Lock()
	defer connectionIdMutex.Unlock()
	if len(unusedConnectionIds) > 0 {
		id = unusedConnectionIds[0]
		unusedConnectionIds = unusedConnectionIds[1:]
	} else {
		id = connectionIdCounter
		connectionIdCounter++
	}
	return &connection{
		Tcp: c,
		Id:  id,
	}
}

func RemoveConnection(c *connection) error {
	connectionIdMutex.Lock()
	defer connectionIdMutex.Unlock()
	unusedConnectionIds = append(unusedConnectionIds, c.Id)
	return c.Tcp.Close()
}

func (c *connection) ReadTCP() (cmd Command, n int, err error) {
	var bytes []byte
	n, err = c.Tcp.Read(bytes)
	if err != nil {
		return
	}
	err = cmd.UnmarshalBinary(bytes)
	return
}

func (c *connection) ReadUDP() (cmd Command, n int, err error) {
	var bytes []byte
	n, err = c.Udp.Read(bytes)
	if err != nil {
		return
	}
	err = cmd.UnmarshalBinary(bytes)
	return
}

func (c *connection) WriteTCP(cmd Command) (n int, err error) {
	bytes, err := cmd.MarshalBinary()
	if err != nil {
		return 0, err
	}
	n, err = c.Tcp.Write(bytes)
	return
}

func (c *connection) WriteUDP(cmd Command) (n int, err error) {
	bytes, err := cmd.MarshalBinary()
	if err != nil {
		return 0, err
	}
	n, err = c.Udp.Write(bytes)
	return
}