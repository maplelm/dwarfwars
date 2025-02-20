package main

import (
	// STD Packages
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	fp "path/filepath"
	"strings"
	"time"

	// 3rd Party
	"github.com/BurntSushi/toml"
	zl "github.com/rs/zerolog"

	// Project Packages
	"github.com/maplelm/dwarfwars/pkg/cache"
	"github.com/maplelm/dwarfwars/src/server/pkg/types"

	// Drivers
	_ "github.com/go-sql-driver/mysql"
)

/*
# ValidateSQL

This fucntion valdiates that the execpted Mariadb sql database is working properly and has the exected databases and tables before the server offically starts up. If the database is not in working order the server will not work.

## Variables

- maxAttempts (int): Number of times server will try and revalidate the SQL database before giving up.
- timeoutRate (int): The base rate that the server will wait inbetween each validation attempt.
- logger (*log.Logger): used to log function activity
- opts (*cache.Cache[types.Options]): struct of settings the server needs to operate.

## Returns

- Error
*/
func ValidateSQL(ctx context.Context, log zl.Logger, opts *cache.Cache[types.Options]) error {

	var (
		o    types.Options = opts.MustGetData()
		conn *sql.DB
		err  error
	)

	// Getting List of Databases to check
	b, err := os.ReadFile("./config/databases.csv")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get list of databases to validate")
	}
	dbSegs := strings.Split(string(b), "\n")
	var dbList []string = []string{}
	for _, s := range dbSegs {
		dbList = append(dbList, strings.Split(s, ",")...)
	}

	///////////////////////////////////////
	// Making sure databases are created //
	///////////////////////////////////////
	for _, e := range dbList {
		conn, err = sql.Open("sqlite", "file:"+fp.Join(o.Db.BaseDir, dbList[0]+".db")+"?mode=rw&_mutex=full")
		if err != nil {
			log.Fatal().Err(err).Str("Database", e).Msg("Failed to create a database during the validation process")
		}
		conn.Close()
	}

	/////////////////////////////////////
	// Connecting to all the databases //
	/////////////////////////////////////
	for i, e := range dbList {
		if i == 0 {
			conn, _ = sql.Open("sqlite3", "file:"+fp.Join(o.Db.BaseDir, e+".db")+"?mode=rw&_mutex=full")
		} else {
			_, err = conn.Exec("attach database '" + fp.Join(o.Db.BaseDir, e+".db") + "' as " + e + ";")
			if err != nil {
				log.Fatal().Err(err).Str("Database", e).Msg("failed to attach database to connection")
			}
		}
	}

	return nil

	/*
	   	for a := 0; a < maxAttempts; a++ {
	   		if o, err = opts.Get(); err != nil {
	   			if a <= (maxAttempts - 1) {
	   				return err
	   			}
	   			logger.Printf("Failed to Pull options for Validating SQL Server, Waiting %d milliseconds...", (a+1)*timeoutRate)
	   			continue
	   		}
	   		if conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s", o.Db.Username, o.Db.Password, o.Db.Addr, o.Db.Port, "")); err != nil {
	   			if a <= (maxAttempts - 1) {
	   				return err
	   			}
	   			logger.Printf("Failed to connect to SQL Server for validation, Waiting %d milliseconds...", (a+1)*timeoutRate)
	   			time.Sleep(time.Duration((a+1)*timeoutRate) * time.Millisecond)
	   		}
	   		defer conn.Close()

	   		if len(o.Db.ValidationDir) == 0 {
	   			logger.Fatalf("SQL Validation: Validation Script Directory is empty")
	   		}
	   		logger.Printf("SQL Validation: Walking Dir: %s", o.Db.ValidationDir)
	   		if err = filepath.Walk(o.Db.ValidationDir, func(path string, info os.FileInfo, err error) error {
	   			if err != nil {
	   				logger.Printf("SQL Validation filepath.Walk Error: %s", err)
	   				return err
	   			}
	   			if info.IsDir() {
	   				return nil
	   			}
	   			var b []byte
	   			logger.Printf("Reading SQL File: %s", path)
	   			if b, err = os.ReadFile(path); err != nil {
	   				return err
	   			}
	   			for _, v := range strings.Split(string(b), ";") {
	   				if len(strings.TrimSpace(v)) > 0 {
	   					_, err = conn.Exec(v)
	   					if err != nil {
	   						return err
	   					}
	   				}
	   			}
	   			return err
	   		}); err != nil {
	   			if a <= (maxAttempts - 1) {
	   				return err
	   			}
	   			logger.Printf("Failed to Run SQL Scripts to Validate SQL Server, Waiting %d milliseconds...", (a+1)*timeoutRate)
	   			time.Sleep(time.Duration((a+1)*timeoutRate) * time.Millisecond)
	   			continue
	   		}

	   }
	   return nil
	*/
}

func InitLogger(opts *cache.Cache[types.Options]) *log.Logger { /*
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
	return log.New(os.Stdout, opts.MustGet().Logging.Prefix, logflags)
}

func InitOptionsCache(configpath string) *cache.Cache[types.Options] {
	return cache.New(time.Duration(5)*time.Second, func(o *types.Options) error {
		if o == nil {
			return fmt.Errorf("Options pointer can not be nil")
		}
		fullpath := filepath.Join(configpath, "General.toml")
		b, err := os.ReadFile(fullpath)
		if err != nil {
			return err
		}
		return toml.Unmarshal(b, o)
	})
}
