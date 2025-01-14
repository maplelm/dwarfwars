package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"library/command"

	"server/internal/cache"
	"server/internal/client"
	"server/internal/types"
)

type Server struct {
	Addr     *net.TCPAddr
	Listener *net.TCPListener

	clientmutex *sync.Mutex
	clients     *client.Factory

	ActiveGames []types.Lobby

	quit chan struct{}
}

func New(addr *net.TCPAddr) (*Server, error) {
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Addr:        addr,
		Listener:    l,
		clientmutex: new(sync.Mutex),
		clients:     client.NewFactory(255, 1024, time.Duration(5)*time.Second),
	}, nil
}

func (s *Server) Start(opts *cache.Cache[types.Options], logger *log.Logger, wgrp *sync.WaitGroup, ctx context.Context) (string, error) {
	var (
		conn   net.Conn
		err    error
		client *client.Client
	)

	// If passed a waitgroup add to the waitgroup
	if wgrp == nil {
		wgrp = new(sync.WaitGroup)
	}

	// Defers
	defer s.Listener.Close()

	// Start a new thread that manages the incoming connections
	go s.clients.DispatchIncomingCommands(ctx, logger)

	/////////////////////////////////////
	// Listen For Incoming Connections //
	/////////////////////////////////////
	for {
		select {
		case <-ctx.Done():
			return "Context Closed", ctx.Err()
		case <-s.quit:
			return "Quit Channel Called", ctx.Err()
		default:
			if conn, err = s.Listener.AcceptTCP(); err != nil {
				if errors.Is(err, net.ErrClosed) {
					return fmt.Sprintf("Listener Closed, %s", err), ctx.Err()
				}
				if e, ok := err.(net.Error); ok && e.Timeout() {
					continue
				}
				logger.Printf("Listner Failed to Accept Incoming Connection: %s", err)
				continue
			}

			if client, err = s.clients.Connect(conn, logger); err != nil {
				logger.Printf("Failed to Connect with Client: %s", err)
			}

			logger.Printf("Starting up Monitor for %d", client.Uid())
			go s.clients.Monitor(ctx, client, logger, wgrp)
			cmd, _ := command.New(client.Uid(), command.FormatText, command.TypeWelcome, []byte(strconv.Itoa(int(client.Uid()))))
			cmd.Send(client.Connection)
		}
	}
}

func (s *Server) Stop() {
	s.quit <- struct{}{}
}

func (s *Server) Clients() []uint32 {
	return s.clients.Keys()
}
