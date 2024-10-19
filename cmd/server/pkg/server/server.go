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
	go s.connMgr(logger, serverCtx)
	// listen for incoming connections
	go s.listen(logger)

	////////////////////////////////////////
	// Blocking until Server context ends //
	////////////////////////////////////////
	select {
	case <-ctx.Done():
		// stop if the application context is closed
		return ctx.Err()
	case <-s.quit:
		// stop if the quit funciton is called from anywhere else in the application
		return nil
	}
}

func (s *Server) Stop() {
	s.quit <- struct{}{}
}

func (s *Server) listen(logger *log.Logger) error {
	for {
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
