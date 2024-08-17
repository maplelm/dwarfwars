package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

const (
	Version    byte = 1
	HeaderSize byte = 3
)

type TcpCommand struct {
	Command uint16
	Data    []byte
}

func (t *TcpCommand) MarshalBinary(b []byte) error {
	if b[1] != Version {
		return fmt.Errorf("Version mismatch, expected %d got %d", Version, b[0])

	}

	return nil
}

func (t *TcpCommand) UnmarshalBinary() (b []byte, err error) {
	return
}

type Server struct {
	Welcomes    []WelcomeCmd
	Connections []net.Conn
	mutex       sync.RWMutex
}

type WelcomeCmd func() *TcpCommand
