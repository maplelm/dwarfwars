package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BasicMenuSelecter func(Menu, int) (tea.Cmd, Menu, error)

func (f BasicMenuSelecter) Select(m Menu, i int) (tea.Cmd, Menu, error) {
	return f(m, i)
}

func (f BasicMenuSelecter) IsSelect(m Menu, i int) bool {
	return m.options[i].enabled
}
