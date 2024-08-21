package main

import (
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
		serv      *server.Server
	)

	// getting settings //
	/*
		checking for specified settings path/name in envirment variable
			-> SETTINGS_PATH
			-> SETTINGS_NAME
	*/
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
		log.Fatalf("Main Thread: Failed to load main TOML settings file, %s", err)
	}

	opts, err := settings.Get[settings.Config]("Main")
	if err != nil {
		log.Fatalf("Main Thread: Failed to get Main settings from memory, %s", err)
	}

	//////////////////////////////
	// Setting up System Logger //
	//////////////////////////////
	log.SetOutput(NewRotationWriter(
		opts.Log.Path,
		opts.Log.FileName,
		opts.Log.AdjustedPollRate(),
		opts.Log.MaxFileSize),
	)
	log.SetPrefix("System: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//////////////////////////
	// Setup Main TUI Model //
	//////////////////////////

	/////////////////////////
	// Starting TCP Server //
	/////////////////////////
	mainMenu := tui.NewMenu('>', "Dwarf Wars Server", lg.NewStyle(), lg.NewStyle(), lg.NewStyle())
	mainMenu = mainMenu.
		Add("Start Server", false, func(state bool) (cmd tea.Cmd, s bool, err error) {
			serv, err = server.New(opts.Server.Addr, fmt.Sprintf("%d", opts.Server.Port))
			if err != nil {
				return tea.Quit, false, fmt.Errorf("Main Menu Start Server: %s", err)
			}
			go func() {
				waitgroup.Add(1)
				defer waitgroup.Done()
				serv.Start()
			}()
			s = true
			return
		}).
		Add("Settings", false, func(state bool) (cmd tea.Cmd, s bool, err error) {
			var inter interface{} = *opts
			settingsMenu := tui.NewSettingsMenu(&inter, '>', lg.NewStyle(), lg.NewStyle(), lg.NewStyle())
			p := tea.NewProgram(settingsMenu)
			if _, err = p.Run(); err != nil {
				return nil, state, fmt.Errorf("Settings Menu: %s", err)
			}
			s = state
			return
		}).
		Add("Quit", false, func(state bool) (cmd tea.Cmd, s bool, err error) {
			return tea.Quit, state, nil
		})
	op := tea.NewProgram(mainMenu)
	_, err = op.Run()
	if err != nil {
		fmt.Printf("Bubbletea Error: %s\n", err)
		os.Exit(1)
	}
}

type ActionServerShutdown struct {
	state bool
}

func (ashut *ActionServerShutdown) Enable(m tea.Model) (tea.Model, tea.Cmd) {
	fmt.Println("Shutdown Enabled")
	ashut.state = true
	return m, tea.Quit
}

func (ashut *ActionServerShutdown) Disable(m tea.Model) (tea.Model, tea.Cmd) {
	fmt.Println("Shutdown Disabled")
	ashut.state = false
	return m, tea.Quit
}

func (ashut *ActionServerShutdown) State() bool {
	return bool(ashut.state)
}

type ActionStartServer struct {
	serv server.Server
}

func (ass *ActionStartServer) Enable(m tea.Model) (tea.Model, tea.Cmd) {
	return m, nil
}

func (ass *ActionStartServer) Disable(m tea.Model) (tea.Model, tea.Cmd) {
	return m, nil
}

func (ass *ActionStartServer) State() bool {
	return true
}

type ActionSettings struct{}

func (as *ActionSettings) Enable(m tea.Model) (tea.Model, tea.Cmd)  { return m, nil }
func (as *ActionSettings) Disable(m tea.Model) (tea.Model, tea.Cmd) { return m, nil }
func (as *ActionSettings) State() bool                              { return true }
