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
	"time"

	"github.com/BurntSushi/toml"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/game"
	"github.com/maplelm/dwarfwars/cmd/client/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/cache"
)

var configPath *string = flag.String("c", "./config/", "location of settings files")

func main() {
	flag.Parse()
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
	g := game.New(opts.MustGet().General.ScreenWidth, opts.MustGet().General.ScreenHeight, "Dwarf Wars", opts, 1, []game.Scene{&MainMenu{}})
	rl.SetTargetFPS(60)
	if opts.MustGet().General.Fullscreen {
		rl.ToggleFullscreen()
	}
	g.Run()
}
