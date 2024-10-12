package server

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/maplelm/dwarfwars/cmd/server/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/cache"
)

type Server struct {
	Addr     *net.TCPAddr
	Listener *net.TCPListener
	CC       chan net.Conn // connection channel

	clientmutex sync.Mutex
	clients     []string

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
	}, nil
}

func (s *Server) Start(opts *cache.Cache[types.Options], logger *log.Logger, wgrp *sync.WaitGroup, ctx context.Context) error {
	if wgrp != nil {
		wgrp.Add(1)
		defer wgrp.Done()
	}
	serverCtx, close := context.WithCancel(ctx)

	defer s.Listener.Close()
	defer close()

	go s.connMgr(logger, serverCtx)
	go s.listen(logger)

	////////////////////////////////////////
	// Blocking until Server context ends //
	////////////////////////////////////////
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.quit:
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
