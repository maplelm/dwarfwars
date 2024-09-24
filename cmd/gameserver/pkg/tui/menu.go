package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Menu struct {
	//public
	CursorIcon rune
	List       list.Model
	Style      lg.Style
	// Private
	cursor int
}

func InitMenu(icon rune, options []Option, title string) Menu {
	items := make([]list.Item, len(options))
	for i, e := range options {
		items[i] = e
	}
	m := Menu{
		CursorIcon: icon,
		List:       list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.List.Title = title
	return m
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := m.Style.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Menu) View() string {
	return m.Style.Render(m.List.View())
}
