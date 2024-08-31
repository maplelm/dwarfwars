package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"

	"github.com/maplelm/dwarfwars/pkg/logging"
	"github.com/maplelm/dwarfwars/pkg/settings"
	"github.com/maplelm/dwarfwars/pkg/tui"
)

func main() {
	var (
		opts         *settings.Config = nil
		err          error            = nil
		appWaitGroup sync.WaitGroup
	)

	// CLI Flags
	var (
		tuimode  *bool   = flag.Bool("tui", false, "If true the server will allow the user to interact with features through a TUI interface")
		optsPath *string = flag.String("path", "", "path to to the main settings TOML file")
		optsName *string = flag.String("name", "", "name of main settings TOML file")
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
		Add("Start Server", false, &ServerCommand{Server: nil, waitgroup: mainWaitGroup}).
		AddFunc("Settings", false, EditSettings).
		AddFunc("Quit", false, Quit)
	if _, err := tea.NewProgram(mainMenu).Run(); err != nil {
		logging.Error(err, "Game Server, Main Thread, TUI Mode, Bubbletea Crashed")
		os.Exit(1)
	}
}
func cliMode() {}

func LoadSettings(path, name string) (opts *settings.Config, err error) {
	/*
		Priority:
			+ Commandline Arg
			+ Environment Variable
			+ Default Value
	*/
	if len(path) == 0 {
		path = os.Getenv("SETTINGS_PATH")
		if len(path) == 0 {
			path = "./"
		}
	}
	if len(name) == 0 {
		name = os.Getenv("SETTINGS_NAME")
		if len(name) == 0 {
			name = "settings.toml"
		}
	}
	if _, err := settings.LoadFromTomlFile("Main", path, name); err != nil {
		return nil, err
	}
	return settings.Get[settings.Config]("Main")
}

func InitLogger(opts *settings.Config) (err error) {
	fmt.Printf("Creating logger with path: %s and file name: %s\n", opts.Log.Path, opts.Log.FileName)
	log.SetOutput(logging.NewRotationWriter(opts.Log.Path, opts.Log.FileName, opts.Log.MaxFileSize))
	log.SetPrefix("Game Server:")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// Validating that log path exists
	if _, err = os.Stat(opts.Log.Path); err != nil && errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(opts.Log.Path, 0777)
	}
	return
}
