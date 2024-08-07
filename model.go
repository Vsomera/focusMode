package main

type sessionState uint

type model struct {
	domains []string // holds machine host file contents (implement)
	state   sessionState
}

func newModel() *model {
	return &model{}
}

func (m model) Init() {}

func (m model) Update() {
	// q => quits the program
	// arrows => pick options
	// enter or space => select options
	// add mode if needed ...
}

func (m model) View() {
	// implement list-fancy view from bbt
	// https://github.com/charmbracelet/bubbletea/blob/master/examples/list-fancy/README.md
}
