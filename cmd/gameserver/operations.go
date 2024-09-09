package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/maplelm/dwarfwars/pkg/logging"
	"github.com/maplelm/dwarfwars/pkg/server"
	"github.com/maplelm/dwarfwars/pkg/settings"
	"github.com/maplelm/dwarfwars/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Server struct {
	Server    *server.Server
	Ctx       context.Context
	CtxCancel func()
	waitgroup *sync.WaitGroup
}

func (sc *Server) IsSelect(m tui.Menu, i int) bool {
	return m.IsEnabled(i)
}

func (sc *Server) Toggle(m tui.Menu, i int) (tea.Cmd, tui.Menu, error) {
	logging.SetTitle("Dwarf Wars Game Server")
	defer logging.SetTitle("System")

	opts, err := settings.Get[settings.Config]("Main")
	if err != nil {
		return nil, m, err
	}

	logging.Info("Running the ServerSelecter.Command")
	logging.Infof(" IsEnabled: %t", m.IsEnabled(i))
	if !m.IsEnabled(i) {
		sc.Server, err = server.New(opts.Server.Addr, fmt.Sprintf("%d", opts.Server.Port), time.Duration(opts.Server.IdleTimeout)*time.Millisecond)
		if err != nil {
			logging.Error(err, "Failed to create server")
			return tea.Quit, m, fmt.Errorf("Main Menu Start Server: %s", err)
		}
		logging.Info("Validating SQL Servers & Databases")
		err = ValidateSQLServers(opts.Databases)
		if err != nil {
			logging.Error(err, "Failed to Validate SQL Server & Databases")
			return nil, m, err
		}
		logging.Info("Validating SQL Servers & Databases Successfull")

		if sc.CtxCancel == nil {
			if sc.Ctx == nil {
				sc.Ctx = context.Background()
			}
			sc.Ctx, sc.CtxCancel = context.WithCancel(sc.Ctx)
		}

		logging.Infof("waiting group status: %v", sc.waitgroup)
		logging.Info("Starting Server")
		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			sc.Server.StartTCP(sc.Ctx, make(chan struct{})) // Passing channel this way because I am using the fucntion syncronously in this instance
		}(sc.waitgroup)
		m.SetEnabled(true, i)
	} else {
		sc.CtxCancel()
		sc.Server.Wait()
		logging.Info("Server stopped")
		m.SetEnabled(false, i)
	}

	return nil, m, nil
}

func EditSettings(m tui.Menu, i int) (tea.Cmd, tui.Menu, error) {
	logging.SetTitle("TUI")
	defer logging.SetTitle("System")
	opts, err := settings.Get[settings.Config]("Main")
	if err != nil {
		return nil, m, err
	}
	var inter interface{} = opts
	settingsMenu := tui.NewSettingsMenu(&inter, '>', lg.NewStyle(), lg.NewStyle(), lg.NewStyle())
	p := tea.NewProgram(settingsMenu)
	logging.Info("Opening Settings Window")
	if _, err = p.Run(); err != nil {
		logging.Error(err, "Failed to open settings window")
		return nil, m, fmt.Errorf("Settings Menu: %s", err)
	}
	return nil, m, nil
}

func Quit(m tui.Menu, i int) (tea.Cmd, tui.Menu, error) {
	logging.SetTitle("TUI")
	defer logging.SetTitle("System")
	logging.Info("Quit Command Selected, Manual Shutdown Started")
	if m.Ctxcancel != nil {
		m.Ctxcancel()
	}
	return tea.Quit, m, nil
}
