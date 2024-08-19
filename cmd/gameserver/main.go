package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/maplelm/dwarfwars/pkg/server"
	"github.com/maplelm/dwarfwars/pkg/tui"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {

	var (
		opts Config
		//sysLogger *log.Logger
	)

	/////////////////////////////////////
	// getting settings from toml file //
	/////////////////////////////////////
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

	b, err := os.ReadFile(filepath.Join(settingsPath, settingsName))

	err = toml.Unmarshal(b, &opts)
	if err != nil {
		log.Fatalf("Failed to Unmarshal %s\n", filepath.Join(settingsPath, settingsName))
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

	/////////////////////////
	// Starting TCP Server //
	/////////////////////////
	fmt.Println("Starting Server...")
	serv, err := server.New(opts.Server.Addr, fmt.Sprintf("%d", opts.Server.Port))
	if err != nil {
		fmt.Printf("Failed to create Server, %s\n", err)
		os.Exit(1)
	}
	m := tui.OptionScreenInit("Dwarf Wars Server:").
		Add("Start Server", &ActionStartServer{}).
		Add("Settings", &ActionSettings{}).
		Add("Quit", &ActionServerShutdown{false})
	m.MainStyle = lg.NewStyle().
		Background(lg.Color("#FFFFFF")).
		Foreground(lg.Color("#000000")).
		Border(lg.RoundedBorder(), true).
		Bold(false).
		Width(40).
		Height(20).
		Padding(1)
	m.SelectedStyle = (lg.NewStyle().
		Background(lg.Color("#000000")).
		Foreground(lg.Color("#FFFFFF")).
		Border(lg.DoubleBorder(), false, false, true, false).
		PaddingRight(2))
	m.ItemStyle = lg.NewStyle().
		Background(lg.Color("#FFFFFF")).
		Foreground(lg.Color("#000000")).
		Border(lg.NormalBorder(), false).
		PaddingRight(2)

	men := tui.Menu{
		CursorIcon:    '>',
		Title:         "Main Menu",
		MainStyle:     lg.NewStyle(),
		SelectedStyle: lg.NewStyle(),
		ItemStyle:     lg.NewStyle(),
	}.
		Add("Test Menu", func() (tea.Cmd, error) {
			p := tea.NewProgram(m)
			p.Run()
			return nil, nil
		}).
		Add("Test Action", func() (tea.Cmd, error) {
			log.Printf("Test Action: %s\n", time.Now().Format(time.TimeOnly))
			return nil, nil
		})
	op := tea.NewProgram(men)
	_, err = op.Run()
	if err != nil {
		fmt.Printf("Bubbletea Error: %s\n", err)
		os.Exit(1)
	}
	go serv.Start()
	serv.Stop()
	fmt.Println("Closing Server...")
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
