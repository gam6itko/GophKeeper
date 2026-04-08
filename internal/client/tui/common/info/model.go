package info

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	kvp map[string]string
	// quitCmd команда которую нужно отправить когда пользователь нажал 'q'.
	quitCmd tea.Cmd
}

func NewModel(kvp map[string]string, exitCmd tea.Cmd) *Model {
	return &Model{
		kvp:     kvp,
		quitCmd: exitCmd,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "q" {
			return m, m.quitCmd
		}
	}

	return m, nil
}

func (m Model) View() string {
	result := "\n"
	for k, v := range m.kvp {
		result += fmt.Sprintf("%s: %s\n", k, v)
	}

	result += "\n\nPress q to quit."

	return result
}
