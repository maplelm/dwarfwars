package main

import (
	// STD Packages
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	// Project Packages
	"github.com/maplelm/dwarfwars/pkg/cache"

	// Drivers
	_ "github.com/go-sql-driver/mysql"
)

func ValidateSQL(maxAttempts, timeoutRate int, logger *log.Logger, opts *cache.Cache[Options], servCreds *cache.Cache[Credentials]) error {

	var (
		o     *Options
		creds *Credentials
		conn  *sql.DB
		err   error
	)

	for a := 0; a < maxAttempts; a++ {
		if o, err = opts.Get(); err != nil {
			if a <= (maxAttempts - 1) {
				return err
			}
			logger.Printf("Failed to Pull options for Validating SQL Server, Waiting %d milliseconds...", (a+1)*timeoutRate)
			time.Sleep(time.Duration((a+1)*timeoutRate) * time.Millisecond)
			continue
		}
		if creds, err = servCreds.Get(); err != nil {
			if a <= (maxAttempts - 1) {
				return err
			}
			time.Sleep(time.Duration((a+1)*timeoutRate) * time.Millisecond)
			continue
		}
		if conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s", creds.Username, creds.Password, o.Db.Addr, o.Db.Port, "")); err != nil {
			if a <= (maxAttempts - 1) {
				return err
			}
			logger.Printf("Failed to connect to SQL Server for validation, Waiting %d milliseconds...", (a+1)*timeoutRate)
			time.Sleep(time.Duration((a+1)*timeoutRate) * time.Millisecond)
		}
		defer conn.Close()

		if err = filepath.Walk(o.Db.ValidationDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			var b []byte
			logger.Printf("Reading SQL File: %s", path)
			if b, err = os.ReadFile(path); err != nil {
				return err
			}
			_, err = conn.Exec(string(b))
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
}
