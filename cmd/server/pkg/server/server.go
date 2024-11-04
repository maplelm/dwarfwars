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

	"github.com/maplelm/dwarfwars/cmd/server/pkg/client"
	"github.com/maplelm/dwarfwars/cmd/server/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/server/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/cache"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type Server struct {
	Addr     *net.TCPAddr
	Listener *net.TCPListener
	CC       chan net.Conn // connection channel

	clientmutex sync.Mutex
	clients     *client.Factory

	ActiveGames []*game.Game

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

			go s.clients.Monitor(ctx, client, logger, wgrp)
			cmd, _ := command.New(client.GetID(), command.FormatText, command.CommandWelcome, []byte(strconv.Itoa(int(client.GetID()))))
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
