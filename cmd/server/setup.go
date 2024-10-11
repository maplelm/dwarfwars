package main

import (
	// STD Packages
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	// Project Packages
	"github.com/maplelm/dwarfwars/pkg/cache"

	// Drivers
	_ "github.com/go-sql-driver/mysql"
)

func ValidateSQL(maxAttempts, timeoutRate int, logger *log.Logger, opts *cache.Cache[Options]) error {

	var (
		o *Options
		//creds *Credentials
		conn *sql.DB
		err  error
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
}
