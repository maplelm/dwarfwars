package server

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/maplelm/dwarfwars/cmd/server/pkg/client"
	"github.com/maplelm/dwarfwars/cmd/server/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/cache"
)

type Server struct {
	Addr     *net.TCPAddr
	Listener *net.TCPListener
	CC       chan net.Conn // connection channel

	clientmutex sync.Mutex
	clients     *client.Factory

	quit chan struct{}
}

func New(addr *net.TCPAddr, chanSize int) (*Server, error) {
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Addr:     addr,
		Listener: l,
		CC:       make(chan net.Conn, chanSize),
		clients:  client.NewFactory(255, 1024, time.Duration(5)*time.Second),
	}, nil
}

func (s *Server) Start(opts *cache.Cache[types.Options], logger *log.Logger, wgrp *sync.WaitGroup, ctx context.Context) error {
	if wgrp != nil {
		wgrp.Add(1)
		defer wgrp.Done()
	}

	// Create server base context
	serverCtx, close := context.WithCancel(ctx)

	// defer the stopping of the server context and then close the TCP listener
	defer s.Listener.Close()
	defer close()

	// Start a new thread that manages the incoming connections
	go s.clients.MonitorIncomingConnections(serverCtx, logger)
	go s.clients.DispatchIncomingCommands(serverCtx, logger)

	/////////////////////////////////////
	// Listen For Incoming Connections //
	/////////////////////////////////////
ListenLoop:
	for {
		select {
		case <-serverCtx.Done():
			break ListenLoop
		case <-s.quit:
			break ListenLoop
		default:
			if conn, err := s.Listener.AcceptTCP(); err != nil {
				if errors.Is(err, net.ErrClosed) {
					logger.Printf("Listener Close: %s", err)
					break ListenLoop
				}
				var netErr *net.OpError
				if errors.As(err, &netErr) && netErr.Timeout() {
				} else {
					logger.Printf("Listner Failed to Accept Incoming Connection: %s", err)
				}
			} else {
				s.clients.IncomingConnections <- conn
			}
		}
	}
	return serverCtx.Err()
}

func (s *Server) Stop() {
	s.quit <- struct{}{}
}
