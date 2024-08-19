package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"log"
	"reflect"
)

type SettingHandler interface {
}

type SettingsMenu struct {
	CursorIcon  rune
	MainStyle   lg.Style
	SelectStyle lg.Style
	ItemStyle   lg.Style
	settings    *any
	types       []reflect.Type
	labels      []string
	cursor      int
	length      int
}

func NewSettingsMenu(settings *any, icon rune, ms, ss, is lg.Style) SettingsMenu {
	var (
		v      reflect.Value
		l      int
		labels []string       = []string{}
		types  []reflect.Type = []reflect.Type{}
	)
	v = reflect.ValueOf(settings)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		l = v.Len()
		for i := range l {
			labels = append(labels, fmt.Sprintf("option %d", i))
			types = append(types, v.Type())
		}
	case reflect.Struct:
		l = v.NumField()
		for i := 0; i < l; i++ {
			labels = append(labels, v.Type().Field(i).Name)
			types = append(types, v.Type().Field(i).Type)
		}
	default:
		log.Printf("failed to find the number of fields in passed settings data, Setting to 1")
		l = 1
	}
	return SettingsMenu{
		settings:    settings,
		CursorIcon:  icon,
		MainStyle:   ms,
		SelectStyle: ss,
		ItemStyle:   is,
		length:      l,
		labels:      labels,
		types:       types,
	}
}

func (sm SettingsMenu) Init() tea.Cmd {
	return nil
}

func (sm SettingsMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "j", "down":
			if sm.cursor < sm.length-1 {
				sm.cursor++
			}
		case "k", "up":
			if sm.cursor > 0 {
				sm.cursor--
			}
		case "l", "right", "enter", " ":
			switch sm.types[sm.cursor].Kind() {
			case reflect.String:
			case reflect.Map:
			case reflect.Array:
			case reflect.Float32, reflect.Float64:
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			case reflect.Slice:
			case reflect.Struct:
			default: // unsupported
			}

		case "h", "q", "ctrl+c":
			return sm, tea.Quit
		}
	}
	return sm, nil
}

func (sm SettingsMenu) View() string {
	return ""
}
