package loading

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	DoneCmd struct {
		Cmd tea.Cmd
	}
)

type Model struct {
	spinner spinner.Model
	initCmd tea.Cmd

	message string
}

func New(initCmd tea.Cmd) Model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		spinner: s,
		initCmd: initCmd,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.initCmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DoneCmd:
		return m, msg.Cmd
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf("\n\n   %s Waiting for response...\n\n", m.spinner.View())
}
