package main

import (
	// STD Packages
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	// Project Packages
	"github.com/maplelm/dwarfwars/cmd/server/pkg/server"
	"github.com/maplelm/dwarfwars/cmd/server/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/cache"
)

// CLI Flags
var (
	configPath *string = flag.String("c", "./config/", "location of settings files")
	headless   *bool   = flag.Bool("h", false, "server will not use a tui and can be automated with scripts")
)

// Main Function
func main() {
	// Parse CLI Flags
	flag.Parse()

	// Getting Settings from TOML file
	opts := InitOptionsCache()

	// Initializing the Logger object
	MainLogger := InitLogger(opts)

	// Validating the SQL Server
	MainLogger.Println("Validating Database Before Server Bootup")
	if err := ValidateSQL(3, 500, MainLogger, opts); err != nil {
		MainLogger.Fatalf("Failed to Validate SQL Server: %s", err)
	}

	// Start Server
	switch *headless {
	case true:
		if err := CliMode(MainLogger, opts); err != nil {
			MainLogger.Fatalf("Server Error: %s", err)
		}
	case false:
		if err := TuiMode(MainLogger, opts); err != nil {
			MainLogger.Fatalf("Server Error: %s", err)
		}
	}
}

func CliMode(logger *log.Logger, opts *cache.Cache[types.Options]) error {
	var (
		err    error
		addr   *net.TCPAddr
		wgrp   *sync.WaitGroup = new(sync.WaitGroup)
		s      *server.Server
		clistd *log.Logger = log.New(os.Stdout, "Dwarf Wars Server: ", 0)
		clierr *log.Logger = log.New(os.Stderr, "Error: ", 0)
		input  *string     = new(string)
	)

	logger.Println("Server Mode: Headless")

	// Creating Server Context
	cliCtx, close := context.WithCancel(context.Background())

	// Resolving Server Listening Address
	logger.Printf("Resolving TCP Address: %s:%d", opts.MustGet().Game.Addr, opts.MustGet().Game.Port)
	if addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", opts.MustGet().Game.Addr, opts.MustGet().Game.Port)); err != nil {
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
		if _, err := fmt.Scanln(input); err != nil {
			clierr.Printf("Failed to Read user input: %s", err)
		} else {
			switch *input {
			case "stop", "quit":
				clistd.Printf("Shutting Down")
				return nil
			case "ls":
				clistd.Printf("clients: %d", len(s.Clients()))
				for _, v := range s.Clients() {
					clistd.Printf("\t%d", v)
				}
				// list Connections
			case "count":
				//fmt.Printf("Connections: %d\n", len(server.Connections))
			default:
				clierr.Printf("Invalid input (%s)\n", *input)
			}
		}
	}
}

func TuiMode(logger *log.Logger, opts *cache.Cache[types.Options]) error {
	logger.Println("Server Mode: Interactive")
	return nil
}
