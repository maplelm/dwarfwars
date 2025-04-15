package main

import (
	// STD Packages
	"context"
	"database/sql"
	"fmt"
	"os"
	fp "path/filepath"
	"strings"
	"time"

	// 3rd Party
	"github.com/BurntSushi/toml"
	zl "github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	ljack "gopkg.in/natefinch/lumberjack.v2"

	// Project Packages
	"server/internal/cache"
	"server/internal/types"
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

	//////////////////////////////////////////
	// Creating Database Connection Objects //
	//////////////////////////////////////////
	for _, e := range dbList {
		conn, err = sql.Open("sqlite", "file:"+fp.Join(o.Db.BaseDir, dbList[0]+".db")+"?mode=rw&_mutex=full")
		if err != nil {
			log.Fatal().Err(err).Str("Database", e).Msg("Failed to create a database during the validation process")
		}
		defer conn.Close()
	}

	//////////////////////////////////////
	// Testing Connections to Databases //
	//////////////////////////////////////
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
}

func InitLogger(opts *cache.Cache[types.Options], logger **zl.Logger) {

	// Configuring ZeroLog //
	zl.SetGlobalLevel(zl.TraceLevel)
	zl.TimeFieldFormat = zl.TimeFormatUnix
	zl.ErrorStackMarshaler = pkgerrors.MarshalStack

	*logger = zl.New(&ljack.Logger{
		MaxSize:    opts.MustGetData().Logging.MaxSize,     // Megabytes
		MaxAge:     opts.MustGetData().Logging.MaxAge,      // days
		Compress:   opts.MustGetData().Logging.Compression, // disabled by default
		filename:   fp.Join(opts.MustGetData().Logging.Path, opts.MustGetData().Logging.Name),
		MaxBackups: opts.MustGetData().Logging.Backups,
	})
}

func InitOptionsCache(configpath string) (*cache.Cache[types.Options], error) {
	return cache.New(time.Duration(5)*time.Second, func(o *types.Options) error {
		if o == nil {
			return fmt.Errorf("Options pointer can not be nil")
		}
		fullpath := fp.Join(configpath, "General.toml")
		b, err := os.ReadFile(fullpath)
		if err != nil {
			return err
		}
		return toml.Unmarshal(b, o)
	})
}
