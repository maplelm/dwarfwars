package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"server/internal/cache"
	"server/internal/server"
	"server/internal/types"
)

func CliMode(logger *log.Logger, opts *cache.Cache[types.Options]) error {
	var (
		err    error
		addr   *net.TCPAddr
		wgrp   *sync.WaitGroup = new(sync.WaitGroup)
		s      *server.Server
		clistd *log.Logger = log.New(os.Stdout, "Dwarf Wars Server: ", 0)
		clierr *log.Logger = log.New(os.Stderr, "Error: ", 0)
		input  string
	)

	logger.Println("Server Mode: Headless")

	// Creating Server Context
	cliCtx, close := context.WithCancel(context.Background())

	// Resolving Server Listening Address
	logger.Printf("Resolving TCP Address: %s:%d", opts.MustGetData().Game.Addr, opts.MustGetData().Game.Port)
	if addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", opts.MustGetData().Game.Addr, opts.MustGetData().Game.Port)); err != nil {
		logger.Printf("Failed to Resolve TCP Address for server: %s", err)
		close()
		return err
	} else {
		logger.Printf("Resolved Server Address: %s", addr.String())
	}

	s, err = server.New(addr)
	if err != nil {
		logger.Printf("Failed to Create Server Object: %s", err)
		close()
		return err
	}

	// Defers
	defer wgrp.Wait()
	defer close()

	// Starting Game Server
	logger.Printf("Starting Dwarf Wars Server...")
	go s.Start(opts, logger, wgrp, cliCtx)

	// Listening for commands from CLI
	for {
		clistd.Printf("(Input)")
		if _, err := fmt.Scanln(&input); err != nil {
			clierr.Printf("Failed to Read user input: %s", err)
		} else {
			switch input {
			case "stop", "quit":
				clistd.Printf("Shutting Down")
				return nil
			case "ls":
				clistd.Printf("clients: %d", len(s.Clients()))
				for _, v := range s.Clients() {
					clistd.Printf("\t%d", v)
				}
			case "count":
				fmt.Printf("Connections: %d\n", len(server.Connections))
			default:
				clierr.Printf("Invalid input (%s)\n", input)
			}
		}
	}
}
