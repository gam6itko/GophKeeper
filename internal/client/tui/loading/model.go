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

	LoginResponseMsg struct {
		Err error
	}
)

func closeCmd() tea.Msg {
	return CloseMsg{}
}

//type nextFunc func() (tea.Model, tea.Cmd)
//
//type DoneMsg struct {
//	message string
//	cmd    tea.Cmd
//}
//
//func NewDoneMsg(message string, cmd Cmd) DoneMsg {
//	return DoneMsg{
//		message: message,
//		cmd:    cmd,
//	}
//}

type tState int

const (
	stateLoading tState = iota
	stateSuccess
	stateFail
)

type Model struct {
	spinner spinner.Model
	//Err      error
	initCmd tea.Cmd

	state tState
	prev  tea.Model

	message string
}

func New(initCmd tea.Cmd, prev tea.Model) Model {
	s := spinner.New()
	s.Spinner = spinner.Hamburger
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		spinner: s,
		initCmd: initCmd,

		state: stateLoading,
		prev:  prev,
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
			return m, func() tea.Msg {
				return SuccessMsg{}
			}
		case stateFail:
			return m, func() tea.Msg {
				return FailMsg{GoToModel: m.prev}
			}
		}

	case LoginResponseMsg:
		if msg.Err == nil {
			m.state = stateSuccess
			m.message = "Login successful"
		} else {
			m.state = stateFail
			m.message = msg.Err.Error()
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.state != stateLoading {
		return fmt.Sprintf("\n\n  Error: %s\n   Press any key to continue.\n\n", m.message)
	}

	return fmt.Sprintf("\n\n   %s Waiting for response...\n\n", m.spinner.View())
}
