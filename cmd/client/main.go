package main

/*
NOTES:
	- If Network Send channel is full the next message will lock up the program
*/

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/BurntSushi/toml"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/cache"

	"github.com/maplelm/dwarfwars/cmd/client/scenes/mainmenu"
)

var (
	configPath *string = flag.String("c", "./config/", "location of settings files")
	profile    *bool   = flag.Bool("profile", false, "profile program to file")
)

func main() {
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

	// Initializing Cachable settings from toml file
	var opts *cache.Cache[types.Options] = cache.New(time.Duration(5)*time.Second, func(o *types.Options) error {
		if o == nil {
			return fmt.Errorf("Options pointer can not be nil")
		}
		b, err := os.ReadFile(filepath.Join(*configPath, "General.toml"))
		if err != nil {
			return err
		}
		return toml.Unmarshal(b, o)
	})

	// Initializing the game engine
	g := game.New(opts.MustGet().General.ScreenWidth, opts.MustGet().General.ScreenHeight, "Dwarf Wars", opts, 1, []game.Scene{mainmenu.New()})

	// Setting Target FPS
	rl.SetTargetFPS(60)

	// Checking if settings specify fullscreen
	if opts.MustGet().General.Fullscreen {
		rl.ToggleFullscreen()
	}

	// Start the game
	g.Run()

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
