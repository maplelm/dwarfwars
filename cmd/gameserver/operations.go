package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/maplelm/dwarfwars/pkg/server"
	"github.com/maplelm/dwarfwars/pkg/settings"
	"github.com/maplelm/dwarfwars/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type ServerCommand struct {
	Server    *server.Server
	Ctx       context.Context
	CtxCancel func()
	waitgroup *sync.WaitGroup
}

func (sc *ServerCommand) Enabled(m tui.Menu, i int) bool {
	return m.IsEnabled(i)
}

func (sc *ServerCommand) Command(m tui.Menu, i int) (tea.Cmd, tui.Menu, error) {
	SetLogEntryTitle("Dwarf Wars Game Server")
	defer SetLogEntryTitle("System")

	opts, err := settings.Get[settings.Config]("Main")
	if err != nil {
		return nil, m, err
	}

	LogInfo("Running the ServerCommand.Command")
	LogInfof(" IsEnabled: %t", m.IsEnabled(i))
	if !m.IsEnabled(i) {
		sc.Server, err = server.New(opts.Server.Addr, fmt.Sprintf("%d", opts.Server.Port), time.Duration(opts.Server.IdleTimeout)*time.Millisecond)
		if err != nil {
			LogError(err, "Failed to create server")
			return tea.Quit, m, fmt.Errorf("Main Menu Start Server: %s", err)
		}
		LogInfo("Validating SQL Servers & Databases")
		err = ValidateSQLServers(opts.SQLServers.ToList())
		if err != nil {
			LogError(err, "Failed to Validate SQL Server & Databases")
			return nil, m, err
		}
		LogInfo("Validating SQL Servers & Databases Successfull")

		if sc.CtxCancel == nil {
			if sc.Ctx == nil {
				sc.Ctx = context.Background()
			}
			sc.Ctx, sc.CtxCancel = context.WithCancel(sc.Ctx)
		}

		LogInfof("waiting group status: %v", sc.waitgroup)
		LogInfo("Starting Server")
		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			sc.Server.StartTCP(sc.Ctx)
		}(sc.waitgroup)
		m.SetEnabled(true, i)
	} else {
		sc.CtxCancel()
		sc.Server.Wait()
		LogInfo("Server stopped")
		m.SetEnabled(false, i)
	}

	return nil, m, nil
}

func EditSettings(m tui.Menu, i int) (tea.Cmd, tui.Menu, error) {
	SetLogEntryTitle("TUI")
	defer SetLogEntryTitle("System")
	opts, err := settings.Get[settings.Config]("Main")
	if err != nil {
		return nil, m, err
	}
	var inter interface{} = *opts
	settingsMenu := tui.NewSettingsMenu(&inter, '>', lg.NewStyle(), lg.NewStyle(), lg.NewStyle())
	p := tea.NewProgram(settingsMenu)
	LogInfo("Opening Settings Window")
	if _, err = p.Run(); err != nil {
		LogError(err, "Failed to open settings window")
		return nil, m, fmt.Errorf("Settings Menu: %s", err)
	}
	return nil, m, nil
}

func Quit(m tui.Menu, i int) (tea.Cmd, tui.Menu, error) {
	SetLogEntryTitle("TUI")
	defer SetLogEntryTitle("System")
	LogInfo("Quit Command Selected, Manual Shutdown Started")
	if m.Ctxcancel != nil {
		m.Ctxcancel()
	}
	return tea.Quit, m, nil
}
