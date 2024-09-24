package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/maplelm/dwarfwars/cmd/gameserver/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"

	"github.com/maplelm/dwarfwars/pkg/logging"
	"github.com/maplelm/dwarfwars/pkg/settings"
	"github.com/maplelm/dwarfwars/pkg/tui"
)

// Flags
var (
	mainSettingsPath *string
	mainSettingsName *string
)

func new_main() {
	p := tea.NewProgram(tui.InitMenu('>', []tui.Option{}, "Game Server"))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error Running Program: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	new_main()
	return
	var (
		opts *settings.CachedValue[settings.Settings] = settings.NewCachedValue(time.Duration(5)*time.Minute, func(c *settings.CachedValue[settings.Settings]) error {
			if mainSettingsPath == nil {
				return fmt.Errorf("mainSettingsPath must be initialized before cached value can dereference it")
			}
			if mainSettingsName == nil {
				return fmt.Errorf("mainSettingsName must be initialized before cached value can be dereferenced")
			}
			value, err := settings.LoadFromFile[settings.Settings](*mainSettingsPath, *mainSettingsName)
			if err != nil {
				return err
			}
			c.Set(value)
			return nil
		})
		err          error = nil
		appWaitGroup sync.WaitGroup
	)

	// CLI Flags
	var (
		tuimode          *bool = flag.Bool("tui", false, "If true the server will allow the user to interact with features through a TUI interface")
		mainSettingsPath       = flag.String("path", "./config/", "path to to the main settings TOML file")
		mainSettingsName       = flag.String("name", "mainSettings.toml", "name of main settings TOML file")
	)
	flag.Parse()

	if opts, err = LoadSettings(*optsPath, *optsName); err != nil {
		logging.Error(err, "Game Server Main Thread Loading Settings")
		os.Exit(1)
	}

	err = InitLogger(opts)
	if err != nil {
		logging.Error(err, "Game Server Main Thread Initializing Logger")
		os.Exit(1)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	if *tuimode {
		tuiMode(ctx, &appWaitGroup)
	} else {
		cliMode()
	}

	ctxCancel()
	appWaitGroup.Wait()
}

func tuiMode(mainCtx context.Context, mainWaitGroup *sync.WaitGroup) {
	mainMenu := tui.NewMenu('>', "Dwarf Wars Server", lg.NewStyle(), lg.NewStyle(), lg.NewStyle(), mainCtx).
		Add("Start Server", false, &ServerSelecter{Server: nil, waitgroup: mainWaitGroup}).
		AddFunc("Settings", false, EditSettings).
		AddFunc("Quit", false, Quit)
	if _, err := tea.NewProgram(mainMenu).Run(); err != nil {
		logging.Error(err, "Game Server, Main Thread, TUI Mode, Bubbletea Crashed")
		os.Exit(1)
	}
}
func cliMode() {}

func InitLogger(opts settings.Config) (err error) {
	fmt.Printf("Creating logger with path: %s and file name: %s\n", opts.Log.Path, opts.Log.FileName)
	log.SetOutput(logging.NewRotationWriter(opts.Log.Path, opts.Log.FileName, opts.Log.MaxFileSize))
	log.SetFlags(log.LUTC)
	if _, err = os.Stat(opts.Log.Path); err != nil && errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(opts.Log.Path, 0777)
	}
	return
}
