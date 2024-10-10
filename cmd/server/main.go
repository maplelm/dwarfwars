package main

import (
	// STD Packages
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	// 3rd Party Packages
	"github.com/BurntSushi/toml"

	// Project Packages
	"github.com/maplelm/dwarfwars/pkg/cache"
)

/*
 * Flags
 */
var (
	configPath *string = flag.String("c", "./config/", "location of settings files")
	//savepath   *string = flag.String("w", "./saves/World/", "location of world/game data")
	headless *bool = flag.Bool("h", false, "server will not use a tui and can be automated with scripts")
)

/*
 * Program Entry Point
 */
func main() {
	// Getting command line arguments
	flag.Parse()

	/*
	 * Variables
	 */
	var (
		err       error
		waitgroup *sync.WaitGroup       = new(sync.WaitGroup)
		opts      *cache.Cache[Options] = cache.New(time.Duration(5)*time.Second, func(o *Options) error {
			if o == nil {
				return fmt.Errorf("Options pointer can not be nil")
			}
			fullpath := filepath.Join(*configPath, "General.toml")
			b, err := os.ReadFile(fullpath)
			if err != nil {
				return err
			}
			//fmt.Printf("General Settings Data from file: %s\n", string(b))
			err = toml.Unmarshal(b, o)
			//b, err = toml.Marshal(o)
			//fmt.Printf("General Setting Data Marshalled: %s\n", string(b))
			return err
		})
	)

	/*
	 * Documentating Panics as they happen and then closing program.
	 */
	defer func() {
		if recover() != nil {
			fmt.Printf("Panicing: %s\n", recover())
			os.Exit(1)
		}
	}()

	/*
	 * Configuring Logging
	 */
	logflags := 0
	if opts.MustGet().Logging.Flags.UTC {
		logflags = logflags | log.LUTC
	}
	if opts.MustGet().Logging.Flags.Date {
		logflags = logflags | log.Ldate
	}
	if opts.MustGet().Logging.Flags.Time {
		logflags = logflags | log.Ltime
	}
	if opts.MustGet().Logging.Flags.Longfile {
		logflags = logflags | log.Llongfile
	}
	if opts.MustGet().Logging.Flags.Msgprefix {
		logflags = logflags | log.Lmsgprefix
	}
	if opts.MustGet().Logging.Flags.Shortfile {
		logflags = logflags | log.Lshortfile
	}
	if opts.MustGet().Logging.Flags.Microseconds {
		logflags = logflags | log.Lmicroseconds
	}
	MainLogger := log.New(os.Stdout, opts.MustGet().Logging.Prefix, logflags)

	/*
	 * Validating the SQL Database
	 */
	MainLogger.Println("Validating Database Before Server Bootup")
	if err = ValidateSQL(3, 500, MainLogger, opts); err != nil {
		MainLogger.Fatalf("Failed to Validate SQL Server: %s", err)
	}

	/*
	 * Creating the main application context.
	 */
	ctx, close := context.WithCancel(context.Background())

	/*
	 * Start the server based on what value the headless flag has
	 */
	switch *headless {
	case true:
		MainLogger.Println("Server Mode: Headless")
		go CliMode(MainLogger, ctx, waitgroup, opts)
		for {
			fmt.Printf("Dwarf Wars Server: ")
			var line string
			_, err := fmt.Scanln(&line)
			if err != nil {
				MainLogger.Printf("Failed to Read user input: %s", err)
			}
			switch line {
			case "stop", "quit":
				close()
				waitgroup.Wait()
				return
			case "ls":
				// list Connections
			default:
				fmt.Printf("Invalid input (%s)\n", line)
			}
		}
	case false:
		MainLogger.Println("Server Mode: Interactive")
		TuiMode(MainLogger, ctx, waitgroup, opts)
		close()
		waitgroup.Wait()
		return
	}

}

func TuiMode(logger *log.Logger, ctx context.Context, wgrp *sync.WaitGroup, opts *cache.Cache[Options]) error {
	return nil
}

func CliMode(logger *log.Logger, ctx context.Context, wgrp *sync.WaitGroup, opts *cache.Cache[Options]) error {
	var (
		err  error
		addr *net.TCPAddr
	)

	/*
	 *	Adding CliMode to the sync group.
	 *	This should allow for proper shutdowns.
	 */
	wgrp.Add(1)
	defer wgrp.Done()

	/*
	*	Resolving TCP Address for the Server to use.
	*	Required, will shut the server down if failed.
	 */
	logger.Printf("Resolving TCP Address: %s:%d", opts.MustGet().Game.Addr, opts.MustGet().Game.Port)
	if addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", opts.MustGet().Game.Addr, opts.MustGet().Game.Port)); err != nil {
		logger.Printf("Failed to Resolve TCP Address for server: %s", err)
		return err
	} else {
		logger.Printf("Resolved Server Address: %s", addr.String())
	}

	server, err := NewServer(addr, 10, ConnectionHandlerFunc(CommandHandlerTest))
	if err != nil {
		logger.Printf("Failed to Create Server Object: %s", err)
		return err
	}

	logger.Printf("Starting Dwarf Wars Server...")
	return server.Start(opts, logger, nil, ctx)
}
