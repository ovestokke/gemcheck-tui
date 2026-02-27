package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Tab1    key.Binding
	Tab2    key.Binding
	Tab3    key.Binding
	NextTab key.Binding
	Search  key.Binding
	Refresh key.Binding
	Select  key.Binding
	Back    key.Binding
	Up      key.Binding
	Down    key.Binding
	Quit    key.Binding
}

var Keys = KeyMap{
	Tab1: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "red"),
	),
	Tab2: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "green"),
	),
	Tab3: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "blue"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
