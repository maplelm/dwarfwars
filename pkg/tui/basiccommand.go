package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BasicMenuCommand func(Menu, int) (tea.Cmd, Menu, error)

func (f BasicMenuCommand) Command(m Menu, i int) (tea.Cmd, Menu, error) {
	return f(m, i)
}

func (f BasicMenuCommand) Enabled(m Menu, i int) bool {
	return m.enabled[i]
}
