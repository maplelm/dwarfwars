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

var (
	mainSettingsPath string
)

func main() {
	var (
		mainWaitGroup sync.WaitGroup
		gameServer    *server.Server   = nil
		opts          *settings.Config = nil
	)

	opts, err := LoadSettings("./", "settings.toml")
	if err != nil {
		log.Fatalf("Failed to load Main Settings TOML File, %s", err)
	}

	err = InitLogger(opts)
	if err != nil {
		log.Fatalf("Failed to initiate Logger, %s", err)
	}

	mainCtx, ctxCancel := context.WithCancel(context.Background())

	mainMenu := InitTui(mainCtx, &mainWaitGroup, opts, gameServer)

	//////////////////////////////
	// Setting Up TUI / Actions //
	//////////////////////////////
	op := tea.NewProgram(mainMenu)
	if _, err := op.Run(); err != nil {
		log.Fatalf("(Bubbletea) Error, %s", err)
	}
	ctxCancel()
	log.Printf("Losing Game Server, %s", err)
	mainWaitGroup.Wait()
}

func LoadSettings(defaultpath, defaultname string) (opts *settings.Config, err error) {
	settingsPath := os.Getenv("SETTINGS_PATH")
	if len(settingsPath) == 0 {
		settingsPath = defaultpath
	}
	settingsName := os.Getenv("SETTINGS_NAME")
	if len(settingsName) == 0 {
		settingsName = defaultname
	}
	_, err = settings.LoadFromTomlFile("Main", settingsPath, settingsName)
	if err != nil {
		return
	}
	opts, err = settings.Get[settings.Config]("Main")
	return
}

func InitLogger(opts *settings.Config) (err error) {
	//////////////////////////////
	// Setting up System Logger //
	//////////////////////////////
	fmt.Printf("Creating logger with path: %s and file name: %s\n", opts.Log.Path, opts.Log.FileName)
	log.SetOutput(NewRotationWriter(
		opts.Log.Path,
		opts.Log.FileName,
		opts.Log.MaxFileSize),
	)
	log.SetPrefix("Game Server:")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Validating that log path exists
	_, err = os.Stat(opts.Log.Path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(opts.Log.Path, 0777)
	}
	return
}

func InitTui(parentCtx context.Context, wg *sync.WaitGroup, opts *settings.Config, serv *server.Server) (mm *tui.Menu) {
	mainMenu := tui.NewMenu('>', "Dwarf Wars Server", lg.NewStyle(), lg.NewStyle(), lg.NewStyle())
	ctxServer, serverCancel := context.WithCancel(parentCtx)
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
					wg.Add(1)
					defer wg.Done()
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
			serverCancel()
			return tea.Quit, state, nil
		})
	mm = &mainMenu
	return
}
