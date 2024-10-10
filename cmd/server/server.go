package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/maplelm/dwarfwars/pkg/cache"
)

type ConnectionHandler interface {
	Serve(*log.Logger, context.Context, *net.Conn)
}

type ConnectionHandle struct {
	f func(*log.Logger, context.Context, *net.Conn)
}

func (ch *ConnectionHandle) Serve(l *log.Logger, ctx context.Context, conn *net.Conn) {
	ch.f(l, ctx, conn)
}

func ConnectionHandlerFunc(f func(*log.Logger, context.Context, *net.Conn)) ConnectionHandler {
	obj := ConnectionHandle{
		f: f,
	}
	return &obj
}

type Server struct {
	Addr     *net.TCPAddr
	Listener *net.TCPListener
	CC       chan net.Conn
	Handle   ConnectionHandler

	connMutex   sync.Mutex
	Connections map[net.Addr]*net.Conn
	conns       sync.WaitGroup
}

func NewServer(addr *net.TCPAddr, chanSize int, handle ConnectionHandler) (*Server, error) {
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Connections: make(map[net.Addr]*net.Conn),
		Addr:        addr,
		Listener:    l,
		Handle:      handle,
		CC:          make(chan net.Conn, chanSize),
	}, nil
}

func (s *Server) Start(opts *cache.Cache[Options], logger *log.Logger, wgrp *sync.WaitGroup, ctx context.Context) error {
	if wgrp != nil {
		wgrp.Add(1)
		defer wgrp.Done()
	}

	defer s.Listener.Close()

	serverCtx, close := context.WithCancel(ctx)
	defer close()

	go s.connectionManager(logger, serverCtx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			s.Listener.SetDeadline(time.Now().Add(time.Second))
			if conn, err := s.Listener.AcceptTCP(); err != nil {
				if errors.Is(err, net.ErrClosed) {
					logger.Printf("Listener Close: %s", err)
					return err
				}
				var netErr *net.OpError
				if errors.As(err, &netErr) && netErr.Timeout() {
				} else {
					logger.Printf("Listner Failed to Accept Incoming Connection: %s", err)
				}
			} else {
				options, err := opts.Get()
				if err != nil {
					logger.Printf("Warning: Server Failed to Get Connection Deadlines from Settings, %s", err)
					conn.SetReadDeadline(time.Now().Add(time.Duration(1) * time.Minute))
					conn.SetWriteDeadline(time.Now().Add(time.Duration(1) * time.Minute))
				} else {
					conn.SetReadDeadline(time.Now().Add(time.Duration(options.Game.Timeouts.Read) * time.Millisecond))
					conn.SetWriteDeadline(time.Now().Add(time.Duration(options.Game.Timeouts.Write) * time.Millisecond))
				}
				s.CC <- conn
			}
		}
	}
}

func (s *Server) connectionManager(logger *log.Logger, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case conn := <-s.CC:
			s.connMutex.Lock()
			s.Connections[conn.RemoteAddr()] = &conn
			s.connMutex.Unlock()
			go func() {
				s.connMutex.Lock()
				s.Connections[conn.RemoteAddr()] = &conn
				s.connMutex.Unlock()

				s.Handle.Serve(logger, ctx, &conn)

				s.connMutex.Lock()
				delete(s.Connections, conn.RemoteAddr())
				s.connMutex.Unlock()
			}()
		}
	}
}

func EchoConnection(logger *log.Logger, ctx context.Context, conn *net.Conn) {
	defer (*conn).Close()

	var (
		data []byte = make([]byte, 2000)
		n    int
		err  error
	)

	// Reading to Connection
	if n, err = (*conn).Read(data); n > 0 && err == nil {
		logger.Printf("Message from connection (%s): %s", (*conn).RemoteAddr(), string(data))
	} else if err != nil {
		logger.Printf("Failed to read data from client: %s", err)
		data = []byte(err.Error())
	} else {
		logger.Printf("No data to read (%s)", (*conn).RemoteAddr())
	}

	// Writing to Connection
	if n, err = (*conn).Write(data[:n]); n > 0 && err == nil {
		logger.Printf("Message Sent to Conneciton (%s): %s", (*conn).RemoteAddr(), string(data))
	} else if err != nil {
		logger.Printf("Failed to Write data to Client: %s", err)
	} else {
		logger.Printf("No data Writen (%s)", (*conn).RemoteAddr())
	}
}
