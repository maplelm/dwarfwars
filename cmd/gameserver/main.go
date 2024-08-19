package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbletea"
	"github.com/maplelm/dwarfwars/pkg/server"
	"github.com/maplelm/dwarfwars/pkg/tui"
	"log"
	"os"
	"path/filepath"
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
	m := tui.OptionScreenInit()
	m.Add("Shutdown", &serverShutdown{false})
	op := tea.NewProgram(tui.OptionScreenInit().Add("Shutdown", &serverShutdown{false}).Add("Other Shutdown", &serverShutdown{false}))
	_, err = op.Run()
	if err != nil {
		fmt.Printf("Bubbletea Error: %s\n", err)
		os.Exit(1)
	}
	go serv.Start()
	serv.Stop()
	fmt.Println("Closing Server...")
}

type serverShutdown struct {
	state bool
}

func (ss *serverShutdown) Enable() error {
	fmt.Println("Shutdown Enabled")
	ss.state = true
	return nil
}

func (ss *serverShutdown) Disable() error {
	fmt.Println("Shutdown Disabled")
	ss.state = false
	return nil
}

func (ss *serverShutdown) State() bool {
	return bool(ss.state)
}
