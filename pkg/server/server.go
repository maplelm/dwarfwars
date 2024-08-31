package server

/*
* WARNING: CURRENTLY NOT VALIDATING CONNECTIONS OR TRAFFIC.
* WARNING: No Authentication is required.
 *
 * FIX: need to acount for bytes not being sent all at once
 * FIX: Connections currently do not timeout leading to a runaway goruitines
*/

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Connection struct {
	ReadData  chan []byte
	WriteData chan []byte
}

func NewConnection(bufSize int) Connection {
	return Connection{
		ReadData:  make(chan []byte, bufSize),
		WriteData: make(chan []byte, bufSize),
	}
}

type Server struct {
	Addr        string         // Address the server will listen on
	Port        string         // Listening Port
	exit        chan struct{}  // used to close down the server
	wait        sync.WaitGroup // tracking goruitines for connections so they can shutdown properly
	IdleTimeout time.Duration
	Shutdown    bool
	connections map[net.Addr]Connection
	connMut     sync.RWMutex
}

func New(addr, port string, timeout time.Duration) (s *Server, err error) {
	s = &Server{
		Addr:        addr,
		Port:        port,
		exit:        make(chan struct{}),
		IdleTimeout: timeout,
		Shutdown:    false,
		connections: make(map[net.Addr]Connection),
	}
	return
}

// Returns the Full Formated listening Address and Port for the server
func (s *Server) FullAddr() string {
	return fmt.Sprintf("%s:%s", s.Addr, s.Port)
}

func (s *Server) StartTCP(ctx context.Context, qChan chan struct{}) (err error) {
	// Create the TCP Listener
	listener, err := net.Listen("tcp", s.FullAddr())
	if err != nil {
		log.Printf("(Server) Failed to create listener for TCP server, %s", err)
		return
	}
	//Base TCP server context
	tcpCtx, tcpCtxCancel := context.WithCancel(ctx)
	defer tcpCtxCancel()

tcplistenloop:
	for {
		select {
		case <-tcpCtx.Done():
			return tcpCtx.Err()
		case <-qChan:
			break tcplistenloop
		default:
			c, err := listener.Accept()
			if err != nil {
				// NOTE: Print out error
			}
			go s.tcpHandle(tcpCtx, c)
		}
	}
	// TCP server is closing
	tcpCtxCancel()
	s.wait.Wait()
	return
}
func (s *Server) StartUDP() (err error) {

	//addr, err := net.ResolveUDPAddr("udp", s.FullAddr())

	if err != nil {
		return
	}
	return
}
func (s *Server) Wait() {
	s.wait.Wait()
}
func (s *Server) tcpHandle(ctx context.Context, conn net.Conn) (err error) {
	s.wait.Add(1)
	defer s.wait.Done()

	s.connMut.Lock()
	s.connections[conn.RemoteAddr()] = NewConnection(10)
	s.connMut.Unlock()
	defer func() {
		s.connMut.Lock()
		delete(s.connections, conn.RemoteAddr())
		s.connMut.Unlock()
	}()

	defer conn.Close()

	var (
		lastRead  time.Time = time.Now()
		lastWrite time.Time = time.Now()
		connWait  sync.WaitGroup
	)
	readCtx, readCtxCancel := context.WithCancel(ctx)
	writeCtx, writeCtxCancel := context.WithCancel(ctx)

	log.Printf("Server Accepted Connection from %s", conn.RemoteAddr())

	// Reading From Connection
	go func(c context.Context, cn net.Conn) error {
		connWait.Add(1)
		defer connWait.Done()
		rc := s.connections[cn.RemoteAddr()].ReadData
		for {
			select {
			case <-c.Done():
				return c.Err()
			default:
				b := []byte{}
				n, err := conn.Read(b)
				if err != nil {
					fmt.Printf("Failed to Read TCP data from %s, Read %d bytes\n", cn.RemoteAddr(), n)
					continue
				}
				rc <- b
			}
		}
	}(readCtx, conn)
	// Writing to Connection
	go func(c context.Context, cn net.Conn) error {
		connWait.Add(1)
		defer connWait.Done()
		wc := s.connections[cn.RemoteAddr()].WriteData
		for {
			select {
			case <-c.Done():
				return c.Err()
			case b := <-wc:
				n, err := cn.Write(b)
				if err != nil {
					fmt.Printf("Failed to Write TCP data from %s, Wrote %d bytes\n", cn.RemoteAddr(), n)
					continue
				}
			}
		}
	}(writeCtx, conn)

	for {
		if time.Since(lastRead) >= time.Minute && time.Since(lastWrite) >= time.Minute {
			readCtxCancel()
			writeCtxCancel()
			log.Printf("(Dwarf Wars Server) Closing Connection due to inactivity (%s)", conn.RemoteAddr())
			conn.Close()
			connWait.Wait()
			return
		}
		time.Sleep(time.Millisecond)
	}
}
