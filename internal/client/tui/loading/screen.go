package loading

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type nextFunc func() (tea.Model, tea.Cmd)

type DoneMsg struct {
	next    nextFunc
	message string
}

func NewDoneMsg(message string, next nextFunc) DoneMsg {
	return DoneMsg{
		message: message,
		next:    next,
	}
}

type Screen struct {
	spinner  spinner.Model
	doneMsg  *DoneMsg
	err      error
	doneCb   tea.Cmd
	fnCancel context.CancelFunc
}

func NewScreen(doneCb tea.Cmd, cancel context.CancelFunc) *Screen {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &Screen{
		spinner:  s,
		doneCb:   doneCb,
		fnCancel: cancel,
	}
}

func (ths Screen) Init() tea.Cmd {
	return tea.Batch(ths.spinner.Tick, ths.doneCb)
}

func (ths Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return ths, tea.Quit
		}

		if ths.doneMsg != nil {
			return ths.doneMsg.next()
		}

	case DoneMsg:
		ths.doneMsg = &msg
	}

	if ths.doneMsg != nil {
		return ths, nil
	}

	var cmd tea.Cmd
	ths.spinner, cmd = ths.spinner.Update(msg)
	return ths, cmd
}

func (ths Screen) View() string {
	if ths.doneMsg != nil {
		return fmt.Sprintf("\n\n   %s\n   Press any key to continue.\n\n", ths.doneMsg.message)
	}

	return fmt.Sprintf("\n\n   %s Waiting for response...\n\n", ths.spinner.View())
}
