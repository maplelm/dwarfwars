package tui

import (
	"context"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"

	"github.com/maplelm/dwarfwars/pkg/logging"
)

type Selecter interface {
	Select(Menu, int) (tea.Cmd, Menu, error)
	IsSelect(Menu, int) bool
}

type menuOption struct {
	Label    string
	enabled  bool
	selecter Selecter
}

type Menu struct {
	cursor  int
	options []menuOption

	Ctxcancel     func()
	CursorIcon    rune
	Ctx           context.Context
	Title         string
	MainStyle     lg.Style
	SelectedStyle lg.Style
	ItemStyle     lg.Style
}

func NewMenu(icon rune, title string, ms, ss, is lg.Style, ctx context.Context) *Menu {
	return &Menu{
		CursorIcon:    icon,
		Title:         title,
		MainStyle:     ms,
		SelectedStyle: ss,
		ItemStyle:     is,
		Ctx: func(c context.Context) context.Context {
			if c == nil {
				return context.Background()
			}
			return ctx
		}(ctx),
	}
}

func (m Menu) IsEnabled(i int) bool {
	return m.options[i].enabled
}

func (m *Menu) SetEnabled(s bool, i int) Menu {
	m.options[i].enabled = s
	return *m
}

func (m *Menu) Increase(i int) {
	if m.cursor+i < len(m.options) {
		m.cursor += i
	}
}

func (m *Menu) Decrease(i int) {
	if m.cursor-i >= 0 {
		m.cursor -= i
	}
}

func (m *Menu) RunCommand() (tea.Cmd, Menu, error) {
	return m.options[m.cursor].selecter.Select(*m, m.cursor)
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch ms := msg.(type) {
	case tea.KeyMsg:
		switch ms.String() {
		case "j", "down":
			m.Increase(1)
		case "k", "up":
			m.Decrease(1)
		case "l", "right", "enter", " ":
			logging.Infof("Running Command (%d): %s", m.cursor, m.options[m.cursor].Label)
			logging.Infof("length of cmds: %d", len(m.options))
			logging.Infof("Command Value: %p", m.options[m.cursor].selecter.Select)
			var (
				c   tea.Cmd
				err error
			)
			if c, m, err = m.RunCommand(); err != nil {
				logging.Errorf(err, "Failed Action: %s", m.options[m.cursor].Label)
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
	for i, opt := range m.options {
		if m.cursor == i {
			line = fmt.Sprintf(" %c %s\n", m.CursorIcon, m.SelectedStyle.Render(opt.Label))
		} else {
			line = fmt.Sprintf("   %s\n", m.ItemStyle.Render(opt.Label))
		}
		if m.options[i].enabled {
			output += fmt.Sprintf("%s\n", lg.NewStyle().Foreground(lg.Color("#00FF00")).Render(line))
		} else {
			output += line
		}
	}
	return m.MainStyle.Render(output)
}

func (m Menu) Add(label string, initState bool, cmd Selecter) Menu {
	for _, v := range m.options {
		if v.Label == label {
			log.Printf("%s already exists in menu", label)
			return m
		}
	}
	m.options = append(m.options, menuOption{
		Label:    label,
		enabled:  initState,
		selecter: cmd,
	})
	return m
}

func (m Menu) AddFunc(l string, state bool, a func(Menu, int) (tea.Cmd, Menu, error)) Menu {
	return m.Add(l, state, BasicMenuSelecter(a))
}

func (m Menu) Remove(l string) Menu {
	for i, v := range m.options {
		if v.Label == l {
			m.options = append(m.options[:i], m.options[:i]...)
			return m
		}
	}
	fmt.Printf("%s does not exist in menu\n", l)
	return m
}
