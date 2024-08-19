package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"log"
)

type Action func() (tea.Cmd, error)

func (f Action) action() (tea.Cmd, error) {
	return f()
}

type ActionHandler interface {
	action() (tea.Cmd, error)
}

func CreateAction(f func() (tea.Cmd, error)) ActionHandler {
	return Action(f)
}

type Menu struct {
	cursor        int
	CursorIcon    rune
	labels        []string
	actions       []ActionHandler
	Title         string
	MainStyle     lg.Style
	SelectedStyle lg.Style
	ItemStyle     lg.Style
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch ms := msg.(type) {
	case tea.KeyMsg:
		switch ms.String() {
		case "j", "down":
			if m.cursor < len(m.actions) {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "l", "right", "enter", " ":
			c, e := m.actions[m.cursor].action()
			if e != nil {
				log.Printf("Failed Action for %s: %s", m.labels[m.cursor], e)
				return m, nil
			}
			return m, c
		case "q", "ctrl+q":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m Menu) View() string {
	output := fmt.Sprintf(" %s \n", m.Title)
	for i, l := range m.labels {
		if m.cursor == i {
			output += fmt.Sprintf(" %c %s\n", m.CursorIcon, m.SelectedStyle.Render(l))
		} else {
			output += fmt.Sprintf("   %s\n", m.ItemStyle.Render(l))
		}
	}
	return m.MainStyle.Render(output)
}

func (m Menu) Add(l string, a func() (tea.Cmd, error)) Menu {
	for _, v := range m.labels {
		if v == l {
			fmt.Printf("%s already exists in menu\n", l)
			return m
		}
	}
	m.labels = append(m.labels, l)
	m.actions = append(m.actions, Action(a))
	return m
}

func (m Menu) Remove(l string) Menu {
	for i, v := range m.labels {
		if v == l {
			m.labels = append(m.labels[:i], m.labels[i:]...)
			return m
		}
	}
	fmt.Printf("%s does not exist in menu\n", l)
	return m
}
