package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

type model struct {
	domains []string // holds machine host file contents (implement)
	state   sessionState
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update() (tea.Model, tea.Cmd) {
	// q => quits the program
	// arrows => pick options
	// enter or space => select options
	// add mode if needed ...
}

func (m model) View() string {
	// implement list-fancy view from bbt
	// https://github.com/charmbracelet/bubbletea/blob/master/examples/list-fancy/README.md
}
