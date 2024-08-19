package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"reflect"
)

type Option interface {
	Enable() error
	Disable() error
	State() bool
}

type OptionScreen struct {
	cursor  int
	Options map[string]Option
}

func OptionScreenInit() OptionScreen {
	return OptionScreen{
		cursor:  0,
		Options: make(map[string]Option),
	}
}

func (os OptionScreen) Init() tea.Cmd {
	// Do nothing on startup
	return nil
}

func (os OptionScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "j", "down":
			if os.cursor < len(os.Options)-1 {
				os.cursor++
			}
		case "k", "up":
			if os.cursor > 0 {
				os.cursor--
			}
		//case "h", "left":
		case "l", "right", "enter", " ":
			index := 0
			for _, v := range os.Options {
				if index == os.cursor {
					if v.State() {
						v.Disable()
					} else {
						v.Enable()
						return os, tea.Quit
					}
					break
				}
				index++
			}
		case "ctrl+c", "q":
		}
	default:
		log.Printf("Unhandled msg: %s", reflect.TypeOf(msg).Kind())
	}
	return os, nil
}

func (os OptionScreen) View() string {
	s := fmt.Sprintf("Select Option (%d)\n", len(os.Options))
	cpos := 0
	c := ""
	for k := range os.Options {
		if cpos == os.cursor {
			c = ">"
		} else {
			c = " "
		}
		s += fmt.Sprintf("# %s %s\n", c, k)
	}
	return s
}

func (os OptionScreen) Add(key string, opt Option) OptionScreen {
	if _, ok := os.Options[key]; !ok {
		os.Options[key] = opt
	}
	return os
}

func (os *OptionScreen) Remove(key string) error {
	if _, ok := os.Options[key]; !ok {
		return fmt.Errorf("%s does not exist as an option", key)
	}
	delete(os.Options, key)
	return nil
}
