package main

/*
	FIX: Program Locks up if you try and disable server after starting it.
*/

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/maplelm/dwarfwars/pkg/server"
	"github.com/maplelm/dwarfwars/pkg/settings"
	"github.com/maplelm/dwarfwars/pkg/tui"
	"log"
	"os"
	"sync"
)

func main() {

	var (
		waitgroup sync.WaitGroup
		serv      *server.Server = nil
	)

	/////////////////////////////
	// Loading System Settings //
	/////////////////////////////
	settingsPath := os.Getenv("SETTINGS_PATH")
	if len(settingsPath) == 0 {
		settingsPath = "./"
	}
	settingsName := os.Getenv("SETTINGS_NAME")
	if len(settingsName) == 0 {
		settingsName = "settings.toml"
	}

	_, err := settings.LoadFromTomlFile("Main", settingsPath, settingsName)
	if err != nil {
		log.Fatalf("(Main Thread) Failed to load main TOML settings file, %s", err)
	}

	opts, err := settings.Get[settings.Config]("Main")
	if err != nil {
		log.Fatalf("(Main Thread) Failed to get Main settings from memory, %s", err)
	}

	//////////////////////////////
	// Setting up System Logger //
	//////////////////////////////
	fmt.Printf("Creating logger with path: %s and file name: %s\n", opts.Log.Path, opts.Log.FileName)
	log.SetOutput(NewRotationWriter(
		opts.Log.Path,
		opts.Log.FileName,
		opts.Log.AdjustedPollRate(),
		opts.Log.MaxFileSize),
	)
	log.SetPrefix("System:")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Validating that log path exists
	_, err = os.Stat(opts.Log.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(opts.Log.Path, 0777)
			if err != nil {
				fmt.Printf("Failed to setup Logging Path, %s\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Failed to validate Log Path, %s\n", err)
			os.Exit(1)
		}
	}

	//////////////////////////////
	// Setting Up TUI / Actions //
	//////////////////////////////
	mainMenu := tui.NewMenu('>', "Dwarf Wars Server", lg.NewStyle(), lg.NewStyle(), lg.NewStyle())
	mainMenu = mainMenu.
		Add("Start Server", false, func(state bool) (cmd tea.Cmd, s bool, err error) {
			if !state {
				log.Printf("(Dwarf Wars Server) Starting Server")
				serv, err = server.New(opts.Server.Addr, fmt.Sprintf("%d", opts.Server.Port))
				if err != nil {
					log.Printf("(Dwarf Wars Server) Failed to Create server, %s", err)
					return tea.Quit, false, fmt.Errorf("Main Menu Start Server: %s", err)
				}
				go func() {
					waitgroup.Add(1)
					defer waitgroup.Done()
					serv.Start()
				}()
				s = true
			} else {
				serv.Stop()
				log.Printf("(TUI) Stopped Server")
				s = false
			}
			return
		}).
		Add("Settings", false, func(state bool) (cmd tea.Cmd, s bool, err error) {
			var inter interface{} = *opts
			settingsMenu := tui.NewSettingsMenu(&inter, '>', lg.NewStyle(), lg.NewStyle(), lg.NewStyle())
			p := tea.NewProgram(settingsMenu)
			log.Printf("(TUI) Opening settings window")
			if _, err = p.Run(); err != nil {
				log.Printf("(TUI) Failed to open settings window, %s", err)
				return nil, state, fmt.Errorf("Settings Menu: %s", err)
			}
			s = state
			return
		}).
		Add("Quit", false, func(state bool) (cmd tea.Cmd, s bool, err error) {
			log.Printf("(TUI) Quiting Program")
			return tea.Quit, state, nil
		})
	op := tea.NewProgram(mainMenu)
	if _, err := op.Run(); err != nil {
		log.Fatalf("(Bubbletea) Error, %s", err)
	}
}
