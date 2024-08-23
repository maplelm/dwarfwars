package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/maplelm/dwarfwars/pkg/server"
	"github.com/maplelm/dwarfwars/pkg/settings"
	"github.com/maplelm/dwarfwars/pkg/tui"
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

	////////////////////////////////////////////
	// Setting up Default Context for program //
	////////////////////////////////////////////
	main_ctx, ctx_cancel := context.WithCancel(context.Background())

	//////////////////////////////
	// Setting Up TUI / Actions //
	//////////////////////////////
	mainMenu := tui.NewMenu('>', "Dwarf Wars Server", lg.NewStyle(), lg.NewStyle(), lg.NewStyle())
	ctxServer, serverCancel := context.WithCancel(main_ctx)
	mainMenu = mainMenu.
		Add("Start Auth Server", false, func(state bool) (cmd tea.Cmd, s bool, err error) {
			if !state {

				serv, err = server.New(opts.Server.Addr, fmt.Sprintf("%d", opts.Server.Port), time.Duration(opts.Server.IdleTimeout)*time.Millisecond)
				if err != nil {
					log.Printf("(Dwarf Wars Server) Failed to Create server, %s", err)
					return tea.Quit, false, fmt.Errorf("Main Menu Start Server: %s", err)
				}
				log.Printf("(Dwarf Wars Server) Validating SQL Servers and Databases")
				err = ValidateSQLServers(opts.SQLServers.ToList())
				if err != nil {
					log.Printf("(Dwarf Wars Server) Failed to Validate SQL Server & Databases: %s", err)
					return nil, false, err
				}
				log.Printf("(Dwarf Wars Server) Starting Server")
				go func() {
					waitgroup.Add(1)
					defer waitgroup.Done()
					serv.StartTCP(ctxServer)
				}()
				s = true
			} else {
				serverCancel()
				serv.Wait()
				log.Printf("(TUI) Stopped Server")
				s = false
			}
			return
		}).
		Add("Start Game Server", false, func(state bool) (cmd tea.Cmd, s bool, err error) {
			if !state {
				log.Printf("(Dwarf Wars Server) Validating SQL Servers and Databases")
				err = ValidateSQLServers(opts.SQLServers.ToList())
				if err != nil {
					log.Printf("(Dwarf Wars Server) Failed to Validate SQL Server & Databases: %s", err)
					return nil, false, err
				}
			} else {
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
			ctx_cancel()
			return tea.Quit, state, nil
		})
	op := tea.NewProgram(mainMenu)
	if _, err := op.Run(); err != nil {
		log.Fatalf("(Bubbletea) Error, %s", err)
	}
}
