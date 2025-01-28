package main

import (
	// STD Packages
	"context"
	"flag"
	"fmt"
	"os"
	fp "path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	// Third Party Packages
	"github.com/BurntSushi/toml"
	zl "github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	ljack "gopkg.in/natefinch/lumberjack.v2"

	// Project Packages
	"server/internal/cache"
	"server/internal/types"

	// Drivers
	_ "github.com/mattn/go-sqlite3"
)

// ///////////////////////
// CLI Argument Parsing //
// ///////////////////////
var (
	configPath *string = flag.String("c", "./config/", "location of settings files")
	headless   *bool   = flag.Bool("h", false, "server will not use a tui and can be automated with scripts")
	profile    *bool   = flag.Bool("profile", false, "Export Profiling information")
)

// ///////////////////
// Global Variables //
// ///////////////////
var ()

// Main Function
func main() {

	var (
		MLog zl.Logger

		RuntimeCtx    context.Context
		CancelRuntime func()

		Opts cache.Cache[types.Options]
	)

	///////////////////////////
	// Parsing CLI Arguments //
	///////////////////////////
	flag.Parse()

	///////////////////////////////
	// Initialize Settings Cache //
	///////////////////////////////
	Opts = OptsSetup()

	/////////////////////////
	// Configuring Logging //
	/////////////////////////
	{
		zl.SetGlobalLevel(zl.TraceLevel)
		zl.TimeFieldFormat = zl.TimeFormatUnix
		zl.ErrorStackMarshaler = pkgerrors.MarshalStack

		ld := Opts.MustGetData().Logging
		MLog = zl.New(&ljack.Logger{
			MaxSize:    ld.MaxSize,     // megabytes
			MaxAge:     ld.MaxAge,      //days
			Compress:   ld.Compression, // disabled by default
			Filename:   fp.Join(ld.Path, ld.Name),
			MaxBackups: ld.Backups,
		})
	}

	//////////////////////////////////
	// Configuring Runetime Context //
	//////////////////////////////////
	RuntimeCtx, CancelRuntime = context.WithCancel(context.Background())
	defer CancelRuntime()

	///////////////////
	// CPU Profiling //
	///////////////////
	if *profile {
		MLog.Info().Msg("Starting CPU Profiling")
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

	//////////////////
	// Start Server //
	//////////////////
	if *headless {
		if err := CliMode(RuntimeCtx, &MLog, &Opts); err != nil {
			MLog.Err(err).Msg("Server ran into an error in CLI Mode")
		}
	} else {
		if err := TuiMode(RuntimeCtx, &MLog, &Opts); err != nil {
			MLog.Err(err).Msg("Server ran into an error in TUI Mode")
		}
	}

	//////////////
	// Clean Up //
	//////////////

	///////////////////////////
	// Memory Profiling Code //
	///////////////////////////
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

func OptsSetup() cache.Cache[types.Options] {
	r, e := cache.New(time.Second*time.Duration(5), func(o *types.Options) error {
		if o == nil {
			return fmt.Errorf("Options pointer can not be nil")
		}
		f := fp.Join(*configPath, "General.toml")
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		return toml.Unmarshal(b, o)

	})
	if e != nil {
		panic(e)
	}
	return *r
}
