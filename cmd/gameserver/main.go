package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/maplelm/dwarfwars/pkg/server"
	"github.com/maplelm/dwarfwars/pkg/tui"
	"os"
	_ "time"
)

func main() {
	fmt.Println("Starting Server...")
	serv, err := server.New("0.0.0.0", "3000")
	if err != nil {
		fmt.Printf("Failed to create Server, %s\n", err)
		os.Exit(1)
	}
	m := tui.OptionScreenInit()
	m.Add("Shutdown", &serverShutdown{false})
	op := tea.NewProgram(m)
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
