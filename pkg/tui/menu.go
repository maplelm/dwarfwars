package tui

import (
	"context"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type CommandExecuter interface {
	Command(Menu, int) (tea.Cmd, Menu, error)
	Enabled(Menu, int) bool
}
type Menu struct {
	cursor  int
	labels  []string
	enabled map[int]bool
	cmds    map[int]CommandExecuter

	Ctxcancel     func()
	CursorIcon    rune
	Ctx           context.Context
	Title         string
	MainStyle     lg.Style
	SelectedStyle lg.Style
	ItemStyle     lg.Style
}

func NewMenu(icon rune, title string, ms, ss, is lg.Style, ctx context.Context) Menu {
	if ctx == nil {
		ctx = context.Background()
	}
	return Menu{
		CursorIcon:    icon,
		Title:         title,
		MainStyle:     ms,
		SelectedStyle: ss,
		ItemStyle:     is,
		Ctx:           ctx,
		enabled:       make(map[int]bool),
		cmds:          make(map[int]CommandExecuter),
	}
}

func (m Menu) IsEnabled(i int) bool {
	return m.enabled[i]
}

func (m *Menu) SetEnabled(s bool, i int) Menu {
	m.enabled[i] = s
	return *m
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch ms := msg.(type) {
	case tea.KeyMsg:
		switch ms.String() {
		case "j", "down":
			if m.cursor < len(m.cmds) {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "l", "right", "enter", " ":
			log.Printf(`{"Type": "Info", "Source": "TUI", "Msg": "Running Command (%d) %s)"}`, m.cursor, m.labels[m.cursor])
			log.Printf(`{"Type": "Info", "Source": "TUI", "Msg": "%s"}`, fmt.Sprintf("lenth of cmds: %d", len(m.cmds)))
			log.Printf(`{"Type": "Info", "Source": "TUI", "Msg": "%s"}`, fmt.Sprintf("Command Value: %p", m.cmds[m.cursor].Command))
			c, m, e := m.cmds[m.cursor].Command(m, m.cursor)
			if e != nil {
				log.Printf("Failed Action for %s: %s", m.labels[m.cursor], e)
			}
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
			output += fmt.Sprintf("%s\n", lg.NewStyle().Foreground(lg.Color("#00FF00")).Render(line))
		} else {
			output += line
		}
	}
	return m.MainStyle.Render(output)
}

func (m Menu) Add(label string, initState bool, cmd CommandExecuter) Menu {
	for _, v := range m.labels {
		if v == label {
			log.Printf("%s already exists in menu", label)
			return m
		}
	}
	m.labels = append(m.labels, label)
	m.enabled[len(m.labels)-1] = initState
	m.cmds[len(m.labels)-1] = cmd
	return m
}

func (m Menu) AddFunc(l string, state bool, a func(Menu, int) (tea.Cmd, Menu, error)) Menu {
	return m.Add(l, state, BasicMenuCommand(a))
}

func (m Menu) Remove(l string) Menu {
	for i, v := range m.labels {
		if v == l {
			m.labels = append(m.labels[:i], m.labels[i:]...)
			delete(m.enabled, i)
			delete(m.cmds, i)
			return m
		}
	}
	fmt.Printf("%s does not exist in menu\n", l)
	return m
}
