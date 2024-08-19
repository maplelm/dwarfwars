package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	//"log"
	//"reflect"
)

type OptionSetter interface {
	Enable(tea.Model) (tea.Model, tea.Cmd)
	Disable(tea.Model) (tea.Model, tea.Cmd)
	State() bool
}

type OptionMenu struct {
	cursor        int
	labels        []string
	options       map[int]OptionSetter
	MainStyle     lg.Style
	ItemStyle     lg.Style
	SelectedStyle lg.Style
	Title         string
}

func OptionScreenInit(title string) OptionMenu {
	return OptionMenu{
		cursor:  0,
		labels:  []string{},
		options: make(map[int]OptionSetter),
		Title:   title,
	}
}

func (os OptionMenu) Init() tea.Cmd {
	// Do nothing on startup
	return nil
}

func (os OptionMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "j", "down":
			if os.cursor < len(os.options)-1 {
				os.cursor++
			}
		case "k", "up":
			if os.cursor > 0 {
				os.cursor--
			}
		//case "h", "left":
		case "l", "right", "enter", " ":
			if os.options[os.cursor].State() {
				return os.options[os.cursor].Disable(os)
			}
			return os.options[os.cursor].Enable(os)
		case "ctrl+c", "q":
			return os, tea.Quit
		}
	default:
		//log.Printf("Unhandled msg: %s", reflect.TypeOf(msg).Kind())
	}
	return os, nil
}

func (os OptionMenu) View() string {
	output := fmt.Sprintf("%s\n", os.Title)
	for i, k := range os.labels {
		if i == os.cursor {
			output += fmt.Sprintf("%s\n", os.SelectedStyle.Render(fmt.Sprintf(" > %s", k)))
		} else {
			output += fmt.Sprintf("%s\n", os.ItemStyle.Render(fmt.Sprintf("   %s", k)))
		}
	}
	return fmt.Sprintf("%s\n", os.MainStyle.Render(output))
}

func (os OptionMenu) Add(key string, opt OptionSetter) OptionMenu {
	os.options[len(os.labels)] = opt
	os.labels = append(os.labels, key)
	return os
}

func (os *OptionMenu) Remove(key string) error {
	for k, v := range os.labels {
		if key == v {
			delete(os.options, k)
			os.labels = append(os.labels[:k], os.labels[k:]...)
			return nil
		}
	}
	return fmt.Errorf("%s does not exist as an option", key)
}
