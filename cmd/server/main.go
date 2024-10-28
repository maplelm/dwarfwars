package main

import (
	// STD Packages
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	// Project Packages
	s "github.com/maplelm/dwarfwars/cmd/server/pkg/server"
	"github.com/maplelm/dwarfwars/cmd/server/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/cache"
)

/*
 * Flags
 */
var (
	configPath *string = flag.String("c", "./config/", "location of settings files")
	//savepath   *string = flag.String("w", "./saves/World/", "location of world/game data")
	headless *bool = flag.Bool("h", false, "server will not use a tui and can be automated with scripts")
	server   *s.Server
)

/*
 * Program Entry Point
 */
func main() {
	var (
		err       error
		waitgroup *sync.WaitGroup = new(sync.WaitGroup)
	)

	flag.Parse()

	opts := InitOptionsCache()

	MainLogger := InitLogger(opts)

	MainLogger.Println("Validating Database Before Server Bootup")
	if err = ValidateSQL(3, 500, MainLogger, opts); err != nil {
		MainLogger.Fatalf("Failed to Validate SQL Server: %s", err)
	}

	ctx, close := context.WithCancel(context.Background())

	// Start the server based on what value the headless flag has
	switch *headless {
	case true:
		MainLogger.Println("Server Mode: Headless")
		CliMode(MainLogger, ctx, waitgroup, opts)
	case false:
		MainLogger.Println("Server Mode: Interactive")
		TuiMode(MainLogger, ctx, waitgroup, opts)
	}

	// Clean up
	MainLogger.Println("Server Shutting down...")
	close()
	waitgroup.Wait()
	return
}

func CliMode(logger *log.Logger, ctx context.Context, wgrp *sync.WaitGroup, opts *cache.Cache[types.Options]) error {
	var (
		err  error
		addr *net.TCPAddr
	)

	cliCtx, close := context.WithCancel(ctx)

	/*
	*	Resolving TCP Address for the Server to use.
	*	Required, will shut the server down if failed.
	 */
	logger.Printf("Resolving TCP Address: %s:%d", opts.MustGet().Game.Addr, opts.MustGet().Game.Port)
	if addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", opts.MustGet().Game.Addr, opts.MustGet().Game.Port)); err != nil {
		logger.Printf("Failed to Resolve TCP Address for server: %s", err)
		close()
		return err
	} else {
		logger.Printf("Resolved Server Address: %s", addr.String())
	}

	server, err = s.New(addr, 10)
	if err != nil {
		logger.Printf("Failed to Create Server Object: %s", err)
		close()
		return err
	}

	// Running the user interface thread for the server...
	go func() error {
		wgrp.Add(1)
		defer wgrp.Done()
		for {
			fmt.Printf("Dwarf Wars Server: ")
			var line string
			_, err := fmt.Scanln(&line)
			if err != nil {
				logger.Printf("Failed to Read user input: %s", err)
			}
			switch line {
			case "stop", "quit":
				close()
				return nil
			case "ls":
				// list Connections
			case "count":
				//fmt.Printf("Connections: %d\n", len(server.Connections))
			default:
				fmt.Printf("Invalid input (%s)\n", line)
			}
		}
	}()

	logger.Printf("Starting Dwarf Wars Server...")
	return server.Start(opts, logger, nil, cliCtx)
}

func TuiMode(logger *log.Logger, ctx context.Context, wgrp *sync.WaitGroup, opts *cache.Cache[types.Options]) error {
	return nil
}
