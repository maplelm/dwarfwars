package main

import (
	// STD Packages
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	// Project Packages
	"server/internal/cache"
	"server/internal/server"
	"server/internal/types"
)

// Main Function
func main() {
	var (
		configPath *string = flag.String("c", "./config/", "location of settings files")
		headless   *bool   = flag.Bool("h", false, "server will not use a tui and can be automated with scripts")
		profile    *bool   = flag.Bool("profile", false, "Export Profiling information")
	)
	flag.Parse()

	// CPU Profiling Code
	if *profile {
		fc, err := os.Create("cpuprofile.txt")
		if err != nil {
			panic(err)
		}
		defer fc.Close()

		if err := pprof.StartCPUProfile(fc); err != err {
			panic(err)
		}
		defer pprof.StopCPUProfile()

	}

	// Getting Settings from TOML file
	opts := InitOptionsCache(*configPath)

	// Initializing the Logger object
	MainLogger := InitLogger(opts)
	MainLogger.Printf("Starting Logging...")

	// Validating the SQL Server
	MainLogger.Println("Validating Database Before Server Bootup")
	if err := ValidateSQL(3, 500, MainLogger, opts); err != nil {
		MainLogger.Fatalf("Failed to Validate SQL Server: %s", err)
	}

	// Start Server
	switch *headless {
	case true:
		if err := CliMode(MainLogger, opts); err != nil {
			MainLogger.Printf("Server Error: %s", err)
		}
	case false:
		if err := TuiMode(MainLogger, opts); err != nil {
			MainLogger.Printf("Server Error: %s", err)
		}
	}

	// Memory Profiling Code
	if *profile {
		runtime.GC()
		fm, err := os.Create("memoryprofile.txt")
		if err != nil {
			panic(err)
		}
		defer fm.Close()
		if err := pprof.WriteHeapProfile(fm); err != nil {
			panic(err)
		}
	}
}

func TuiMode(logger *log.Logger, opts *cache.Cache[types.Options]) error {
	logger.Println("Server Mode: Interactive")
	return nil
}
