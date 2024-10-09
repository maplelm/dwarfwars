package main

import (
	// STD Packages
	"context"
	"errors"
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
 * Global variables
 */
var ()

/*
 * Program Entry Point
 */
func main() {
	/*
	 * Flags
	 */
	var (
		configPath *string = flag.String("config", "./config/", "location of settings files")
		savepath   *string = flag.String("world", "./saves/World/", "location of world/game data")
		headless   *bool   = flag.Bool("headless", false, "server will not use a tui and can be automated with scripts")
	)

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
			return toml.Unmarshal(b, o)
		})
		sqlcreds *cache.Cache[Credentials] = cache.New(time.Duration(5)*time.Second, func(c *Credentials) error {
			o, err := opts.Get()
			if err != nil {
				return err
			}
			c = &Credentials{
				Username: o.Db.Username,
				Password: o.Db.Password,
			}
			return nil
		})
		connectionsMutex sync.Mutex
		connections      map[net.Addr]*net.Conn
	)
	// Getting command line arguments
	flag.Parse()

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
	if err = ValidateSQL(3, 500, MainLogger, opts, sqlcreds); err != nil {
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
				break
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
	if addr, err = net.ResolveTCPAddr("tcp", opts.MustGet().Game.Addr); err != nil {
		logger.Printf("Failed to Resolve TCP Address for server: %s", err)
		return err
	}

	server, err := NewServer(addr, 10)
	if err != nil {
		logger.Printf("Failed to Create Server Object: %s", err)
		return err
	}

	return server.Start(opts, logger, nil, ctx)
}
