package loading

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	SuccessMsg struct{}
	FailMsg    struct {
		GoToModel tea.Model
	}
	CloseMsg struct{}

	ResponseResultMsg struct {
		Success bool
		Message string
	}
)

type tState int

const (
	stateLoading tState = iota
	stateSuccess
	stateFail
)

type Model struct {
	spinner      spinner.Model
	initCmd      tea.Cmd
	onSuccessCmd tea.Cmd
	onFailCmd    tea.Cmd

	state tState

	message string
}

func New(initCmd tea.Cmd, onSuccessCmd tea.Cmd, onFailCmd tea.Cmd) Model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		spinner:      s,
		initCmd:      initCmd,
		onSuccessCmd: onSuccessCmd,
		onFailCmd:    onFailCmd,

		state: stateLoading,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.initCmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateSuccess:
			return m, m.onSuccessCmd
		case stateFail:
			return m, m.onFailCmd
		}

	case ResponseResultMsg:
		if msg.Success {
			m.state = stateSuccess
		} else {
			m.state = stateFail
		}
		m.message = msg.Message
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.state != stateLoading {
		return fmt.Sprintf("\n\n   %s\n   Press any key to continue.\n\n", m.message)
	}

	return fmt.Sprintf("\n\n   %s Waiting for response...\n\n", m.spinner.View())
}
