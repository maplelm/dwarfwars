package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/maplelm/dwarfwars/cmd/server/pkg/client"
	"github.com/maplelm/dwarfwars/pkg/cache"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type ConnectionHandler interface {
	Serve(*log.Logger, context.Context, *net.Conn) error
}

type ConnectionHandle struct {
	f func(*log.Logger, context.Context, *net.Conn) error
}

func (ch *ConnectionHandle) Serve(l *log.Logger, ctx context.Context, conn *net.Conn) error {
	return ch.f(l, ctx, conn)
}

func ConnectionHandlerFunc(f func(*log.Logger, context.Context, *net.Conn) error) ConnectionHandler {
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

	clientmutex sync.Mutex
	clients     []string
}

func NewServer(addr *net.TCPAddr, chanSize int, handle ConnectionHandler) (*Server, error) {
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Addr:     addr,
		Listener: l,
		Handle:   handle,
		CC:       make(chan net.Conn, chanSize),
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

			id, err := client.New(&conn, time.Duration(5)*time.Second, 1000)
			if err != nil {
				fmt.Errorf("Error Creating Client from Connection: %s", err)
				conn.Close()
				continue
			}
			s.clientmutex.Lock()
			s.clients = append(s.clients, id)
			index := len(s.clients) - 1
			s.clientmutex.Unlock()

			c := client.GetClient(id)
			if c != nil {
				go func(index int) {
					c.Serve(logger, ctx, &conn)
					client.Disconect(id)
					s.clientmutex.Lock()
					s.clients = append(s.clients[:index], s.clients[index:]...)
					s.clientmutex.Unlock()

				}(index)
			} else {
				panic(fmt.Errorf("Failed to get newly created client from active client list"))
			}
		}
	}
}

func EchoConnection(logger *log.Logger, ctx context.Context, conn *net.Conn) {
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

func CommandHandlerTest(logger *log.Logger, ctx context.Context, conn *net.Conn) {
	var fullData []byte
	readCount := 0
	data := make([]byte, 1024)
	for {
		n, err := (*conn).Read(data)
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			logger.Printf("Connection (%s): %s", (*conn).RemoteAddr(), err)
			break
		}
		readCount += n
		fullData = append(fullData, data...)
	}

	cmd, err := command.Unmarshal(fullData[:readCount])
	if err != nil {
		logger.Printf("Connection Command Err (%s): %s", (*conn).RemoteAddr(), err)

	}

	logger.Printf("Incomming Command (%s):\n\tVersion: %d\n\tType: %d\n\tSize: %d\n\t%s\n", (*conn).RemoteAddr(), cmd.Version, cmd.Type, cmd.Size, string(cmd.Data))

	_, err = (*conn).Write(cmd.Marshal())
	if err != nil {
		logger.Printf("Connection Error Sending Commend (%s): %s", (*conn).RemoteAddr(), err)
	}
}
