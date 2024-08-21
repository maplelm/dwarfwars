package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"log"
)

// Takes in state and returns state and bubble team command (can fail)
type Action func(bool) (cmd tea.Cmd, state bool, err error)

func (f Action) action(b bool) (cmd tea.Cmd, state bool, err error) {
	return f(b)
}

type ActionHandler interface {
	action(bool) (cmd tea.Cmd, state bool, err error)
}

func CreateAction(f func(bool) (cmd tea.Cmd, state bool, err error)) ActionHandler {
	return Action(f)
}

type Menu struct {
	cursor        int
	CursorIcon    rune
	labels        []string
	enabled       []bool
	actions       []ActionHandler
	Title         string
	MainStyle     lg.Style
	SelectedStyle lg.Style
	ItemStyle     lg.Style
}

func NewMenu(icon rune, title string, ms, ss, is lg.Style) Menu {
	return Menu{
		CursorIcon:    icon,
		Title:         title,
		MainStyle:     ms,
		SelectedStyle: ss,
		ItemStyle:     is,
	}
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
			c, s, e := m.actions[m.cursor].action(m.enabled[m.cursor])
			if e != nil {
				log.Printf("Failed Action for %s: %s", m.labels[m.cursor], e)
				return m, nil
			}
			m.enabled[m.cursor] = s // update state
			return m, c
		case "q", "ctrl+q":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m Menu) View() string {
	var line string
	output := fmt.Sprintf(" %s \n", m.Title)
	for i, l := range m.labels {
		if m.cursor == i {
			line = fmt.Sprintf(" %c %s\n", m.CursorIcon, m.SelectedStyle.Render(l))
		} else {
			line = fmt.Sprintf("   %s\n", m.ItemStyle.Render(l))
		}
		if m.enabled[i] {
			output += lg.NewStyle().Foreground(lg.Color("#00FF00")).Render(line)
		} else {
			output += line
		}
	}
	return m.MainStyle.Render(output)
}

func (m Menu) Add(l string, state bool, a func(bool) (tea.Cmd, bool, error)) Menu {
	for _, v := range m.labels {
		if v == l {
			fmt.Printf("%s already exists in menu\n", l)
			return m
		}
	}
	m.labels = append(m.labels, l)
	m.enabled = append(m.enabled, state)
	m.actions = append(m.actions, Action(a))
	return m
}

func (m Menu) Remove(l string) Menu {
	for i, v := range m.labels {
		if v == l {
			m.labels = append(m.labels[:i], m.labels[i:]...)
			m.enabled = append(m.enabled[:i], m.enabled[i:]...)
			m.actions = append(m.actions[:i], m.actions[i:]...)
			return m
		}
	}
	fmt.Printf("%s does not exist in menu\n", l)
	return m
}
