package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/rs/zerolog"

	"server/internal/cache"
	"server/internal/server"
	"server/internal/types"
)

var (
	connections map[string]net.TCPConn
)

func CliMode(ctx context.Context, logger *zerolog.Logger, opts *cache.Cache[types.Options]) error {
	var (
		err   error
		addr  *net.TCPAddr
		wgrp  *sync.WaitGroup = new(sync.WaitGroup)
		s     *server.Server
		input string
	)

	logger.Info().Str("Mode", "Headless").Msg("Starting Server")

	// Creating Server Context
	cliCtx, cliclose := context.WithCancel(context.Background())
	defer cliclose()

	// Resolving Server Listening Address
	logger.Printf("Resolving TCP Address: %s:%d", opts.MustGetData().Game.Addr, opts.MustGetData().Game.Port)
	if addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", opts.MustGetData().Game.Addr, opts.MustGetData().Game.Port)); err != nil {
		logger.Printf("Failed to Resolve TCP Address for server: %s", err)
		return err
	} else {
		logger.Printf("Resolved Server Address: %s", addr.String())
	}

	s, err = server.New(addr)
	if err != nil {
		logger.Printf("Failed to Create Server Object: %s", err)
		return err
	}

	// Defers
	defer wgrp.Wait()

	// Starting Game Server
	logger.Printf("Starting Dwarf Wars Server...")
	go Server(logger, opts, wgrp, cliCtx)

	// Listening for commands from CLI
	for {
		fmt.Print("(Input): ")
		if _, err := fmt.Scanln(&input); err != nil {
			logger.Info().Err(err).Msg("Failed to Read user input")
		} else {
			switch input {
			case "stop", "quit":
				logger.Info().Msg("Shutting down server")
				fmt.Println("Shutting Down")
				return nil
			case "ls":
				fmt.Printf("clients: %d\n", len(s.Clients()))
				for _, v := range s.Clients() {
					fmt.Printf("\t%d\n", v)
				}
			case "count":
				fmt.Printf("Connections: %d\n", len(server.Connections))
			default:
				fmt.Printf("Invalid input (%s)\n", input)
			}
		}
	}
}

func Server(logger *zerolog.Logger, opts *cache.Cache[types.Options], wgrp *sync.WaitGroup, cliCtx context.Context) error {

	return nil
}
