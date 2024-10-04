package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/maplelm/dwarfwars/pkg/cache"

	_ "github.com/go-sql-driver/mysql"
)

var (
	headless   *bool                 = flag.Bool("headless", false, "server will not use a tui and can be automated with scripts")
	configPath *string               = flag.String("config", "./config/", "location of settings files")
	savepath   *string               = flag.String("world", "./saves/World/", "location of world/game data")
	opts       *cache.Cache[Options] = cache.New(time.Duration(5)*time.Second, func(o *Options) error {
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
	sqlcreds *cache.Cache[struct {
		user     string
		password string
	}] = cache.New(time.Duration(5)*time.Second, func(c *struct {
		user     string
		password string
	}) error {
		o, err := opts.Get()
		if err != nil {
			return err
		}
		c = &struct {
			user     string
			password string
		}{
			user:     o.Db.Username,
			password: o.Db.Password,
		}
		return nil
	})
)

func main() {
	var (
		o         *Options
		err       error
		waitgroup *sync.WaitGroup = new(sync.WaitGroup)
	)
	// Getting command line arguments
	flag.Parse()

	logflags := 0
	if o, err = opts.Get(); err != nil {
		panic(err)
	} else {
		if o.Logging.Flags.UTC {
			logflags = logflags | log.LUTC
		}
		if o.Logging.Flags.Date {
			logflags = logflags | log.Ldate
		}
		if o.Logging.Flags.Time {
			logflags = logflags | log.Ltime
		}
		if o.Logging.Flags.Longfile {
			logflags = logflags | log.Llongfile
		}
		if o.Logging.Flags.Msgprefix {
			logflags = logflags | log.Lmsgprefix
		}
		if o.Logging.Flags.Shortfile {
			logflags = logflags | log.Lshortfile
		}
		if o.Logging.Flags.Microseconds {
			logflags = logflags | log.Lmicroseconds
		}
	}

	// Setting up logging
	MainLogger := log.New(os.Stdout, o.Logging.Prefix, logflags)

	// validate the sql server here
	MainLogger.Println("Validating Database Before Server Bootup")
	sqlvalidationattempts := 0
sqlvalidation:
	for sqlvalidationattempts < 3 {
		var (
			creds *struct {
				user     string
				password string
			}
			o    *Options
			err  error
			conn *sql.DB
		)
		if creds, err = sqlcreds.Get(); err != nil {
			panic(err)
		}
		if o, err = opts.Get(); err != nil {
			panic(err)
		}

		if conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s", creds.user, creds.password, o.Db.Addr, o.Db.Port, "")); err != nil {
			sqlvalidationattempts++
			MainLogger.Printf("Failed to connect to sql server, waiting %d seconds and trying again\n\t%s", sqlvalidationattempts*3, err)
			time.Sleep(time.Duration(sqlvalidationattempts*3) * time.Second)
			continue sqlvalidation
		}
		defer conn.Close()

		MainLogger.Printf("Walking %s to get sql validation scripts", o.Db.ValidationDir)
		if err = filepath.Walk(o.Db.ValidationDir, func(path string, info os.FileInfo, err error) error {
			var b []byte

			if err != nil {
				return err
			}
			MainLogger.Printf("Reading SQL file: %s", path)
			if b, err = os.ReadFile(path); err != nil {
				return nil
			}

			_, err = conn.Exec(string(b))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			sqlvalidationattempts++
			MainLogger.Printf("Failed to run SQL validation script %s, Waiting %d seconds and trying again\n\t%s", err, sqlvalidationattempts*3, err)
			time.Sleep(time.Duration(sqlvalidationattempts*3) * time.Second)
			continue sqlvalidation
		}
	}
	if sqlvalidationattempts >= 3 {
		MainLogger.Fatalln(fmt.Errorf("Failed to validate sql server before game server boots up.\n\t%s", err))
	}

	ctx, close := context.WithCancel(context.Background())

	switch *headless {
	case true:
		MainLogger.Println("Server Mode: Headless")
		go CliMode(MainLogger, ctx, waitgroup)
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
		TuiMode(MainLogger)
		close()
		waitgroup.Wait()
	}

}

func TuiMode(logger *log.Logger) error {
	return nil
}

func CliMode(logger *log.Logger, ctx context.Context, wgrp *sync.WaitGroup) error {
	var (
		o        *Options
		err      error
		addr     *net.TCPAddr
		listener *net.TCPListener
		connChan chan net.Conn = make(chan net.Conn, 10)
	)

	wgrp.Add(1)
	defer wgrp.Done()

	if o, err = opts.Get(); err != nil {
		logger.Printf("Failed to get options: %s", err)
		return err
	}

	if addr, err = net.ResolveTCPAddr("tcp", o.Game.Addr); err != nil {
		logger.Printf("Failed to Resolve TCP Address for server: %s", err)
		return err
	}

	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		logger.Printf("Failed to creat TCP Listner: %s", err)
		return err
	}

	go func(cc chan net.Conn) {
		wgrp.Add(1)
		defer wgrp.Done()
		for {
			if conn, err := listener.AcceptTCP(); err != nil {
				if errors.Is(err, net.ErrClosed) {
					logger.Printf("TCP Listener closed: %s", err)
					return
				}

				var netErr *net.OpError
				if errors.As(err, &netErr) && netErr.Timeout() {
					logger.Printf("Listener timed out before accepting a connection: %s", err)
				} else {
					logger.Printf("Listener Failed to Accept Incoming Connection: %s", err)
				}
			} else {
				conn.SetReadDeadline(time.Now().Add(time.Duration(o.Game.Timeouts.Read) * time.Millisecond))
				conn.SetWriteDeadline(time.Now().Add(time.Duration(o.Game.Timeouts.Write) * time.Millisecond))
				cc <- conn
			}
		}
	}(connChan)

	for {
		select {
		case <-ctx.Done():
			return listener.Close()
		case conn := <-connChan:
			// handle each connection has they come in from the listner
			go func(c net.Conn) error {
				wgrp.Add(1)
				defer wgrp.Done()
				defer c.Close()

				var (
					data []byte
					n    int
					err  error
				)

				if n, err = c.Read(data); n > 0 && err == nil {
					logger.Printf("Message from connection (%s): %s", conn.RemoteAddr(), string(data))
				} else if err != nil {
					logger.Printf("Failed to read data from client: %s", err)
					data = []byte(err.Error())
				}

				if n, err = c.Write(data); n > 0 && err == nil {
					logger.Printf("Message Sent to Connection (%s): %s", conn.RemoteAddr(), string(data))
				} else if err != nil {
					logger.Printf("Failed to Write data to Client: %s", err)
				}
				return nil
			}(conn)
		}
	}
}
